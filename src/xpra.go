package main

import (
	"log"
	"os"
	"os/exec"
)

func XpraConnect(port string, password string) int {
	log.Println("connecting to xpra server")
	addressArg := "tcp:localhost:" + port
	passwordArg := "--password=" + password
	cmd := exec.Command("xpra", "attach", addressArg, "--auth=password", passwordArg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("xpra server connected", cmd.Process.Pid)
	return cmd.Process.Pid
}
