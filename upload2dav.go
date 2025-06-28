package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/axllent/ghru/v2"
	"github.com/spf13/pflag"
	"github.com/studio-b12/gowebdav"
)

var (
	config  Config
	client  *gowebdav.Client
	quiet   bool
	version = "dev"
)

func main() {
	ghruConf := ghru.Config{
		Repo:           "axllent/upload2dav",
		ArchiveName:    "upload2dav-{{.OS}}-{{.Arch}}",
		BinaryName:     "upload2dav",
		CurrentVersion: version,
	}

	var showHelp, writeConfig, showVersion, update bool
	var configFile, uploadPath string

	defaultConfig := Home() + "/.config/upload2dav.json"

	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// set the default help
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file(s)>\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.SortFlags = false
		flag.PrintDefaults()
		fmt.Println("")
	}

	flag.StringVarP(&uploadPath, "dir", "d", "", "Alternative upload directory")
	flag.StringVarP(&configFile, "conf", "c", defaultConfig, "Specify config file")
	flag.BoolVarP(&writeConfig, "write-config", "w", false, "Write config")
	flag.BoolVarP(&quiet, "quiet", "q", false, "Quiet (do not show upload progress)")
	flag.BoolVarP(&showVersion, "version", "v", false, "Show version")
	flag.BoolVarP(&update, "update", "u", false, "Update to latest version")
	flag.BoolVarP(&showHelp, "help", "h", false, "Show help")

	// parse args excluding os.Args[0]
	flag.Parse(os.Args[1:])

	// parse arguments
	files := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(1)
	}

	if showVersion {
		fmt.Printf("Version: %s\n", version)

		release, err := ghruConf.Latest()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// The latest version is the same version
		if release.Tag == version {
			os.Exit(0)
		}

		// A newer release is available
		fmt.Printf(
			"Update available: %s\nRun `%s -u` to update (requires read/write access to install directory).\n",
			release.Tag,
			os.Args[0],
		)
		os.Exit(0)
	}

	if update {
		// Update the app
		rel, err := ghruConf.SelfUpdate()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Updated %s to version %s\n", os.Args[0], rel.Tag)
		os.Exit(0)
	}

	if writeConfig {
		if err := WriteConfig(configFile); err != nil {
			fmt.Printf("Error: %s\n\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully wrote config: %s\n\n", configFile)
		os.Exit(0)
	}

	if err := ReadConfig(configFile); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("\nUse -c to specify a configuration file, or -w to create a new one\n\n")
		os.Exit(1)
	}

	if len(files) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if uploadPath != "" {
		config.UploadDir = uploadPath
	}

	client = gowebdav.NewClient(config.ServerAddress, config.Username, config.Password)

	if err := CheckDirExists(config.UploadDir); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for _, file := range files {
		if err := Upload(file, config.UploadDir); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}

// Upload sends a local file to the webdav server
func Upload(file, dir string) error {
	info, err := os.Stat(file)

	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a file", file)
	}

	wfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer wfile.Close()

	outFilename := filepath.Base(file)

	uploadName := path.Join(config.UploadDir, outFilename)

	if !quiet {
		fmt.Printf("Uploading %s to %s ... ", file, uploadName)
	}

	if err := client.WriteStream(uploadName, wfile, 0664); err != nil {
		return err
	}

	if !quiet {
		fmt.Println("done")
	}

	return nil
}

// CheckDirExists checked first is a directory exists
func CheckDirExists(dir string) error {
	if _, err := client.ReadDir(dir); err != nil {
		if err := client.MkdirAll(dir, 0644); err != nil {
			return err
		}
	}

	return nil
}
