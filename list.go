package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/go2c/optparse"
)

// list lists ports.
func list(input []string) {
	// Define valid arguments.
	o := optparse.New()
	argi := o.Bool("installed", 'i', false)
	argr := o.Bool("repo", 'r', false)
	argv := o.Bool("version", 'v', false)
	argh := o.Bool("help", 'h', false)

	// Parse arguments.
	_, err := o.Parse(input)
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: prt list [arguments]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -i,   --installed       list installed ports only")
		fmt.Println("  -r,   --repo            list with repo info")
		fmt.Println("  -v,   --version         list with version info")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	// Get all ports.
	all, err := ports()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var db database
	if *argi {
		// Get installed ports.
		db, err = parseDatabase()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Get port locations.
		var pl []port
		for _, n := range db.Name {
			p, err := location(n, all)
			if err != nil {
				continue
			}
			pl = append(pl, p[0])
		}

		// I'm using all in the the following for loop, so alias db to
		// all.
		all = pl
	}

	var pl []string
	for i, p := range all {
		var s string

		if !*argr {
			s = p.getPortDir()
		} else {
			s = p.getBaseDir()
		}

		if *argv {
			if *argi {
				s += " " + db.Version[i]
			} else {
				if err := p.parsePkgfile(); err != nil {
					printe(err.Error())
					continue
				}
				s += " " + p.Pkgfile.Version
			}
		}

		pl = append(pl, s)
	}

	sort.Strings(pl)
	for _, p := range pl {
		fmt.Println(p)
	}
}
