package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// getHomeDir returns the current user home directory.
func getHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// buildDevcertDir returns the constructed devcert directory.
func buildDevcertDir() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}

	var devcerts strings.Builder
	devcerts.WriteString(home)
	devcerts.WriteString(string(os.PathSeparator))
	devcerts.WriteString(".devcert")
	devcerts.WriteString(string(os.PathSeparator))

	return devcerts.String(), nil
}

func sudoify(cmd ...string) *exec.Cmd {
	_, err := exec.LookPath("sudo")
	if err != nil {
		fmt.Printf("\n'sudo' is not installed on the system, devcert might fail.")
		return exec.Command(cmd[0], cmd[1:]...)
	}

	sudoCmd := []string{"--prompt=Enter your sudo password:", "--"}
	sudoCmd = append(sudoCmd, cmd...)

	return exec.Command("sudo", sudoCmd...)
}
