package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

const usage = `gomodgraph takes the output of "go mod graph" and prints it as a tree.

https://github.com/arp242/gomodgraph

Usage:

    % go mod graph | gomodgraph

Flags:

    -h, -help       Show this help.
    -v, -version    Also show the version of modules.
    -d, -depth      Set maximum depth to print.
`

func main() {
	f := zli.NewFlags(os.Args)
	var (
		showVersion = f.Bool(false, "v", "version")
		depth       = f.Int(0, "d", "depth")
		help        = f.Bool(false, "h", "help")
	)
	zli.F(f.Parse())

	if help.Bool() {
		fmt.Println(usage)
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

	root := order[0]
	for _, p := range order {
		if p != root {
			break
		}
		printpkg(p, packages, 0, nil, depth.Int(), 1)
	}
}

const indent = "\t"

// TODO: "// indirect" packages are listed here too as direct dependencies; this
// is how go mod graph outputs it. Should filter those really.
func printpkg(p string, packages map[string][]string, i int, parents []string, maxDepth, curDepth int) {
	fmt.Printf("%s%s\n", strings.Repeat(indent, i), p)

loop:
	for _, d := range packages[p] {
		for _, p := range parents {
			if d == p {
				continue loop
			}
		}

		// Don't print dependencies for the golang.org/x/... packages as 1) many
		// packages depend on them, and 2) many of them depend on each other.
		// This leads to a huge amount of noise in the output.
		_, ok := packages[d]
		if ok && !strings.HasPrefix(d, "golang.org/x/") {
			if maxDepth == 0 || curDepth < maxDepth {
				printpkg(d, packages, i+1, append(parents, p), maxDepth, curDepth+1)
			}
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(indent, i+1), d)
		}
	}
}
