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
	"io"
	"os"
)

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func Clone(repoUrl string, revision string, destDir string) error {
	terminalSpinner.Start()
	defer func() {
		terminalSpinner.Stop()
	}()

	err := GitClone(repoUrl, revision, destDir)
	if err != nil {
		return err
	}

	return nil
}

func CheckoutByTag(repoDir string, tag string) error {
	repo, err := GitOpenRepository(repoDir)
	if err != nil {
		return err
	}

	err = GitCheckoutByTag(repo, tag)
	if err != nil {
		return err
	}

	return nil
}

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:          "clone [config filename] (default ) " + careenConfig.GetString("config"),
	Short:        "Clones repositories",
	SilenceUsage: true,
	Long:         `Clones repositories at a specific commit specified by configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		manifestFilename = careenConfig.GetString("manifest")
		fmt.Printf("INFO: Cloning packages from manifest %v\n", manifestFilename)

		manifest, err := GetManifestFromFile(manifestFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			ExitCode = 1
			return
		}

		outputDir = careenConfig.GetString("output.directory")

		for _, pkg := range manifest.Packages {
			repoDir := outputDir + pkg.Name
			fmt.Printf("INFO: Checking if repository directory %v is empty\n", repoDir)
			empty, err := IsEmpty()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				ExitCode = 1
				return
			} else if empty {
				fmt.Printf("INFO: Attempting to clone repository %v to directory %v\n", pkg.Repo, repoDir)
				err = Clone(pkg.Repo, pkg.Revision, repoDir)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					ExitCode = 1
					return
				}
			}
			fmt.Printf("INFO: Attempting to checkout tag %v from repository directory %v\n", pkg.Tag, repoDir)
			err = CheckoutByTag(repoDir, pkg.Tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				ExitCode = 1
				return
			}
			fmt.Printf("INFO: Checked out tag %v from repository directory %v\n", pkg.Tag, repoDir)
		}

		ExitCode = 0
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)
}
