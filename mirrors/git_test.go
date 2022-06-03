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
	"time"
)

const (
	//GithubRepoCloneUrl = "https://github.com/x-actions/git-mirrors"
	GithubRepoCloneUrl = "git@github.com:x-actions/git-mirrors.git"
	GiteeRepoUrl       = "https://gitee.com/x-actions/git-mirrors"
	TempPath           = "../temp/git-mirrors"
	privateKeyPath     = "../temp/id_ed25519"
)

var defaultTimeOut = 1 * time.Minute

func TestNewAccessTokenClient(t *testing.T) {
	var err error
	c, err := NewGitAccessTokenClient("", defaultTimeOut, true)
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
	//c, err := NewGitAccessTokenClient("", defaultTimeOut)
	c, err := NewGitPrivateKeysClient(privateKeyPath, "", defaultTimeOut, true)
	if err != nil {
		t.Skip(err.Error())
		return
	}

	isNewClone, err := c.CloneOrPull(GithubRepoCloneUrl, "", TempPath)
	//isNewClone, err := c.CloneOrPull(GithubRepoSSHURL, "", TempPath)
	if err != nil {
		t.Skip(err.Error())
		return
	}
	t.Logf("isNewClone: %t", isNewClone)
}

func TestGitClient_CloneOrFetch(t *testing.T) {
	var err error
	//c, err := NewGitUsernamePasswordClient("", "")
	//c, err := NewGitAccessTokenClient("", defaultTimeOut)
	c, err := NewGitPrivateKeysClient(privateKeyPath, "", defaultTimeOut, true)
	if err != nil {
		t.Skip(err.Error())
		return
	}

	isNewClone, err := c.CloneOrFetch(GithubRepoCloneUrl, "", TempPath)
	//isNewClone, err := c.CloneOrPull(GithubRepoSSHURL, "", TempPath)
	if err != nil {
		t.Skip(err.Error())
		return
	}
	t.Logf("isNewClone: %t", isNewClone)
}

func TestGitClient_CreateRemote(t *testing.T) {
	var err error
	c, err := NewGitUsernamePasswordClient("", "", defaultTimeOut, true)
	if err != nil {
		t.Skip(err.Error())
		return
	}

	err = c.CreateRemote([]string{GiteeRepoUrl}, "gitee", TempPath)
	if err != nil {
		t.Skip(err.Error())
		return
	}
}

func TestGitClient_Mirror(t *testing.T) {
	var err error
	//c, err := NewGitUsernamePasswordClient("", "", defaultTimeOut)
	c, err := NewGitPrivateKeysClient(privateKeyPath, "", defaultTimeOut, true)
	if err != nil {
		t.Skip(err.Error())
		return
	}

	err = c.Mirror("gitee", TempPath, false)
	if err != nil {
		t.Skip(err)
	}
}
