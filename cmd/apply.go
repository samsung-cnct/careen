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
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Computes the hash of file named patchPath and compares it with the expected hash
func VerifyPatch(patch string, expectedHash string) (valid bool, err error) {
	fileData, err := ioutil.ReadFile(patch)
	if err != nil {
		return false, err
	}
	fileLen := len(fileData)

	fileHash := sha1.New()
	io.WriteString(fileHash, string(fileData[:fileLen]))
	computedHash := hex.EncodeToString(fileHash.Sum(nil))
	if computedHash != expectedHash {
		return false, fmt.Errorf("Computed hash %v does not equal expected hash %v", computedHash, expectedHash)
	}

	return true, nil
}

func Apply(repoDir string, patchPath string) (err error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	absRepoDir, err := filepath.Abs(repoDir)
	if err != nil {
		return err
	}

	absPatchPath, err := filepath.Abs(patchPath)
	if err != nil {
		return err
	}

	oldPwd, err := os.Getwd()
	if err != nil {
		return err
	}

	defer func() {
		err = os.Chdir(oldPwd)
	}()

	err = os.Chdir(absRepoDir)
	if err != nil {
		return err
	}

	cmdName := "git"
	cmdArgs := []string{"apply", absPatchPath}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Printf("INFO: Running command \"%v %v\"\n", cmdName, strings.Join(cmdArgs, " "))
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		fmt.Fprintf(os.Stderr, "Stdout:\n %v", stdout.String())
		fmt.Fprintf(os.Stderr, "Stderr:\n %v", stderr.String())
		return fmt.Errorf("There was an error running \"git apply\" command")
	}
	fmt.Println("INFO: Command completed successfully\n")

	return nil
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:          "apply [config filename] (default ) " + careenConfig.GetString("config"),
	Short:        "Applies patches to repositories",
	SilenceUsage: true,
	Long:         `Applies patches to the repositories after verifying that the patch file matches the specified hash`,
	Run: func(cmd *cobra.Command, args []string) {
		manifestFilename := careenConfig.GetString("manifest")
		fmt.Printf("INFO: Using manifest %v\n", manifestFilename)

		manifest, err := GetManifestFromFile(manifestFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			fmt.Fprintf(os.Stderr, "ERROR: Failed to get manifest %v\n", manifestFilename)
			ExitCode = 1
			return
		}

		patchDir := careenConfig.GetString("patches.directory")
		outputDir := careenConfig.GetString("output.directory")

		for _, pkg := range manifest.Packages {
			fmt.Printf("INFO: Applying patches to package: %v\n", pkg.Name)
			repoDir := outputDir + pkg.Name
			for _, patch := range pkg.Patches {
				patchName := patchDir + patch.Filename
				fmt.Printf("INFO: Applying patch %v to repo %v\n", patchName, repoDir)
				valid, err := VerifyPatch(patchName, patch.Hash)
				if !valid || err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					fmt.Fprintf(os.Stderr, "ERROR: Refusing to apply patch %v\n", patchName)
					ExitCode = 1
					return
				}
				err = Apply(repoDir, patchName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
					fmt.Fprintf(os.Stderr, "ERROR: Failed to apply patch %v\n", patchName)
					ExitCode = 1
					return
				}
				fmt.Printf("INFO: Applied patch %v to repo %v\n", patchName, repoDir)
			}
		}

		ExitCode = 0
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
