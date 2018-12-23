package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/rand-hashing/internal/hashgen"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./gethashes key iterations")
		os.Exit(0)
	}
	conf := os.Args[1]
	key := os.Args[2]
	n, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("argument parsing error: ", err)
		os.Exit(0)
	}

	err = hashgen.StartGeneration(conf, key, n)
	if err != nil {
		fmt.Println("generation error: ", err)
		os.Exit(0)
	}
}
