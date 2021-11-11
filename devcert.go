package main

import (
	"fmt"
)

func devcertExec(args []string) (err error) {

	// Decide to execute setup
	need, err := needsSetup()
	if err != nil {
		err = fmt.Errorf("Setup failed: %w", err)
		return
	}

	if need == true {
		err = setup()
		if err != nil {
			err = fmt.Errorf("Setup failed: %w", err)
			return
		}
	}

	err = generateCertificate(args)
	if err != nil {
		err = fmt.Errorf("Generate certificate failed: %w", err)
		return
	}

	return
}
