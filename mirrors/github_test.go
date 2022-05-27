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
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/github"
)

const (
	GITHUB_TOKEN_KEY = "GITHUB_TOKEN"
	GITHUB_USER_NAME = "estack"
	GITHUB_ORG_NAME  = "pncx"
)

func TestGithub_Organizations(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	orgs, err := c.Organizations("")
	if err != nil {
		t.Skipf("get github Organizations err: %s", err.Error())
		return
	}
	for _, org := range orgs {
		j, _ := json.Marshal(org)
		t.Log(string(j))
	}
}

func TestGithub_GetOrganization(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	org, err := c.GetOrganization(GITHUB_ORG_NAME)
	if err != nil {
		t.Skipf("get github Organization err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(org)
	t.Log(string(j))

	org, err = c.GetOrganization(GITHUB_USER_NAME)
	if err != nil {
		t.Skipf("get github Organization err: %s", err.Error())
		return
	}
	j, _ = json.Marshal(org)
	t.Log(string(j))
}

func TestGithub_Repositories(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	repos, err := c.Repositories("")
	if err != nil {
		t.Skipf("get github Repositories err: %s", err.Error())
		return
	}
	for _, repo := range repos {
		j, _ := json.Marshal(repo)
		t.Log(string(j))
	}
}

func TestGithub_CreateRepository(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	baseRepo := &Repository{
		Name:        github.String("test-create-repo"),
		Description: github.String("i am description."),
		Homepage:    github.String("https://www.xiexianbin.cn"),
		Private:     github.Bool(true),
	}

	repo, err := c.CreateRepository(baseRepo, GITHUB_ORG_NAME)
	if err != nil {
		t.Skipf("create github Repositories err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(repo)
	t.Log(string(j))
}

func TestGithub_UpdateRepository(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	baseRepo := &Repository{
		Description: github.String(fmt.Sprintf("i am description, date: %s.", time.Now().Format("2006-01-02 15:04:06"))),
	}

	repo, err := c.UpdateRepository(GITHUB_ORG_NAME, "test-create-repo", baseRepo)
	if err != nil {
		t.Skipf("update github Repositories err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(repo)
	t.Log(string(j))
}

func TestGithub_RepositoriesByOrg(t *testing.T) {
	accessToken := os.Getenv(GITHUB_TOKEN_KEY)
	c, err := NewGithubAPI(accessToken)
	if err != nil {
		t.Skipf("init github api client err: %s", err.Error())
		return
	}

	org := GITHUB_ORG_NAME
	repos, err := c.RepositoriesByOrg(org)
	if err != nil {
		t.Skipf("get github RepositoriesByOrg %s err: %s", org, err.Error())
		return
	}
	for _, repo := range repos {
		j, _ := json.Marshal(repo)
		t.Log(string(j))
	}
}
