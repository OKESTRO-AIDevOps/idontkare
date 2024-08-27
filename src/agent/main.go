package main

import (
	"fmt"
	"os"
)

var AGENT_CONFIG *AgentConfig

func main() {

	var err error

	AGENT_CONFIG, err = GetAgentConfig()

	if err != nil {

		fmt.Fprintf(os.Stderr, "%s\n", err.Error())

		return
	}

	if AGENT_CONFIG.ApixVersion == 1 {

		err := V1Run(
			AGENT_CONFIG.ConnectUrl,
			AGENT_CONFIG.Name,
			AGENT_CONFIG.UserName,
			AGENT_CONFIG.KeyPath,
		)

		if err != nil {

			fmt.Fprintf(os.Stderr, "failed to run: %s\n", err.Error())
		}

	} else {

		fmt.Fprintf(os.Stderr, "wrong apix version: %d\n", AGENT_CONFIG.ApixVersion)
		return
	}

	return
}
