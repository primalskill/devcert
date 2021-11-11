package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/fs"
	"os"
)

// needsSetup checks if the setup process needs to be executed.
func needsSetup() (need bool, err error) {
	isFolder, err := isDevcertFolder()
	if err != nil {
		return
	}

	if isFolder == false {
		return true, nil
	}

	isCA, err := isValidCA()
	if err != nil {
		return
	}

	if isCA == false {
		return true, nil
	}

	return false, nil
}

// isDevcertFolder checks if the necessary folder exists.
func isDevcertFolder() (is bool, err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(devcertDir)
	notExist := errors.Is(err, fs.ErrNotExist)

	if err != nil && notExist == false {
		return false, err
	}

	is = notExist == false
	return is, nil
}

// isValidCA check if the certificate authority files are valid.
func isValidCA() (is bool, err error) {
	ca, err := loadCA()
	if err != nil {
		return
	}

	is = ca.Valid

	return
}

// validateSetupPrompt returns a validation error if setupPrompt passes the wrong value.
func validateSetupPrompt(input string) (err error) {
	if input != "Y" && input != "y" && input != "n" && input != "N" {
		err = fmt.Errorf("Invalid value")
	}
	return
}

// setupPrompt will ask the user to continue or not.
func setupPrompt() (err error) {
	prompt := promptui.Prompt{
		Label:    "Do you want to continue? [Y/n]",
		Validate: validateSetupPrompt,
	}

	result, err := prompt.Run()
	if err != nil {
		return
	}

	if result == "n" || result == "N" {
		os.Exit(0)
	}

	return nil
}

// setup will start the setup process.
func setup() (err error) {
	devcertDir, err := buildDevcertDir()
	if err != nil {
		return
	}

	fmt.Printf("\ndevcert needs to execute the setup process first.")
	fmt.Printf("\n  - It will create %s directory.", devcertDir)
	fmt.Printf("\n  - It will create a local certificate authority (CA) to sign future certificates.")
	fmt.Printf("\n  - It will mark the CA as trusted locally.\n")

	err = setupPrompt()
	if err != nil {
		return
	}

	// Continue the setup process.
	err = createDevcertDir()
	if err != nil {
		err = fmt.Errorf("Setup failed: %w", err)
		attemptCleanupDevcertDir()
		return
	}

	err = createCA()
	if err != nil {
		err = fmt.Errorf("Setup failed: %w", err)
		attemptCleanupDevcertDir()
		return
	}

	err = trustCA()
	if err != nil {
		err = fmt.Errorf("Setup failed: %w", err)
		attemptCleanupCA()
		return
	}

	return
}

// createDevcertDir creates the .devcert directory in the user's home directory.
func createDevcertDir() (err error) {
	fmt.Printf("Creating directory...\n")

	devcertDir, err := buildDevcertDir()
	if err != nil {
		err = fmt.Errorf("Creating .devcert directory failed: %w", err)
		return
	}

	isDevcertDir, err := isDevcertFolder()
	if err != nil {
		err = fmt.Errorf("Creating .devcert directory failed: %w", err)
		return
	}

	// The directory already exists.
	if isDevcertDir == true {
		fmt.Printf("Directory %s already created.\n", devcertDir)
		return
	}

	// Create the directory
	err = os.MkdirAll(devcertDir, 0755)
	if err != nil {
		fmt.Errorf("Creating .devcert directory failed: %w", err)
		return
	}

	fmt.Printf("Directory %s created.\n", devcertDir)

	return
}
