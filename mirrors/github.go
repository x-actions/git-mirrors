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

// https://github.com/xiexianbin/gsync

package mirrors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/xiexianbin/golib/logger"
	"golang.org/x/oauth2"
)

type GithubAPI struct {
	Client      *github.Client
	Context     context.Context
	accessToken string
}

// NewGithubAPI return new Github API
func NewGithubAPI(accessToken string) (*GithubAPI, error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	if accessToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	return &GithubAPI{Client: client, Context: ctx, accessToken: accessToken}, nil
}

// Organizations list Organizations
func (g *GithubAPI) Organizations(user string) ([]*Organization, error) {
	opt := &github.ListOptions{
		Page:    1,
		PerPage: 1000,
	}
	orgs, _, err := g.Client.Organizations.List(g.Context, user, opt)
	baseOrgs := make([]*Organization, len(orgs))
	for i, org := range orgs {
		baseOrgs[i] = formatGithubOrg(org)
	}

	return baseOrgs, err
}

// GetOrganization get an organization by name
func (g *GithubAPI) GetOrganization(orgName string) (*Organization, error) {
	if orgName == "" {
		return nil, fmt.Errorf("new repo name must not be empty")
	}

	githubOrg, resp, err := g.Client.Organizations.Get(g.Context, orgName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound("Organization", orgName)
		}
		return nil, err
	}
	logger.Debugf("get github %s Organization: %v", orgName, githubOrg)

	return formatGithubOrg(githubOrg), nil
}

// Repositories list all repositories for the authenticated user, if user is empty show all repos
// https://docs.github.com/en/rest/repos/repos#list-repositories-for-the-authenticated-user
func (g *GithubAPI) Repositories(user string) ([]*Repository, error) {
	opt := &github.RepositoryListOptions{
		//Type:        "all",
		Affiliation: "owner",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 1000,
		},
	}
	repos, _, err := g.Client.Repositories.List(g.Context, user, opt)
	baseRepos := make([]*Repository, len(repos))
	for i, repo := range repos {
		baseRepos[i] = formatGithubRepo(repo)
	}
	return baseRepos, err
}

// GetRepository fetches a repository
func (g *GithubAPI) GetRepository(orgName, repoName string) (*Repository, error) {
	repo, _, err := g.Client.Repositories.Get(g.Context, orgName, repoName)
	if err != nil {
		return nil, err
	}

	return formatGithubRepo(repo), nil
}

// CreateRepository create a new repository, if repo is already exist, just return it
func (g *GithubAPI) CreateRepository(baseRepo *Repository, orgName string) (*Repository, error) {
	if *baseRepo.Name == "" {
		return nil, fmt.Errorf("new repo name must not be empty")
	}
	repo := &github.Repository{
		Name: github.String(*baseRepo.Name),
	}
	if baseRepo.Description != nil {
		repo.Description = baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		repo.Homepage = baseRepo.Homepage
	}
	if baseRepo.Topics != nil {
		repo.Topics = baseRepo.Topics
	}
	if baseRepo.Private != nil {
		repo.Private = baseRepo.Private
	}
	githubRepo, resp, err := g.Client.Repositories.Create(g.Context, orgName, repo)
	if err != nil {
		if resp.StatusCode == http.StatusUnprocessableEntity {
			// 422 Repository creation failed.
			// [{Resource:Repository Field:name Code:custom Message:name already exists on this account}]
			if baseRepo, err := g.GetRepository(orgName, *repo.Name); err == nil {
				return baseRepo, nil
			}
		}
		return nil, err
	}
	return formatGithubRepo(githubRepo), nil
}

// UpdateRepository updates a repository
func (g *GithubAPI) UpdateRepository(orgName, repoName string, baseRepo *Repository) (*Repository, error) {
	_, err := g.GetRepository(orgName, repoName)
	if err != nil {
		return nil, fmt.Errorf("get %s/%s err: %s", orgName, repoName, err.Error())
	}

	_githubRepo := &github.Repository{}
	if baseRepo.Description != nil {
		_githubRepo.Description = baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		_githubRepo.Homepage = baseRepo.Homepage
	}
	if baseRepo.Topics != nil {
		_githubRepo.Topics = baseRepo.Topics
	}
	if baseRepo.Private != nil {
		_githubRepo.Private = baseRepo.Private
	}
	githubRepo, _, err := g.Client.Repositories.Edit(g.Context, orgName, repoName, _githubRepo)
	if err != nil {
		return nil, err
	}
	return formatGithubRepo(githubRepo), nil
}

// RepositoriesByOrg list repositories for special org
func (g *GithubAPI) RepositoriesByOrg(orgName string) ([]*Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		Type: "all", // Possible values are: all, public, private, forks, sources, member. Default is "all".
	}
	repos, _, err := g.Client.Repositories.ListByOrg(g.Context, orgName, opt)
	baseRepos := make([]*Repository, len(repos))
	for i, repo := range repos {
		baseRepos[i] = formatGithubRepo(repo)
	}
	return baseRepos, err
}

func formatGithubOrg(repo *github.Organization) *Organization {
	baseRepo := &Organization{
		Name:        repo.Login,
		Description: repo.Description,
		Type:        repo.Type,
	}

	return baseRepo
}

func formatGithubRepo(repo *github.Repository) *Repository {
	baseRepo := &Repository{
		Owner: &User{
			Name: repo.Owner.Login,
			Type: repo.Owner.Type,
		},
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		HTMLURL:     repo.HTMLURL,
		CloneURL:    repo.CloneURL,
		GitURL:      repo.GitURL,
		SSHURL:      repo.SSHURL,
		Homepage:    repo.Homepage,
		Fork:        repo.Fork,
		Topics:      repo.Topics,
		Private:     repo.Private,
		Archived:    repo.Archived,
	}

	if repo.Organization != nil {
		baseRepo.Organization = &Organization{
			Name:        repo.Organization.Name,
			Description: repo.Organization.Description,
		}
	}

	return baseRepo
}
