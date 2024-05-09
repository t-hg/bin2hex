package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func printHelpAndExit() {
	help := `Usage: bin2hex [FILE]

Converts a binary to hex decimal representation like:
...
06993f70  09 43 24 7b af 98 76 b7  c4 ab 15 89 83 88 50 5d  |.C${..v.......P]|
06993f80  96 00 e5 cf 97 b7 90 12  cc 06 e4 82 93 16 63 7c  |..............c||
06993f90  cc 03 f0 91 43 40 ba 1f  5d c3 aa 75 4e c2 c6 0c  |....C@..]..uN...|
06993fa0  13 5f 7a 79 40 c7 6b e2  67 db 4e 1b 27 a3 31 09  |._zy@.k.g.N.'.1.|
...

If no FILE has been given, the tool will read from STDIN.
`
	fmt.Fprintf(os.Stderr, help)
	flag.PrintDefaults()
	os.Exit(1)
}

func printLine(address int, hex []string, str string) {
	hexJoined1 := strings.Join(hex[:8], " ")
	hexJoined2 := ""
	if len(hex) > 8 {
		hexJoined2 = strings.Join(hex[8:], " ")
	}
	fmt.Printf("%08x  %-23s  %-23s  |%-16s|\n", address, hexJoined1, hexJoined2, str)
}

func main() {
	flag.Usage = printHelpAndExit
	flag.Parse()
	if flag.NArg() > 1 {
		return
	}
	var reader io.Reader
	if flag.NArg() > 0 {
		file := flag.Arg(0)	
		var err error
		reader, err = os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read %s: %v", file, err)
			os.Exit(1)
			return
		}
	} else {
		reader = os.Stdin
	}
	bufferedReader := bufio.NewReader(reader)
	address := 0
	read := 0
	hex := make([]string, 0)
	str := ""
	for  {
		if read == 16 {
			printLine(address, hex, str)
			address++
			read = 0
			hex = make([]string, 0)
			str = ""
			continue
		}
		b, err := bufferedReader.ReadByte()
		if err != nil {
			printLine(address, hex, str)
			break
		}
		hex = append(hex, fmt.Sprintf("%02x", b))
		if b >= 32 && b <= 126 {
			str += string(b)
		} else {
			str += "."
		}
		read++
	}
}
