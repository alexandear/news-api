package main

import (
	"fmt"
	"os"

	"github.com/alexandear/news-api/cmd"
)

func main() {
	if err := cmd.NewRoot().Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
