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
	"github.com/x-actions/git-mirrors/constants"
	"github.com/xiexianbin/golib/logger"
	"time"
)

type Mirror struct {
	SrcGit      string
	SrcOrg      string
	srcToken    string
	DstGit      string
	DstOrg      string
	dstKey      string
	dstToken    string
	AccountType string
	CloneStyle  string
	CachePath   string
	BlackList   []string
	WhiteList   []string
	ForceUpdate bool
	Timeout     time.Duration
	Mappings    map[string]string

	srcRepos []*Repository
	dstRepos []*Repository

	gitClient *GitClient
	githubAPI *GithubAPI
	giteeAPI  *GiteeAPI
}

func New(srcGit, srcOrg, srcToken, dstGit, dstOrg, dstKey, dstToken, accountType, cloneStyle, cachePath string,
	blackList, whiteList []string, forceUpdate bool, timeout time.Duration, mappings map[string]string) *Mirror {
	return &Mirror{
		SrcGit:      srcGit,
		SrcOrg:      srcOrg,
		srcToken:    srcToken,
		DstGit:      dstGit,
		DstOrg:      dstOrg,
		dstKey:      dstKey,
		dstToken:    dstToken,
		AccountType: accountType,
		CloneStyle:  cloneStyle,
		CachePath:   cachePath,
		BlackList:   blackList,
		WhiteList:   whiteList,
		ForceUpdate: forceUpdate,
		Timeout:     timeout,
		Mappings:    mappings,
	}
}

// prepare init src/dst APIs and Repos
func (m *Mirror) prepare() error {
	// init src
	switch m.SrcGit {
	case constants.GITHUB:
		// init Github Client
		client, err := NewGithubAPI(m.srcToken)
		if err != nil {
			return err
		}
		m.githubAPI = client

		// src repos
		switch m.AccountType {
		case constants.AccountTypeUser:
			repos, err := m.githubAPI.Repositories(m.SrcOrg)
			if err != nil {
				return err
			}
			m.srcRepos = repos
		case constants.AccountTypeOrg:
			repos, err := m.githubAPI.RepositoriesByOrg(m.SrcOrg)
			if err != nil {
				return err
			}
			m.srcRepos = repos
		default:
			return fmt.Errorf("un-support account-type %s", m.AccountType)
		}
	case constants.GITEE:
		// init Gitee Client
		client, err := NewGiteeAPI(m.srcToken)
		if err != nil {
			return err
		}
		m.giteeAPI = client

		// src repos
		switch m.AccountType {
		case constants.AccountTypeUser:
			repos, err := m.giteeAPI.Repositories(m.SrcOrg)
			if err != nil {
				return err
			}
			m.srcRepos = repos
		case constants.AccountTypeOrg:
			repos, err := m.giteeAPI.RepositoriesByOrg(m.SrcOrg)
			if err != nil {
				return err
			}
			m.srcRepos = repos
		default:
			return fmt.Errorf("un-support account-type %s", m.AccountType)
		}
	default:
		return fmt.Errorf("un-support git source %s", m.SrcGit)
	}

	// init dst
	switch m.DstGit {
	case constants.GITHUB:
		// init Github Client
		api, err := NewGithubAPI(m.srcToken)
		if err != nil {
			return err
		}
		m.githubAPI = api

		// dst repos
		switch m.AccountType {
		case constants.AccountTypeUser:
			repos, err := m.githubAPI.Repositories(m.DstOrg)
			if err != nil {
				return err
			}
			m.dstRepos = repos
		case constants.AccountTypeOrg:
			repos, err := m.githubAPI.RepositoriesByOrg(m.DstOrg)
			if err != nil {
				return err
			}
			m.dstRepos = repos
		default:
			return fmt.Errorf("un-support account-type %s", m.AccountType)
		}
	case constants.GITEE:
		// init Gitee Client
		api, err := NewGiteeAPI(m.srcToken)
		if err != nil {
			return err
		}
		m.giteeAPI = api

		// dst repos
		switch m.AccountType {
		case constants.AccountTypeUser:
			repos, err := m.giteeAPI.Repositories(m.DstOrg)
			if err != nil {
				return err
			}
			m.dstRepos = repos
		case constants.AccountTypeOrg:
			repos, err := m.giteeAPI.RepositoriesByOrg(m.DstOrg)
			if err != nil {
				return err
			}
			m.dstRepos = repos
		default:
			return fmt.Errorf("un-support account-type %s", m.AccountType)
		}
	default:
		return fmt.Errorf("un-support git destination %s", m.DstGit)
	}

	return nil
}

// syncRepoInfo create or sync Repo Info
func (m *Mirror) syncRepoInfo() error {

	return nil
}

// syncGit clone/pull from src repo and push to dst repo
func (m *Mirror) syncGit() error {

	return nil
}

// Do mirror logic
func (m *Mirror) Do() error {
	// get src/dst Repos
	err := m.prepare()
	if err != nil {
		return err
	}

	total := len(m.srcRepos)
	for i, srcRepo := range m.srcRepos {
		logger.Debugf("begin sync %d/%s", i+1, total)
		// sync repo infos

		// sync git commits
	}

	return nil
}
