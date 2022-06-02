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
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/xiexianbin/golib/logger"

	"github.com/x-actions/git-mirrors/constants"
)

type Mirror struct {
	SrcGit         string
	SrcOrg         string
	srcToken       string
	DstGit         string
	DstOrg         string
	dstKey         string
	dstToken       string
	SrcAccountType string
	DstAccountType string
	CloneStyle     string
	CachePath      string
	BlackList      []string
	WhiteList      []string
	ForceUpdate    bool
	Timeout        time.Duration
	Mappings       map[string]string

	blackListMap map[string]string
	whiteListMap map[string]string

	srcRepos    []*Repository
	srcReposMap map[string]*Repository
	dstRepos    []*Repository
	dstReposMap map[string]*Repository

	srcGitClient *GitClient
	dstGitClient *GitClient
	srcAPI       interface{}
	dstAPI       interface{}
}

func New(srcGit, srcOrg, srcToken, dstGit, dstOrg, dstKey, dstToken, srcAccountType, dstAccountType, cloneStyle,
	cachePath string, blackList, whiteList []string, forceUpdate bool, timeout time.Duration,
	mappings map[string]string) *Mirror {
	return &Mirror{
		SrcGit:         srcGit,
		SrcOrg:         srcOrg,
		srcToken:       srcToken,
		DstGit:         dstGit,
		DstOrg:         dstOrg,
		dstKey:         dstKey,
		dstToken:       dstToken,
		SrcAccountType: srcAccountType,
		DstAccountType: dstAccountType,
		CloneStyle:     cloneStyle,
		CachePath:      cachePath,
		BlackList:      RemoveDuplicates(blackList),
		blackListMap:   StringListToMap(blackList),
		WhiteList:      RemoveDuplicates(whiteList),
		whiteListMap:   StringListToMap(whiteList),
		ForceUpdate:    forceUpdate,
		Timeout:        timeout,
		Mappings:       mappings,
	}
}

// prepare init src/dst APIs and Repos
func (m *Mirror) prepare() error {
	initAPI := func(t, accessToken string) (IGitAPI, error) {
		switch t {
		// init Github api Client
		case constants.GITHUB:
			logger.Infof("init %s API use accessToken(len: %d)", constants.GITHUB, len(accessToken))
			client, err := NewGithubAPI(accessToken)
			if err != nil {
				return nil, err
			}
			return client, nil

		// init Gitee api Client
		case constants.GITEE:
			logger.Infof("init %s API use accessToken(len: %d)", constants.GITEE, len(accessToken))
			client, err := NewGiteeAPI(accessToken)
			if err != nil {
				return nil, err
			}
			return client, nil

		default:
			return nil, fmt.Errorf("un-support git %s", m.SrcGit)
		}
	}

	getRepos := func(t, orgName string, client IGitAPI) ([]*Repository, error) {
		switch t {
		// init User Repos
		case constants.AccountTypeUser:
			repos, err := client.Repositories(orgName)
			if err != nil {
				return nil, err
			}
			return repos, nil

		// init Org Repos
		case constants.AccountTypeOrg:
			repos, err := client.RepositoriesByOrg(orgName)
			if err != nil {
				return nil, err
			}
			return repos, nil

		default:
			return nil, fmt.Errorf("un-support account-type %s", t)
		}
	}

	initGitClient := func(keyPath, accessToken string) (*GitClient, error) {
		if keyPath != "" {
			logger.Infof("use ssh private key to init git client")
			// maybe need to support ssh key with password
			return NewGitPrivateKeysClient(keyPath, "", m.Timeout)
		} else if accessToken != "" {
			logger.Infof("use accessToken to init git client")
			return NewGitAccessTokenClient(accessToken, m.Timeout)
		}

		return nil, fmt.Errorf("git client init fail")
	}

	// init src
	srcAPI, err := initAPI(m.SrcGit, m.srcToken)
	if err != nil {
		return err
	}
	m.srcAPI = srcAPI

	srcRepos, err := getRepos(m.SrcAccountType, m.SrcOrg, srcAPI)
	if err != nil {
		return err
	}
	m.srcRepos = srcRepos
	m.srcReposMap = ReposToMap(srcRepos)

	srcGitClient, err := initGitClient(m.dstKey, m.srcToken)
	if err != nil {
		return err
	}
	m.srcGitClient = srcGitClient

	// init dst
	dstAPI, err := initAPI(m.DstGit, m.dstToken)
	if err != nil {
		return err
	}
	m.dstAPI = dstAPI

	dstRepos, err := getRepos(m.DstAccountType, m.DstOrg, dstAPI)
	if err != nil {
		return err
	}
	m.dstRepos = dstRepos
	m.dstReposMap = ReposToMap(dstRepos)

	dstGitClient, err := initGitClient(m.dstKey, m.dstToken)
	if err != nil {
		return err
	}
	m.dstGitClient = dstGitClient

	return nil
}

// isMirrorRepo check is mirror repo
func (m *Mirror) isMirrorRepo(repoName string) bool {
	if repoName == "" {
		return false
	}

	if len(m.BlackList) > 0 {
		if _, ok := m.blackListMap[repoName]; ok {
			return false
		}
	}

	return true
}

func (m *Mirror) getDstRepoName(repoName string) string {
	if name, ok := m.Mappings[repoName]; ok {
		return name
	}
	return repoName
}

// mirrorRepoInfo create or sync Repo Info
func (m *Mirror) mirrorRepoInfo(srcRepo *Repository, dstRepoName string) (*Repository, error) {
	var dstRepo *Repository
	dstRepo, ok := m.dstReposMap[dstRepoName]
	if ok {
		// already created || dstRepo.Private != srcRepo.Private
		if StringsEqual(dstRepo.Homepage, srcRepo.Homepage) == false || len(dstRepo.Topics) != len(srcRepo.Topics) ||
			StringsEqual(dstRepo.Description, srcRepo.Description) == false {
			if client, ok := m.dstAPI.(IGitAPI); ok {
				dstRepo.Homepage = srcRepo.Homepage
				dstRepo.Description = srcRepo.Description
				dstRepo.Topics = srcRepo.Topics
				dstRepo.Private = srcRepo.Private
				_, err := client.UpdateRepository(*dstRepo.Organization.Name, *dstRepo.Name, dstRepo)
				if err != nil {
					logger.Warnf("update repo %s/%s err: %s", *dstRepo.Owner.Name, *dstRepo.Name, err.Error())
					return dstRepo, nil
				}
			} else {
				return nil, fmt.Errorf("git dstAPI is not implement interface IGitAPI.UpdateRepository")
			}
		}
	} else {
		// new create
		if client, ok := m.dstAPI.(IGitAPI); ok {
			dstRepo = &Repository{
				Name:        &dstRepoName,
				Homepage:    srcRepo.Homepage,
				Description: srcRepo.Description,
				Topics:      srcRepo.Topics,
				Private:     srcRepo.Private,
			}
			return client.CreateRepository(dstRepo, m.DstOrg)
		} else {
			return nil, fmt.Errorf("git dstAPI is not implement interface IGitAPI.CreateRepository")
		}
	}

	return dstRepo, nil
}

// mirrorGit clone/pull from src repo and push to dst repo
func (m *Mirror) mirrorGit(srcRepo, dstRepo *Repository) error {
	var err error
	cachePath := path.Join(m.CachePath, *srcRepo.Name)
	// clone or fetch from origin
	_, err = m.srcGitClient.CloneOrFetch(GitURL(srcRepo, m.srcGitClient.GitAuthType), "origin", cachePath)
	if err != nil {
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			logger.Warnf("source remote repository %s/%s is empty, skip.", *srcRepo.Owner.Name, *srcRepo.Name)
			return nil
		}
		return err
	}

	// create dst git remote
	err = m.dstGitClient.CreateRemote([]string{GitURL(dstRepo, m.dstGitClient.GitAuthType)}, m.DstGit, cachePath)
	if err != nil {
		return err
	}

	// push to dst
	err = m.dstGitClient.Mirror(m.DstGit, cachePath, m.ForceUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mirror) mirror(srcRepo *Repository, dstRepoName string) error {
	var err error

	// mirror repo infos
	dstRepo, err := m.mirrorRepoInfo(srcRepo, dstRepoName)
	if err != nil {
		return err
	}

	// mirror git commits
	err = m.mirrorGit(srcRepo, dstRepo)
	if err != nil {
		return err
	}

	return nil
}

// Do mirror logic
func (m *Mirror) Do() error {
	var err error
	// get src/dst Repos
	err = m.prepare()
	if err != nil {
		return err
	}

	if len(m.WhiteList) > 0 {
		// mirror white list repos
		total := len(m.WhiteList)
		for i, srcRepoName := range m.WhiteList {
			if srcRepo, ok := m.srcReposMap[srcRepoName]; ok {
				dstRepoName := m.getDstRepoName(srcRepoName)
				logger.Debugf("(%d/%d) begin mirror WhiteList %s/%s/%s to %s/%s/%s",
					i+1, total, m.SrcGit, m.SrcOrg, srcRepoName, m.DstGit, m.DstOrg, dstRepoName)
				err := m.mirror(srcRepo, dstRepoName)
				if err != nil {
					return err
				}
			} else {
				logger.Warnf("(%d/%d) source repo %s not in Org %s/%s, skip.", i+1, total, srcRepoName, m.SrcGit, m.SrcOrg)
			}
		}
	} else {
		// mirror all repos
		total := len(m.srcRepos)
		for i, srcRepo := range m.srcRepos {
			if m.isMirrorRepo(*srcRepo.Name) == true {
				dstRepoName := m.getDstRepoName(*srcRepo.Name)
				logger.Debugf("(%d/%d) begin mirror %s/%s/%s to %s/%s/%s",
					i+1, total, m.SrcGit, m.SrcOrg, *srcRepo.Name, m.DstGit, m.DstOrg, dstRepoName)
				err := m.mirror(srcRepo, dstRepoName)
				if err != nil {
					return err
				}
			} else {
				logger.Warnf("(%d/%d) source repo %s of Org %s/%s maybe in black-list, skip.", i+1, total, *srcRepo.Name, m.SrcGit, m.SrcOrg)
			}
		}
	}

	return nil
}
