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
		fmt.Printf("'sudo' is not installed on the system, devcert might fail.\n")
		return exec.Command(cmd[0], cmd[1:]...)
	}

	sudoCmd := []string{"--prompt=Enter your sudo password:", "--"}
	sudoCmd = append(sudoCmd, cmd...)

	return exec.Command("sudo", sudoCmd...)
}

func isDirectoryExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return err == nil
}

func isBinaryExist(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

// detectLinux tries to fingerprint Linux CA trust methods instead of detecting the Linux distro itself.
func detectLinux() (base string, err error) {
	if isDirectoryExist("/usr/local/share/ca-certificates/") == true && isBinaryExist("update-ca-certificates") == true {
		return "debian", nil
	}

	if isDirectoryExist("/etc/pki/ca-trust/source/anchors/") == true && isBinaryExist("update-ca-trust") == true {
		return "rhel", nil
	}

	if isBinaryExist("pacman") == true && isBinaryExist("trust") == true {
		return "arch", nil
	}

	err = fmt.Errorf("Detecting Linux failed: Cannot detect base Linux distro.")
	return
}
