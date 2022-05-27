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
	"strings"

	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/antihax/optional"
	"golang.org/x/oauth2"
)

type Gitee struct {
	Client      *gitee.APIClient
	Context     context.Context
	accessToken string
}

// NewGiteeAPI return new Gitee API
func NewGiteeAPI(accessToken string) (*Gitee, error) {
	ctx := context.Background()
	// configuration
	conf := gitee.NewConfiguration()
	if accessToken != "" {
		// oauth
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		conf.HTTPClient = oauth2.NewClient(ctx, ts)
	}

	// git client
	client := gitee.NewAPIClient(conf)

	return &Gitee{
		Client:      client,
		Context:     ctx,
		accessToken: accessToken,
	}, nil
}

// Organizations list all Organizations
func (g *Gitee) Organizations(user string) ([]*Organization, error) {
	opt := &gitee.GetV5UsersUsernameOrgsOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(1),
		PerPage:     optional.NewInt32(1000),
	}
	groups, _, err := g.Client.OrganizationsApi.GetV5UsersUsernameOrgs(g.Context, user, opt)
	baseOrgs := make([]*Organization, len(groups))
	for i, group := range groups {
		baseOrgs[i] = formatGiteeGroup(group)
	}
	return baseOrgs, err
}

func (g *Gitee) GetOrganization(orgName string) (*Organization, error) {
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
func (g *Gitee) Repositories(user string) ([]*Repository, error) {
	opt := &gitee.GetV5UserReposOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(int32(1)),
		PerPage:     optional.NewInt32(int32(100)),
	}
	projects, _, err := g.Client.RepositoriesApi.GetV5UserRepos(g.Context, opt)
	baseRepos := make([]*Repository, len(projects))
	for i, repo := range projects {
		baseRepos[i] = formatGiteeRepo(repo)
	}
	return baseRepos, err
}

// GetRepository fetches a repository
func (g *Gitee) GetRepository(orgName, repoName string) (*Repository, error) {
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
func (g *Gitee) CreateUserRepo(baseRepo *Repository) (*Repository, error) {
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
func (g *Gitee) CreateOrgRepo(baseRepo *Repository, orgName string) (*Repository, error) {
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
func (g *Gitee) CreateRepository(baseRepo *Repository, orgName string) (*Repository, error) {
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
func (g *Gitee) UpdateRepository(orgName, repoName string, baseRepo *Repository) (*Repository, error) {
	opt := gitee.RepoPatchParam{
		AccessToken: g.accessToken,
		Owner:       orgName,
		Repo:        repoName,
		Name:        repoName,
	}
	if baseRepo.Description != nil {
		opt.Description = *baseRepo.Description
	}
	if baseRepo.Homepage != nil {
		opt.Homepage = *baseRepo.Homepage
	}
	if baseRepo.Private != nil {
		opt.Private = fmt.Sprintf("%t", *baseRepo.Private)
	}
	project, _, err := g.Client.RepositoriesApi.PatchV5ReposOwnerRepo(g.Context, orgName, repoName, opt)
	if err != nil {
		return nil, err
	}
	return formatGiteeRepo(project), nil
}

// RepositoriesByOrg list repositories for special org
func (g *Gitee) RepositoriesByOrg(orgName string) ([]*Repository, error) {
	opt := &gitee.GetV5OrgsOrgReposOpts{
		AccessToken: optional.NewString(g.accessToken),
		Page:        optional.NewInt32(int32(1)),
		PerPage:     optional.NewInt32(int32(100)),
	}
	projects, resp, err := g.Client.RepositoriesApi.GetV5OrgsOrgRepos(g.Context, orgName, opt)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound("Organization", orgName)
		}
		return nil, err
	}
	baseRepos := make([]*Repository, len(projects))
	for i, project := range projects {
		baseRepos[i] = formatGiteeRepo(project)
	}
	return baseRepos, err
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
