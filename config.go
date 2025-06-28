package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/howeyc/gopass"
)

// Config struct
type Config struct {
	ServerAddress string `json:"ServerAddress"`
	Username      string `json:"Username"`
	Password      string `json:"Password"`
	UploadDir     string `json:"UploadDir"`
}

// ReadConfig returns an error if file does not exist
func readConfig(file string) error {
	configJSON, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configJSON, &config)

	if err != nil {
		return err
	}

	return nil
}

// WriteConfig writes a config file
func writeConfig(file string) error {
	config = Config{}

	fmt.Printf("Webdav server: ")
	if _, err := fmt.Scanln(&config.ServerAddress); err != nil {
		return err
	}

	fmt.Printf("Username: ")
	if _, err := fmt.Scanln(&config.Username); err != nil {
		return err
	}

	fmt.Printf("Password (not displayed): ")
	pwd, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	config.Password = string(pwd)

	fmt.Printf("Default upload directory: ")
	if _, err := fmt.Scanln(&config.UploadDir); err != nil {
		return err
	}

	configJSON, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(file, configJSON, 0600)
}

// Home returns the user's home directory
func home() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
