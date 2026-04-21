package main

import (
	"flag"
)

func main() {
	part := flag.String("part", "1", "which part to run: 1 or 2")
	flag.Parse()
	
	if *part == "1" {
		main1()
	} else {
		main2()
	}
}