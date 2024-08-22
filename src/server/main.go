package main

import (
	"fmt"
	"os"
)

var SERVER_CONFIG *ServerConfig

func main() {

	var err error

	SERVER_CONFIG, err = GetServerConfig()

	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	if SERVER_CONFIG.ApixVersion == 1 {

	} else {

		fmt.Fprintf(os.Stderr, "wrong apix version: %d\n", SERVER_CONFIG.ApixVersion)
		return
	}

}
