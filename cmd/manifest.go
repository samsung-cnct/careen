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
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Manifest struct {
	Version  string
	Packages []Package
}

type Package struct {
	Name     string
	Repo     string
	Revision string
	Patches  []Patch
}

type Patch struct {
	Name          string
	Filename      string
	Hash          string
	Documentation []string
}

func GetManifestFromFile(filename string) (*Manifest, error) {
	manifest := Manifest{}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("manifest: error reading %v: %v", filename, err)
	}

	err = yaml.Unmarshal([]byte(file), &manifest)
	if err != nil {
		return nil, fmt.Errorf("manifest: error parsing %v: %v", filename, err)
	}

	return &manifest, nil
}
