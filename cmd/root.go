/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rspace-client",
	Short: "RSpace CLI",
	Long: `CLI for RSpace - make API calls to RSpace
To get started, set your API key and RSpace URL in file '.rspace' in your home folder, e.g.

RSPACE_API_KEY=fsdfsd
RSPACE_URL=https://myrspace.org/api/v1
	
Alternatively set these as environment variables.

To see all the ELN commands run rspace eln --help
`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// NewOperatorCommand returns the `quarks-job` command.
func NewOperatorCommand() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// specify config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rspace)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			exitWithErr(err)
		}
		// Search for config in home directory  ".rspace" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rspace")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		messageStdErr("Using config file:" + viper.ConfigFileUsed())
	}
	if err != nil {
		msg := fmt.Sprintf("No config file found. Looking for environment variables %s and %s",
			APIKEY_ENV_NAME, BASE_URL_ENV_NAME)
		messageStdErr(msg)

		if viper.Get(APIKEY_ENV_NAME) == nil || viper.Get(BASE_URL_ENV_NAME) == nil {
			exitWithStdErrMsg("No config file or environment variables found")
		}
	}
}
