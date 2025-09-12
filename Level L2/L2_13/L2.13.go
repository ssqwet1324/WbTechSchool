package main

import (
	"bufio"
	"flag"
	"os"
)

type Flags struct {
	Fields    int
	Delimiter string
	Separator string
}

func Cut(scanner *bufio.Scanner, flags *Flags) {

}

func main() {
	fields := flag.Int("fields", 1, "Number of fields")
	delimiter := flag.String("delimiter", "\t", "Delimiter")
	separator := flag.String("separator", " ", "Separator")
	flag.Parse()

	flags := Flags{
		Fields:    *fields,
		Delimiter: *delimiter,
		Separator: *separator,
	}

	var scanner *bufio.Scanner
	scanner = bufio.NewScanner(os.Stdin)

	Cut(scanner, &flags)

}
