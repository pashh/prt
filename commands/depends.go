package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chiyouhen/getopt"
	"github.com/fatih/color"
	"github.com/onodera-punpun/prt/config"
	"github.com/onodera-punpun/prt/pkgfile"
	"github.com/onodera-punpun/prt/ports"
	"github.com/onodera-punpun/prt/utils"
)

func recursive(path string) {
	// Read out Pkgfile
	f, err := ioutil.ReadFile(path + "/Pkgfile")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read '"+path+"/Pkgfile'!")
		return
	}

	// Read out Pkgfile dependencies
	deps, err := pkgfile.Depends(f, "Depends on")
	if err != nil {
		return
	}

	for _, dep := range deps {
		// Continue if already checked
		if utils.StringInList(dep, checkPorts) {
			continue
		}
		checkPorts = append(checkPorts, dep)

		// Get port location
		locs, err := ports.Loc(allPorts, dep)
		if err != nil {
			continue
		}
		loc := locs[0]

		// Alias if needed
		if !utils.StringInList("n", optsList) {
			loc = ports.Alias(loc)
		}

		// Continue if already installed
		if !utils.StringInList("a", optsList) {
			if utils.StringInList(filepath.Base(loc), instPorts) {
				continue
			}
		}

		// Print tree indentation
		if !utils.StringInList("t", optsList) {
			if i > 0 {
				color.Set(color.FgBlack, color.Bold)
				fmt.Printf(strings.Repeat(config.Struct.IndentChar, i))
				color.Unset()
			}
			i++
		}

		// Finally print the port :)
		fmt.Println(loc)

		// Loop
		recursive(config.Struct.PortDir + "/" + loc)

		if !utils.StringInList("t", optsList) {
			i--
		}
	}
}

// Depends lists dependencies recursivly
func Depends(args []string) {
	// Define opts
	shortopts := "hant"
	longopts := []string{
		"--help",
		"--no-alias",
		"--no-tree",
	}

	// Read out opts
	opts, _, err := getopt.Getopt(args, shortopts, longopts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	for _, opt := range opts {
		switch opt[0] {
		case "-h", "--help":
			fmt.Println("Usage: prt depends [arguments]")
			fmt.Println("")
			fmt.Println("arguments:")
			fmt.Println("  -a,   --all             also list installed dependencies")
			fmt.Println("  -n,   --no-alias        disable aliasing")
			fmt.Println("  -t,   --no-tree         list without tree view")
			fmt.Println("  -h,   --help            print help and exit")
			os.Exit(0)
		case "-a", "--all":
			optsList = append(optsList, "a")
		case "-n", "--no-alias":
			optsList = append(optsList, "n")
		case "-t", "--tree":
			optsList = append(optsList, "t")
		}
	}

	allPorts, err = ports.All()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if utils.StringInList("a", optsList) {
		instPorts, err = ports.Inst()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	recursive("./")
}
