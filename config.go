package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
func ReadConfig(file string) error {
	configJSON, err := ioutil.ReadFile(file)
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
func WriteConfig(file string) error {
	config = Config{}

	fmt.Printf("Webdav server: ")
	fmt.Scanln(&config.ServerAddress)

	fmt.Printf("Username: ")
	fmt.Scanln(&config.Username)

	fmt.Printf("Password (not displayed): ")
	pwd, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	config.Password = string(pwd)

	fmt.Printf("Default upload directory: ")
	fmt.Scanln(&config.UploadDir)

	configJSON, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, configJSON, 0600)
}

// Home returns the user's home directory
func Home() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
