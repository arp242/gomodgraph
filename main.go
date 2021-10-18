package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"zgo.at/zli"
	"zgo.at/zstd/zos"
	"zgo.at/zstd/zstring"
)

func main() {
	showVersion := len(os.Args) > 1
	os.Args = []string{"x"}

	fp, err := zli.InputOrFile(zos.Arg(1), false)
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
		if showVersion {
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
		printpkg(p, packages, 0, "")
	}
}

const indent = "\t"

func printpkg(p string, packages map[string][]string, i int, parent string) {
	fmt.Printf("%s%s\n", strings.Repeat(indent, i), p)
	for _, d := range packages[p] {
		/*
		   golang.org/x/text
		   	golang.org/x/tools
		   		github.com/yuin/goldmark
		   		golang.org/x/mod
		   			golang.org/x/crypto
		   			golang.org/x/net
		   				golang.org/x/sys
		   				golang.org/x/term
		   					golang.org/x/sys
		   				golang.org/x/text
		*/
		if strings.HasPrefix(d, "golang.org/x/") {
			continue
		}

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
