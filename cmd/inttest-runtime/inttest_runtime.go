package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const (
	version     = "0.1.0"
	versionDate = "30.10.2023"
)

type cliArg string

const (
	configPathArg cliArg = "config"
)

func main() {
	rootCmd := &cobra.Command{
		Short:   fmt.Sprintf("IntTest Runtime: mock services master, v%s %s", version, versionDate),
		Version: version + "(" + versionDate + ")",
		Run:     launchServices,
	}

	setRootFlags(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setRootFlags(root *cobra.Command) {
	root.PersistentFlags().String(string(configPathArg), "", "json configuration to launch mock services")
}
