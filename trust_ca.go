package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// trustCA marks the certificate authority trusted locally.
func trustCA() (err error) {
	fmt.Printf("Trusting certificate authority...\n")

	ca, err := loadCA()
	if err != nil {
		err = fmt.Errorf("Trusting CA failed: %w", err)
		return
	}

	// Certificate is valid, nothing to do.
	if ca.Valid == false {
		err = fmt.Errorf("Trusting CA failed: Certificate is invalid.")
		return
	}

	crtPath, _, err := buildCAPaths()
	if err != nil {
		err = fmt.Errorf("Trusting CA failed: %w", err)
		return
	}

	switch runtime.GOOS {
	case "darwin":
		err = trustDarwin(crtPath)
	case "windows":
		err = trustWindows(crtPath)
	}

	if err == nil {
		fmt.Printf("Certificate authority (CA) marked trusted.\n")
	}

	return
}

func trustDarwin(crtPath string) (err error) {
	stdOutStdError, err := sudoify("security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", crtPath).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Trusting CA failed: %s, %w", stdOutStdError, err)
		return
	}

	return
}

func trustLinux(crtPath, keyPath string) (err error) {
	return
}

func trustWindows(crtPath string) (err error) {
	var argList strings.Builder
	argList.WriteString("-ArgumentList '-addstore -f ROOT ")
	argList.WriteString(crtPath)
	argList.WriteString("'")

	stdOutStdError, err := exec.Command("powershell", "Start-Process -FilePath certutil -Verb RunAs -Wait -PassThru", argList.String()).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Trusting CA failed: %s, %w", stdOutStdError, err)
		return
	}

	return
}
