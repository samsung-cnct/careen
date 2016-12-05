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
	"github.com/spf13/cobra"
	"os"
)

func Clone(repoUrl string, revision string, destDir string) error {
	err := GitClone(repoUrl, revision, destDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return fmt.Errorf("Failed to clone repo %v to destination %v", repoUrl, destDir)
	}

	return nil
}

func CheckoutByTag(repoDir string, tag string) error {
	repo, err := GitOpenRepository(repoDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return fmt.Errorf("Failed to open repo at %v", repoDir)
	}

	err = GitCheckoutByTag(repo, tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return fmt.Errorf("Failed to check out tag %v from repo %v", tag, repoDir)
	}

	return nil
}

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:          "clone [configuration filename] (default ) " + careenConfig.GetString("careenConfigFile"),
	Short:        "Clones repositories",
	SilenceUsage: true,
	Long:         `Clones repositories at a specific commit specified by configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			manifestFilename string
			manifest         *Manifest
			outputDir        string
			err              error
		)

		manifestFilename = careenConfig.GetString("manifest")
		fmt.Printf("INFO: Cloning packages from manifest %v\n", manifestFilename)

		manifest, err = GetManifestFromFile(manifestFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			ExitCode = 1
			return
		}

		outputDir = careenConfig.GetString("output")

		for _, pkg := range manifest.Packages {
			repoDir := outputDir + pkg.Name
			fmt.Printf("INFO: Cloning package %v from %v to %v\n", pkg.Name, pkg.Repo, repoDir)
			terminalSpinner.Start()
			err = Clone(pkg.Repo, pkg.Revision, repoDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				ExitCode = 1
				return
			}
			CheckoutByTag(repoDir, pkg.Tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				ExitCode = 1
				return
			}
			terminalSpinner.Stop()
		}

		ExitCode = 0
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
