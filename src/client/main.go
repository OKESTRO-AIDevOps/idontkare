package main

import (
	"fmt"
	"os"
)

var CLIENT_CONFIG *ClientConfig

func main() {

	var err error

	CLIENT_CONFIG, err = GetClientConfig()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	if CLIENT_CONFIG.ApixVersion == 1 {

		if CLIENT_CONFIG.SocketMode {

			// TODO:
			//   add socket mode, V1Run

		} else {

			result, err := V1RunOnce(CLIENT_CONFIG.ConnectUrl)

			if err != nil {

				fmt.Fprintf(os.Stderr, "%s\n", err.Error())

			} else {

				fmt.Fprintf(os.Stdout, "%s\n", result)
			}

		}

	} else {

		fmt.Fprintf(os.Stderr, "wrong apix version: %d\n", CLIENT_CONFIG.ApixVersion)
		return
	}

	return
}
