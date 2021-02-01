package main

import (
	"fmt"
	"os"

	"github.com/Ken2mer/whois"
)

func main() {
	res, err := whois.Fetch("example.com")
	FatalIf(err)

	fmt.Println(res.String())
}

func FatalIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(-1)
}
