package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "data/test.tsv"
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	in := bufio.NewScanner(os.Stdin)
	out := os.Stdout

	for s.Scan() {
		line := s.Text()
		splitedLine := strings.Split(line, "\t")
		fmt.Printf("(p)Positive (f)Flat (n)Negative (d)Drop: %s\n", splitedLine[2])
		in.Scan()
		if s.Err() != nil {
			log.Println(err)
			os.Exit(1)
		}
		switch in.Text() {
		case "p":
			fmt.Fprintln(out, "\"Positive\","+splitedLine[2])
		case "f":
			fmt.Fprintln(out, "\"Flat\","+splitedLine[2])
		case "n":
			fmt.Fprintln(out, "\"Negative\","+splitedLine[2])
		case "d":
			continue
		default:
			log.Println("Interupted")
			os.Exit(1)
		}

	}
	if s.Err() != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("vim-go")
}
