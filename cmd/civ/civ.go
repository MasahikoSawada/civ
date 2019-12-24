package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/MasahikoSawada/civ"
	termbox "github.com/nsf/termbox-go"
	"os"
)

func main() {
	var d rune

	dd := flag.String("d", ",", "use the delimiter instead of comma")
	withoutHeader := flag.Bool("H", false, "set dummy header")
	in := os.Stdin
	flag.Parse()

	if len(flag.Args()) > 1 {
		panic("civ can accept only one file")
	}

	// File is specified
	if len(flag.Args()) != 0 {
		_in, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Printf("could not open file %s: %s\n", flag.Args()[0], err)
			os.Exit(1)
		}
		defer in.Close()
		in = _in
	}

	// Determine the delimiter. We allow special character
	// '\t' which represents a tab for ease of use.
	if *dd == "\\t" {
		d = '\t'
	} else {
		d = rune((*dd)[0])
	}

	// read file
	reader := csv.NewReader(in)
	reader.Comma = d
	csv, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("could not read data: %s\n", err)
		os.Exit(1)
	}

	// initialize termbox
	if err := termbox.Init(); err != nil {
		fmt.Printf("could not initialize termbox: %s\n", err)
		os.Exit(1)
	}

	c := civ.NewCiv(csv, *withoutHeader)
	exited := c.Run()

	termbox.Close()

	if exited {
		c.Draw()
	}
}
