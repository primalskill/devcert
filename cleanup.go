package main

import (
	"os"
)

// attemptCleanupDevcertDir attempts to remove the .devcert directory. Errors are suppressed.
func attemptCleanupDevcertDir() {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return
	}

	os.Remove(devcertDir)
}

// attemptCleanupCA attempts to remove the CA files. Errors are suppressed.
func attemptCleanupCA() {
	crtPath, keyPath, err := buildCAPaths()
	if err != nil {
		return
	}

	os.Remove(crtPath)
	os.Remove(keyPath)
}
