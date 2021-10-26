package main

import (
	"fmt"
	"os"

	"github.com/infastin/gops/pkg"
)

func main() {
	procs, err := ps.Processes()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range procs {
		fmt.Println(v)
	}
}
