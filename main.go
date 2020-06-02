package main

import (
	"bufio"
	"fmt"
	"strings"

	"zgo.at/zli"
	"zgo.at/zstd/zos"
	"zgo.at/zstd/zstring"
)

func main() {
	fp, err := zli.FileOrInput(zos.Arg(1))
	if err != nil {
		zli.Fatal(err)
	}
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
			zli.Fatal("malformed line: %q", line)
		}

		pkg := strings.Split(s[0], "@")[0]
		dep := strings.Split(s[1], "@")[0]

		_, ok := packages[pkg]
		if !ok {
			order = append(order, pkg)
		}
		packages[pkg] = append(packages[pkg], dep)
	}
	if err := scan.Err(); err != nil {
		zli.Fatal(err)
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
		printpkg(p, packages, 0, "")
	}
}

const indent = "\t"

func printpkg(p string, packages map[string][]string, i int, parent string) {
	fmt.Printf("%s%s\n", strings.Repeat(indent, i), p)
	for _, d := range packages[p] {
		// x/net depends on x/crypto which depends on x/net
		if d == parent {
			continue
		}

		_, ok := packages[d]
		if ok {
			printpkg(d, packages, i+1, p)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(indent, i+1), d)
		}
	}
}
