package main

import (
	"fmt"

	"github.com/rellyson/http-echo/pkg/version"
)

func main() {
	version, err := version.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello world from HTTP Echo version: %s build: %s!\n", version.Version, version.Build)
}
