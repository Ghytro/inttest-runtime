package main

import (
	"fmt"
	"inttest-runtime/internal/config"
	"inttest-runtime/internal/errors/internalErr"
	"log"

	"github.com/pkg/errors"
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

func launchServices(cmd *cobra.Command, args []string) {
	confPath, err := cmd.PersistentFlags().GetString(string(configPathArg))
	if err != nil {
		err := errors.Wrap(err, "error getting config file path")
		log.Fatal(err)
	}
	config, err := config.FromFile(confPath)
	if err != nil {
		err := internalErr.WrapWithCode(err, internalErr.ErrCodeConfigurationParsing)
		log.Fatal(err)
	}

}
