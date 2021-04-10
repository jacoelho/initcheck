package main

import (
	"flag"
	"log"

	"github.com/jacoelho/initanalysis"
)

func main() {

	flag.Parse()
	if err := initanalysis.Run(flag.Args()); err != nil {
		log.Fatal(err)
	}

}
