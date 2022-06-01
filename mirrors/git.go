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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
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

var defaultPushRefSpecs = []config.RefSpec{
	"refs/heads/*:refs/heads/*",
	"refs/remotes/*:refs/remotes/*",
	"refs/tags/*:refs/tags/*"}

type GitClient struct {
	auth         transport.AuthMethod
	cloneOptions *git.CloneOptions
	pullOptions  *git.PullOptions
	fetchOptions *git.FetchOptions
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
		auth: publicKey,
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
		fetchOptions: &git.FetchOptions{
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
		auth: auth,
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
		fetchOptions: &git.FetchOptions{
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
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) == true {
			return nil
		}
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return err
		}
		return fmt.Errorf("pull err: %s", err.Error())
	}
	return nil
}

// CloneOrPull if path is not exist run git clone, else pull
func (c *GitClient) CloneOrPull(url, remoteName, path string) (bool, error) {
	if remoteName == "" {
		remoteName = "origin"
	}

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

// Fetch git repo changes to local directory
func (c *GitClient) Fetch(remoteName, path string) error {
	if remoteName == "" {
		remoteName = "origin"
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("open git repository from path %s err: %s", path, err.Error())
	}

	// Pull the latest changes from the remoteName remote and merge into the current branch
	logger.Infof("[git fetch %s] in path %s", remoteName, path)
	o := *c.fetchOptions
	o.RemoteName = remoteName
	o.RefSpecs = []config.RefSpec{"refs/*:refs/*"}
	o.Tags = git.TagFollowing
	//err = w.Pull(&o)
	// pull with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	go func() {
		<-time.After(c.Timeout)
		cancel()
	}()
	err = r.FetchContext(ctx, &o)
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) == true {
			return nil
		}
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return err
		}
		return fmt.Errorf("[git fetch %s] in path %s err: %s", remoteName, path, err.Error())
	}
	return nil
}

// CloneOrFetch if path is not exist run git clone, else fetch
func (c *GitClient) CloneOrFetch(url, remoteName, path string) (bool, error) {
	if remoteName == "" {
		remoteName = "origin"
	}

	_, err := git.PlainOpen(path)
	if err == nil {
		err = c.CreateRemote([]string{url}, remoteName, path)
		if err != nil {
			return false, err
		}
		return false, c.Fetch(remoteName, path)
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

// findRemoteBranchesAndTag
func (c *GitClient) findRemoteBranchesAndTag(repo *git.Repository, remoteName string) (map[string]string, map[string]string, error) {
	//repoTags := make(map[string]string)
	//repoBranches := make(map[string]string)
	//// all Tags
	//tags, err := repo.Tags()
	//if err != nil {
	//	return err
	//}
	//_ = tags.ForEach(func(tag *plumbing.Reference) error {
	//	repoTags[tag.Name().String()] = tag.Hash().String()
	//	return nil
	//})
	//
	//// all Branches
	//branches, err := repo.Branches()
	//if err != nil {
	//	return err
	//}
	//_ = branches.ForEach(func(branch *plumbing.Reference) error {
	//	repoBranches[branch.Name().String()] = branch.Hash().String()
	//	return nil
	//})
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return nil, nil, err
	}
	if remote == nil {
		return nil, nil, fmt.Errorf("can not find remote %s", remoteName)
	}

	// List the references on the remote repository
	refs, err := remote.List(&git.ListOptions{
		Auth:            c.auth,
		InsecureSkipTLS: false,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("list remote %s the references on the remote repository err: %s", remoteName, err.Error())
	}

	tags := make(map[string]string)
	branches := make(map[string]string)
	for _, remoteRef := range refs {
		// skip InvalidReference and SymbolicReference
		if remoteRef.Type() == plumbing.HashReference {
			if remoteRef.Name().IsBranch() {
				branches[remoteRef.Name().String()] = remoteRef.Hash().String()
			} else if remoteRef.Name().IsTag() {
				tags[remoteRef.Name().String()] = remoteRef.Hash().String()
			}
		}
	}

	return branches, tags, nil
}

// fixPrune fix Push with Prune does not achieve the desired effect
// ref: https://github.com/go-git/go-git/issues/172 bug
func (c *GitClient) fixPrune(repo *git.Repository, srcRemoteName, dstRemoteName, path string) error {
	// Src Remote: get remote by srcRemoteName
	srcRemoteBranches, srcRemoteTags, err := c.findRemoteBranchesAndTag(repo, srcRemoteName)
	if err != nil {
		return err
	}

	// Src Remote: get remote by dstRemoteName
	if dstRemoteName == "" {
		dstRemoteName = "origin"
	}
	dstRemoteBranches, dstRemoteTags, err := c.findRemoteBranchesAndTag(repo, dstRemoteName)
	if err != nil {
		return err
	}

	var delRefs []config.RefSpec
	// find which branch to del
	for name, dstHash := range dstRemoteBranches {
		srcHash, ok := srcRemoteBranches[name]
		if ok == true && dstHash == srcHash {
			continue
		}

		delRefs = append(delRefs, config.RefSpec(fmt.Sprintf(":%s", name)))
	}

	// find which tags to del
	for name, dstHash := range dstRemoteTags {
		srcHash, ok := srcRemoteTags[name]
		if ok == true && dstHash == srcHash {
			continue
		}

		delRefs = append(delRefs, config.RefSpec(fmt.Sprintf(":%s", name)))
	}

	// push diffRefs
	if len(delRefs) > 0 {
		remote, err := repo.Remote(dstRemoteName)
		if err != nil {
			return err
		}
		if remote == nil {
			return fmt.Errorf("can not find remote %s in path %s", dstRemoteName, path)
		}

		o := *c.pushOptions
		o.RemoteName = dstRemoteName
		o.RefSpecs = delRefs

		// push
		if err := remote.Push(&o); err != nil {
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				return nil
			} else {
				return fmt.Errorf("fixPrune %s in path %s occur err: %s", dstRemoteName, path, err.Error())
			}
		}

		// fetch dst remote
		if err := c.Fetch(dstRemoteName, path); err != nil {
			logger.Errorf("fixPrune %s in path %s occur err: %s", dstRemoteName, path, err.Error())
			return err
		}
	}

	return nil
}

// Push open a repository in a specific path, and push to its remoteName remote.
// equal git cmd:
//   git push --prune --tags [--force] [origin|gitee|github] "refs/*:refs/*" #refs/remotes/origin/*:refs/heads/*
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
	// run:
	//   git show-ref
	//   git remote -v
	// Push with Prune does not achieve the desired effect, ref: https://github.com/go-git/go-git/issues/172
	o.RefSpecs = defaultPushRefSpecs
	//o.RefSpecs = []config.RefSpec{"+refs/heads/*:refs/remotes/origin/*"}
	o.Prune = false
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
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logger.Debugf("push remoteName %s. path: %s, already up-to-date", remoteName, path)
		} else {
			return fmt.Errorf("push remoteName: %s, path: %s, err: %s", remoteName, path, err.Error())
		}
	}

	// in https://github.com/go-git/go-git/blob/v5.4.2/COMPATIBILITY.md prune in not support in v5.4.2
	err = c.fixPrune(r, "origin", remoteName, path)
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logger.Debugf("fix prune push remoteName %s. path: %s, already up-to-date", remoteName, path)
		} else {
			return err
		}
	}

	return nil
}
