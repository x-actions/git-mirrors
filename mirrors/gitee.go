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
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"golang.org/x/oauth2"
)

const (
	maxGiteePerPage = 100
)

type GiteeAPI struct {
	Client      *gitee.APIClient
	Context     context.Context
	accessToken string
	IsAuthed    bool
}

// NewGiteeAPI return new Gitee API
func NewGiteeAPI(accessToken string) (*GiteeAPI, error) {
	ctx := context.Background()
	// configuration
	conf := gitee.NewConfiguration()
	isAuthed := false
	if accessToken != "" {
		// oauth
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		conf.HTTPClient = oauth2.NewClient(ctx, ts)
		isAuthed = true
	}

	// git client
	client := gitee.NewAPIClient(conf)

	return &GiteeAPI{Client: client, Context: ctx, accessToken: accessToken, IsAuthed: isAuthed}, nil
}

// IsAPIAuthed return is the API auth, true or false
func (g *GiteeAPI) IsAPIAuthed() bool {
	return g.IsAuthed
}

// Organizations list all Organizations
func (g *GiteeAPI) Organizations(user string) ([]*Organization, error) {
	page := 1
	opt := &gitee.GetV5UsersUsernameOrgsOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(int32(page)),
		PerPage:     optional.NewInt32(maxGiteePerPage),
	}
	groups, _, err := g.Client.OrganizationsApi.GetV5UsersUsernameOrgs(g.Context, user, opt)
	baseOrgs := make([]*Organization, len(groups))
	for i, group := range groups {
		baseOrgs[i] = formatGiteeGroup(group)
	}
	return baseOrgs, err
}

func (g *GiteeAPI) GetOrganization(orgName string) (*Organization, error) {
	opt := &gitee.GetV5OrgsOrgOpts{
		AccessToken: optional.NewString(g.accessToken),
	}
	group, resp, err := g.Client.OrganizationsApi.GetV5OrgsOrg(g.Context, orgName, opt)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound("Organization", orgName)
		}
		return nil, err
	}
	return formatGiteeGroup(group), nil
}

// Repositories list all repositories for the authenticated user
func (g *GiteeAPI) Repositories(user string) ([]*Repository, error) {
	page := 1
	opt := &gitee.GetV5UserReposOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(int32(page)),
		PerPage:     optional.NewInt32(int32(maxGiteePerPage)),
	}
	var baseRepos []*Repository
	for {
		projects, _, err := g.Client.RepositoriesApi.GetV5UserRepos(g.Context, opt)
		if err != nil {
			return nil, err
		}
		for _, repo := range projects {
			baseRepos = append(baseRepos, formatGiteeRepo(repo))
		}

		if len(projects) < maxGiteePerPage {
			break
		}

		page += 1
		opt.Page = optional.NewInt32(int32(page))
	}

	return baseRepos, nil
}

// GetRepository fetches a repository
func (g *GiteeAPI) GetRepository(orgName, repoName string) (*Repository, error) {
	opt := &gitee.GetV5ReposOwnerRepoOpts{
		AccessToken: optional.NewString(g.accessToken),
	}
	project, _, err := g.Client.RepositoriesApi.GetV5ReposOwnerRepo(g.Context, orgName, repoName, opt)
	if err != nil {
		return nil, err
	}
	return formatGiteeRepo(project), nil
}

// CreateUserRepo create a new user repository
func (g *GiteeAPI) CreateUserRepo(baseRepo *Repository) (*Repository, error) {
	opt := gitee.RepositoryPostParam{
		AccessToken: g.accessToken,
		Name:        *baseRepo.Name,
	}
	if baseRepo.Description != nil {
		opt.Description = *baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		opt.Homepage = *baseRepo.Homepage
	}
	if baseRepo.Private != nil {
		opt.Private = *baseRepo.Private
	}
	project, _, err := g.Client.RepositoriesApi.PostV5UserRepos(g.Context, *baseRepo.Name, opt)
	if err != nil {
		return nil, err
	}
	return formatGiteeRepo(project), nil
}

// CreateOrgRepo create a new org repository
func (g *GiteeAPI) CreateOrgRepo(baseRepo *Repository, orgName string) (*Repository, error) {
	opt := gitee.RepositoryPostParam{
		AccessToken: g.accessToken,
		Name:        *baseRepo.Name,
	}
	if baseRepo.Description != nil {
		opt.Description = *baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		opt.Homepage = *baseRepo.Homepage
	}
	if baseRepo.Private != nil {
		opt.Private = *baseRepo.Private
	}
	project, _, err := g.Client.RepositoriesApi.PostV5OrgsOrgRepos(g.Context, orgName, opt)
	if err != nil {
		return nil, err
	}
	return formatGiteeRepo(project), nil
}

// CreateRepository create a new repository, if repo is already exist, just return it
func (g *GiteeAPI) CreateRepository(baseRepo *Repository, orgName string) (*Repository, error) {
	_, err := g.GetOrganization(orgName)
	if err != nil {
		// user
		return g.CreateUserRepo(baseRepo)
	} else {
		// org
		return g.CreateOrgRepo(baseRepo, orgName)
	}
}

// UpdateRepository updates a repository
func (g *GiteeAPI) UpdateRepository(orgName, repoName string, baseRepo *Repository) (*Repository, error) {
	opt := gitee.RepoPatchParam{
		AccessToken: g.accessToken,
		Name:        repoName,
	}
	if baseRepo.Description != nil {
		opt.Description = *baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		opt.Homepage = *baseRepo.Homepage
	}
	if baseRepo.Private != nil {
		opt.Private = strconv.FormatBool(*baseRepo.Private)
	}
	project, _, err := g.Client.RepositoriesApi.PatchV5ReposOwnerRepo(g.Context, orgName, repoName, opt)
	if err != nil {
		return nil, err
	}
	return formatGiteeRepo(project), nil
}

// RepositoriesByOrg list repositories for special org
func (g *GiteeAPI) RepositoriesByOrg(orgName string) ([]*Repository, error) {
	page := 1
	opt := &gitee.GetV5OrgsOrgReposOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(int32(page)),
		PerPage:     optional.NewInt32(int32(maxGiteePerPage)),
	}

	var baseRepos []*Repository
	for {
		projects, resp, err := g.Client.RepositoriesApi.GetV5OrgsOrgRepos(g.Context, orgName, opt)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrNotFound("Organization", orgName)
			}
			return nil, err
		}
		for _, repo := range projects {
			baseRepos = append(baseRepos, formatGiteeRepo(repo))
		}

		if len(projects) < maxGiteePerPage {
			break
		}

		page += 1
		opt.Page = optional.NewInt32(int32(page))
	}

	return baseRepos, nil
}

func formatGiteeGroup(repo gitee.Group) *Organization {
	baseRepo := &Organization{
		Name:        &repo.Login,
		Description: &repo.Description,
	}

	return baseRepo
}

func formatGiteeRepo(project gitee.Project) *Repository {
	htmlURL := strings.TrimSuffix(project.HtmlUrl, ".git")
	baseRepo := &Repository{
		Owner: &User{
			Name: &project.Owner.Login,
			Type: &project.Owner.Type_,
		},
		Name:        &project.Name,
		FullName:    &project.FullName,
		Description: &project.Description,
		HTMLURL:     &htmlURL,
		CloneURL:    &project.HtmlUrl,
		GitURL:      &project.SshUrl,
		SSHURL:      &project.SshUrl,
		Homepage:    &project.Homepage,
		Fork:        &project.Fork,
		Topics:      []string{},
		Private:     &project.Private,
		//Archived: optional.NewBool(true),
	}

	if project.Namespace != nil {
		baseRepo.Organization = &Organization{
			Name: &project.Namespace.Name,
			//Description: &repo.Namespace.Description,
			Type: &project.Namespace.Type_,
		}
	}

	return baseRepo
}
