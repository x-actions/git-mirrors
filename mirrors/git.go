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
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/xiexianbin/golib/logger"
)

const (
	PublicKeysAuth int = iota
	AccessTokenAuth
	UsernamePassword
)

type GitClient struct {
	options  *git.CloneOptions
	authType int
}

// NewGitPublicKeysClient ssh key auth
func NewGitPublicKeysClient(privateKeyFile, keyPassword string) (*GitClient, error) {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("read file %s failed %s", privateKeyFile, err.Error())
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, keyPassword)
	if err != nil {
		return nil, fmt.Errorf("generate publickeys failed: %s", err.Error())
	}

	o := &git.CloneOptions{
		Auth:     publicKeys,
		Progress: os.Stdout,
	}

	return &GitClient{options: o, authType: PublicKeysAuth}, nil
}

// NewGitAccessTokenClient access_token auth
func NewGitAccessTokenClient(accessToken string) (*GitClient, error) {
	o := &git.CloneOptions{
		Progress: os.Stdout,
	}
	if accessToken != "" {
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		o.Auth = &http.BasicAuth{
			Username: "xiexianbin", // yes, this can be anything except an empty string
			Password: accessToken,
		}
	}
	return &GitClient{options: o, authType: AccessTokenAuth}, nil
}

// NewGitUsernamePasswordClient username password auth
func NewGitUsernamePasswordClient(username, password string) (*GitClient, error) {
	o := &git.CloneOptions{
		Progress: os.Stdout,
	}
	if username != "" && password != "" {
		o.Auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	return &GitClient{options: o, authType: UsernamePassword}, nil
}

// Clone clone git to local directory
func (c *GitClient) Clone(url, path string) error {
	// Clone the given repository to the given path
	logger.Infof("[git clone %s] in path %s", url, path)
	c.options.URL = url
	_, err := git.PlainClone(path, false, c.options)
	if err != nil {
		return fmt.Errorf("clone %s err: %s", url, err.Error())
	}
	return nil
}

// Pull git repo changes to local directory
func (c *GitClient) Pull(remoteName, path string) error {
	if remoteName == "" {
		remoteName = "origin"
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("open git repository from path %s err: %s", path, err.Error())
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree err: %s", err.Error())
	}

	// Pull the latest changes from the remoteName remote and merge into the current branch
	logger.Infof("[git pull %s] in path %s", remoteName, path)
	err = w.Pull(&git.PullOptions{RemoteName: remoteName})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("pull err: %s", err.Error())
	}
	return nil
}

// CloneOrPull if path is not git clone it, else pull
func (c *GitClient) CloneOrPull(url, remoteName, path string) (bool, error) {
	_, err := git.PlainOpen(path)
	if err == nil {
		return false, c.Pull(remoteName, path)
	} else {
		return true, c.Clone(url, path)
	}
}

// CreateRemote create remote
// url eg. https://github.com/git-fixtures/basic.git
func (c *GitClient) CreateRemote(urls []string, remoteName, path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("create remote, when open git repository from path %s err: %s", path, err.Error())
	}

	// list all remote
	remotes, err := r.Remotes()
	if err != nil {
		return fmt.Errorf("get %s remote err: %s", path, err.Error())
	}

	// List remotes from a repository
	logger.Infof("[git remotes -v] in path %s, remotes len is %d", path, len(remotes))
	for _, remote := range remotes {
		if strings.Split(remote.String(), "\t")[0] == remoteName {
			logger.Warnf("remote %s in path %s is already created", remoteName, path)
			// may be need to check url
			return nil
		}
	}

	// Add a new remote, with the default fetch refspec
	logger.Infof("[git remote add %s %s]", remoteName, urls)
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: urls,
	})
	if err != nil {
		return fmt.Errorf("create remote for %s err: %s", path, err.Error())
	}

	return nil
}

// Push open a repository in a specific path, and push to its remoteName remote.
func (c *GitClient) Push(remoteName, path string, force bool) error {
	if remoteName == "" {
		remoteName = "origin"
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("create remote, when open git repository from path %s err: %s", path, err.Error())
	}

	o := &git.PushOptions{
		RemoteName: remoteName,
		Progress:   os.Stdout,
		Force:      force,
	}
	err = r.Push(o)
	if err != nil {
		return fmt.Errorf("push remoteName: %s, path: %s, err: %s", remoteName, path, err.Error())
	}

	return nil
}
