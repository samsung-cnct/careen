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

/* For instructions on installing git2go see:
   http://www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/
*/

import (
	"github.com/libgit2/git2go"
)

func GitClone(repoUrl string, revision string, destDir string) error {
	_, err := git.Clone(repoUrl, destDir, &git.CloneOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GitOpenRepository(repoPath string) (*git.Repository, error) {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func GitCheckoutByTag(repo *git.Repository, tag string) error {
	/* Note DwimReference() was renamed Dwim() and moved from repository.go
	   to references.go. This can lead to some confusion if you read the
	   docs and not the code. Beware.
	*/
	ref, err := repo.References.Dwim(tag)
	if err != nil {
		return err
	}

	if err := repo.SetHeadDetached(ref.Target()); err != nil {
		return err
	}

	return nil
}
