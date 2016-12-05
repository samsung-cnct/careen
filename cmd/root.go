// Copyright © 2016 Samsung CNCT
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"time"
)

var configFilename string
var manifestFilename string
var outputDirectory string
var patchDirectory string
var ExitCode int

// progress spinner
var terminalSpinner = spinner.New(spinner.CharSets[35], 200*time.Millisecond)

// init the careen config viper instance
var careenConfig = viper.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "careen",
	Short: "CLI for patching github repos",
	Long:  `careen is a command line interface for cloning and patching a set of github repos`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
			os.Mkdir(outputDirectory, 0755)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initCareenConfig)

	RootCmd.SetHelpCommand(helpCmd)

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defaultConfigFilename := "$HOME/.carren.yaml"
	defaultManifestFilename := workingDir + "/manifests/docker.yaml"
	defaultOutputDirectory := workingDir + "/src/"
	defaultPatchDirectory := workingDir + "/patches/"

	RootCmd.PersistentFlags().StringVarP(
		&configFilename,
		"config",
		"c",
		defaultConfigFilename,
		"config file (default \""+defaultConfigFilename+"\"")
	RootCmd.PersistentFlags().StringVarP(
		&manifestFilename,
		"manifest",
		"m",
		defaultManifestFilename,
		"manifest filename")
	RootCmd.PersistentFlags().StringVarP(
		&outputDirectory,
		"output",
		"o",
		defaultOutputDirectory,
		"output directory")
	RootCmd.PersistentFlags().StringVarP(
		&patchDirectory,
		"patches",
		"p",
		defaultPatchDirectory,
		"patch directory")

	configureSpinner(terminalSpinner)

}

// initCareenConfig uses flags (with defaults set by init), ENV variables and finally configuration files
func initCareenConfig() {
	careenConfig.BindPFlag("manifest", RootCmd.Flags().Lookup("manifest"))
	careenConfig.BindPFlag("output", RootCmd.Flags().Lookup("output"))
	careenConfig.BindPFlag("patches", RootCmd.Flags().Lookup("patches"))

	careenConfig.SetEnvPrefix("CAREEN_")
	careenConfig.AutomaticEnv()

	if configFilename != "" { // enable ability to specify config file via flag
		careenConfig.SetConfigName("careen")         // name of config file (without extension)
		careenConfig.AddConfigPath("$HOME/.careen/") // path to look for the config file in
		careenConfig.AddConfigPath(".")              // optionally look for config in the working directory
	} else {
		// Warn the user if they explicitly requested a config file which does not exist
		if _, err := os.Stat(configFilename); os.IsNotExist(err) {
			fmt.Printf("WARNING: Specified config file %v does not exist, using defaults", configFilename)
		} else {
			careenConfig.SetConfigFile(configFilename)
		}
	}
	// Ignore errors. All configuration parameters have defaults.
	careenConfig.ReadInConfig()
}

func configureSpinner(s *spinner.Spinner) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		// Directing the spinner to stderr makes the command compatible with pipe, et al.
		s.Writer = os.Stderr
	} else {
		// Directing the spinner to /dev/null makes the command play nice when run non-interactively (e.g. by Jenkins)
		s.Writer = ioutil.Discard
	}

	s.FinalMSG = "\n"
}
