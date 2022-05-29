// Copyright 2022 xiexianbin<me@xiexianbin.cn>
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

package mirrors

import (
	"testing"
)

const (
	GithubRepoCloneUrl = "https://github.com/x-actions/git-mirrors"
	GithubRepoSSHURL   = "git@github.com:estack/estack.git"
	GiteeRepoUrl       = "https://gitee.com/x-actions/git-mirrors"
	TempPath           = "../temp/git-mirrors"
)

func TestNewAccessTokenClient(t *testing.T) {
	var err error
	c, err := NewGitAccessTokenClient("", "")
	if err != nil {
		t.Skip(err.Error())
	}

	err = c.Clone(GithubRepoCloneUrl, TempPath)
	if err != nil {
		t.Skip(err.Error())
	}
}

func TestGitClient_CloneOrPull(t *testing.T) {
	var err error
	//c, err := NewGitUsernamePasswordClient("", "")
	c, err := NewGitAccessTokenClient("", "")
	//c, err := NewGitPrivateKeysClient("/Users/xiexianbin/workspace/code/github.com/x-actions/git-mirrors/temp/id_ed25519", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	isNewClone, err := c.CloneOrPull(GithubRepoCloneUrl, "", TempPath)
	//isNewClone, err := c.CloneOrPull(GithubRepoSSHURL, "", TempPath)
	if err != nil {
		t.Skip(err.Error())
	}
	t.Logf("isNewClone: %t", isNewClone)
}

func TestGitClient_CreateRemote(t *testing.T) {
	var err error
	c, err := NewGitUsernamePasswordClient("", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.CreateRemote([]string{GiteeRepoUrl}, "gitee", TempPath)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGitClient_Push(t *testing.T) {
	var err error
	c, err := NewGitUsernamePasswordClient("", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	err = c.Push("gitee", TempPath, false)
	if err != nil {
		t.Skip(err)
	}
}
