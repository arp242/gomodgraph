package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

const usage = `gomodgraph takes the output of "go mod graph" and prints it as a tree.

https://github.com/arp242/gomodgraph

Usage:

    % go mod graph | gomodgraph

Flags:

    -h, -help          Show this help.
    -v, -version       Show gomodgraph's version.
    -n, -no-color      Don't output colour.
    -d, -depth         Set maximum depth to print.
    -r, -repeat        Repeat dependencies; by default if a package is seen more
                       than once the dependency tree is not printed the second
                       and subsequent times. This will always print all
                       dependencies.
    -x, -with-x        Also display the golang.org/x/[..] packages; by default
                       they're not included as they're so common and "kind-of
                       stdlib". Note that dependencies of /x/ packages are never
                       printed, as many depend on eachother and gives a very
                       noisy output. This also includes github.com/golang/protobuf.
    -V, -with-version  Also show the version of modules.
`

func main() {
	f := zli.NewFlags(os.Args)
	var (
		help        = f.Bool(false, "h", "help")
		versionF    = f.IntCounter(0, "v", "version")
		showVersion = f.Bool(false, "V", "with-version")
		noColor     = f.Bool(false, "n", "no-color")
		depth       = f.Int(0, "d", "depth")
		withX       = f.Bool(false, "x", "with-x")
		repeat      = f.Bool(false, "r", "repeat")
	)
	zli.F(f.Parse())

	zli.WantColor = !noColor.Bool()
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		zli.WantColor = false
	}

	if help.Bool() {
		fmt.Print(usage)
		return
	}
	if versionF.Int() > 0 {
		zli.PrintVersion(versionF.Int() > 1)
		return
	}

	fp, err := zli.InputOrFile(f.Shift(), false)
	zli.F(err)
	defer fp.Close()
	scan := bufio.NewScanner(fp)

	var (
		order    []string
		packages = make(map[string][]string)
	)
	for scan.Scan() {
		line := scan.Text()
		s := strings.Split(line, " ")
		if len(s) != 2 {
			zli.Fatalf("malformed line: %q", line)
		}

		// TODO: abbrev version in case it's not a tag
		var pkg, dep string
		if showVersion.Bool() {
			pkg = strings.Replace(s[0], "@", " ", -1)
			dep = strings.Replace(s[1], "@", " ", -1)
		} else {
			pkg = strings.Split(s[0], "@")[0]
			dep = strings.Split(s[1], "@")[0]
		}

		_, ok := packages[pkg]
		if !ok {
			order = append(order, pkg)
		}
		packages[pkg] = append(packages[pkg], dep)
	}
	if err := scan.Err(); err != nil {
		zli.F(err)
	}
	if len(order) == 0 {
		return
	}

	for p := range packages {
		packages[p] = zstring.Uniq(packages[p])
	}

	indir := indirects()
	root := order[0]
	for _, p := range order {
		if p != root {
			break
		}
		printpkg(p, packages, 0, nil, withX.Bool(), repeat.Bool(), depth.Int(), 1, indir, make(map[string]struct{}))
	}
}

// "// indirect" packages are listed as direct dependencies; this is how go mod
// graph outputs it.
func indirects() []string {
	out, err := exec.Command("go", "list", "-f", "{{.Indirect}} {{.Path}}", "-m", "all").CombinedOutput()
	if err != nil {
		zli.Errorf("running go list: %s:\n%s", err, out)
	}

	in := make([]string, 0, 8)
	for _, line := range strings.Split(string(out), "\n") {
		indir, pkg := zstring.Split2(line, " ")
		if indir == "true" {
			in = append(in, pkg)
		}
	}
	return in
}

func indent(n int) string {
	if n == 0 {
		return ""
	}
	if n == 1 {
		return "\t"
	}
	return "\t" + strings.Repeat(zli.Color256(254).String()+"│\t"+zli.Reset.String(), n-1)
}

func printpkg(
	p string,
	packages map[string][]string,
	i int,
	parents []string,
	withX, repeat bool,
	maxDepth, curDepth int,
	indir []string,
	seen map[string]struct{},
) {
	fmt.Printf("%s%s\n", indent(i), p)

loop:
	for _, d := range packages[p] {
		if curDepth == 1 && zstring.Contains(indir, d) {
			continue
		}

		for _, p := range parents {
			if d == p {
				continue loop
			}
		}

		xPkg := strings.HasPrefix(d, "golang.org/x/") || d == "github.com/golang/protobuf"
		if xPkg && !withX {
			continue
		}

		// Don't print dependencies for the golang.org/x/... packages as 1) many
		// packages depend on them, and 2) many of them depend on each other.
		// This leads to a huge amount of noise in the output.
		_, ok := packages[d]
		if ok && !xPkg {
			if maxDepth == 0 || curDepth < maxDepth {
				if !repeat {
					if _, ok := seen[d]; ok {
						fmt.Printf("%s%s\n%s…\n", indent(i+1), d, indent(i+2))
						continue
					}
					seen[d] = struct{}{}
				}
				printpkg(d, packages, i+1, append(parents, p), withX, repeat, maxDepth, curDepth+1, indir, seen)
			}
		} else {
			fmt.Printf("%s%s\n", indent(i+1), d)
		}
	}
}
