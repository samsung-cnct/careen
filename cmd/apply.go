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
	"time"
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

// Run command with args and kill if timeout is reached
func RunCommand(name string, args []string, timeout time.Duration) error {
	fmt.Printf("Running command \"%v %v\"\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		fmt.Fprintf(os.Stderr, "Command %v timed out\n", name)
		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", name, err))
		}
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", name, err)
			output, e := cmd.CombinedOutput()
			if e != nil {
				return e
			}
			fmt.Fprintf(os.Stderr, "%v", output)
			return err
		}
	}
	fmt.Printf("Command %v completed successfully\n", name)

	return nil
}

// Apply patch to repo in repoDir
func Apply(repoDir string, patchPath string) (err error) {
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
	cmdTimeout := time.Duration(10) * time.Second
	err = RunCommand(cmdName, cmdArgs, cmdTimeout)
	if err != nil {
		return err
	}

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
}
