package main

import (
	"flag"
	"log"

	"github.com/jacoelho/initcheck"
)

func main() {
	flag.Parse()

	if err := initcheck.Run(flag.Args()); err != nil {
		log.Fatal(err)
	}

}
