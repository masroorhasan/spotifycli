package main

import (
	"log"
	"os"

	"github.com/masroorhasan/spotifycli/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
