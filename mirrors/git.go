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
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/xiexianbin/golib/logger"
)

type GitAuthType int

const (
	GitKeyAuth GitAuthType = iota
	GitAccessTokenAuth
	GitUsernamePasswordAuth
)

type GitClient struct {
	cloneOptions *git.CloneOptions
	pullOptions  *git.PullOptions
	pushOptions  *git.PushOptions
	Timeout      time.Duration
	GitAuthType  GitAuthType
}

// NewGitPrivateKeysClient ssh key auth
func NewGitPrivateKeysClient(privateKeyFile, keyPassword string, timeout time.Duration) (*GitClient, error) {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("read file %s failed %s", privateKeyFile, err.Error())
	}

	// Username must be "git" for SSH auth to work, not your real username.
	// See https://github.com/src-d/go-git/issues/637
	//publicKey, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, keyPassword)
	//if err != nil {
	//	return nil, fmt.Errorf("read private key failed: %s", err.Error())
	//}
	sshKey, _ := ioutil.ReadFile(privateKeyFile)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), keyPassword)
	if err != nil {
		return nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return &GitClient{
		cloneOptions: &git.CloneOptions{
			Auth:            publicKey,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		pullOptions: &git.PullOptions{
			Auth:            publicKey,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		pushOptions: &git.PushOptions{
			Auth:            publicKey,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		Timeout:     timeout,
		GitAuthType: GitKeyAuth,
	}, nil
}

// httpBasicAuthClient
func httpBasicAuthClient(username, password string, timeout time.Duration, authType GitAuthType) (*GitClient, error) {
	var auth *http.BasicAuth
	switch authType {
	case GitAccessTokenAuth:
		if username == "" {
			username = "xiexianbin"
		}
		if password != "" {
			// The intended use of a GitHub personal access token is in replace of your password
			// because access tokens can easily be revoked.
			// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
			auth = &http.BasicAuth{
				Username: username, // yes, this can be anything except an empty string
				Password: password,
			}
		}
	case GitUsernamePasswordAuth:
		if username != "" && password != "" {
			auth = &http.BasicAuth{
				Username: username,
				Password: password,
			}
		}
	}

	if auth == nil {
		return nil, fmt.Errorf("invalid authentication")
	}

	return &GitClient{
		cloneOptions: &git.CloneOptions{
			Auth:            auth,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		pullOptions: &git.PullOptions{
			Auth:            auth,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		pushOptions: &git.PushOptions{
			Auth:            auth,
			Progress:        os.Stdout,
			InsecureSkipTLS: true,
		},
		Timeout:     timeout,
		GitAuthType: authType,
	}, nil
}

// NewGitAccessTokenClient access_token auth
func NewGitAccessTokenClient(accessToken string, timeout time.Duration) (*GitClient, error) {
	return httpBasicAuthClient("", accessToken, timeout, GitAccessTokenAuth)
}

// NewGitUsernamePasswordClient username password auth
func NewGitUsernamePasswordClient(username, password string, timeout time.Duration) (*GitClient, error) {
	return httpBasicAuthClient(username, password, timeout, GitUsernamePasswordAuth)
}

// Clone clone git to local directory
func (c *GitClient) Clone(url, path string) error {
	// Clone the given repository to the given path
	logger.Infof("[git clone %s] in path %s", url, path)
	o := *c.cloneOptions
	o.URL = url
	//_, err := git.PlainClone(path, false, &o)
	// clone with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	go func() {
		<-time.After(c.Timeout)
		cancel()
	}()
	_, err := git.PlainCloneContext(ctx, path, false, &o)
	if err != nil {
		// if is "remote repository is empty" err, skip
		//if errors.Is(err, transport.ErrEmptyRemoteRepository) {
		//	return nil
		//}
		return err
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
	o := *c.pullOptions
	o.RemoteName = remoteName
	//err = w.Pull(&o)
	// pull with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	go func() {
		<-time.After(c.Timeout)
		cancel()
	}()
	err = w.PullContext(ctx, &o)
	if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) == false {
		return fmt.Errorf("pull err: %s", err.Error())
	}
	return nil
}

// CloneOrPull if path is not git clone it, else pull
func (c *GitClient) CloneOrPull(url, remoteName, path string) (bool, error) {
	_, err := git.PlainOpen(path)
	if err == nil {
		err = c.CreateRemote([]string{url}, remoteName, path)
		if err != nil {
			return false, err
		}
		return false, c.Pull(remoteName, path)
	} else {
		return true, c.Clone(url, path)
	}
}

// DeleteRemote delete special remote
func (c *GitClient) DeleteRemote(remoteName, path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("delete remote %s, when open git repository from path %s err: %s", remoteName, path, err.Error())
	}

	return r.DeleteRemote(remoteName)
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

	remotesList := make([]string, len(remotes))
	for i, remote := range remotes {
		remotesList[i] = fmt.Sprintf("%s %s", remote.Config().Name, remote.Config().URLs[0])
	}

	// List remotes from a repository
	logger.Debugf("[git remotes -v] in path %s, remotes is %s", path, strings.Join(remotesList, ", "))
	for _, remote := range remotes {
		if remote.Config().Name == remoteName {
			logger.Debugf("remote %s in path %s is already created, url is %s", remoteName, path, remote.Config().URLs[0])
			// check url, remote url is match
			if remote.Config().URLs[0] == urls[0] {
				return nil
			}

			// remote url is not match
			err := c.DeleteRemote(remoteName, path)
			if err != nil {
				return err
			}
		}
	}

	// Add a new remote, with the default fetch refspec
	logger.Infof("[git remote add %s %s]", remoteName, urls[0])
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
		return fmt.Errorf("when open git repository from path %s err: %s", path, err.Error())
	}

	if force == false {
		logger.Infof("[git push %s] in path %s", remoteName, path)
	} else {
		logger.Warnf("[git push %s -f] in path %s", remoteName, path)
	}
	o := *c.pushOptions
	o.RemoteName = remoteName
	// run: git show-ref
	// https://github.com/go-git/go-git/issues/172
	o.RefSpecs = []config.RefSpec{"refs/remotes/origin/*:refs/heads/*"}
	//o.RefSpecs = []config.RefSpec{"refs/*:refs/*"}
	o.Prune = true
	o.Force = force

	//err = r.Push(&o)
	// push with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	go func() {
		<-time.After(c.Timeout)
		cancel()
	}()
	err = r.PushContext(ctx, &o)
	if err != nil {
		return fmt.Errorf("push remoteName: %s, path: %s, err: %s", remoteName, path, err.Error())
	}

	return nil
}
