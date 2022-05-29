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

// https://gitee.com/xiexianbin/gsync

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
	GiteeTokenKey = "GITEE_TOKEN"
	GiteeUserName = "e-stack"
	GiteeOrgName  = "pncx"
	GiteeTestRepo = "test"
)

func TestGitee_Organizations(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	orgs, err := c.Organizations(GiteeUserName)
	if err != nil {
		t.Skipf("get gitee Organizations err: %s", err.Error())
		return
	}
	for _, org := range orgs {
		j, _ := json.Marshal(org)
		t.Log(string(j))
	}
}

func TestGitee_Repositories(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	repos, err := c.Repositories(GiteeUserName)
	if err != nil {
		t.Skipf("get gitee Repositories err: %s", err.Error())
		return
	}
	for _, repo := range repos {
		j, _ := json.Marshal(repo)
		t.Log(string(j))
	}
}

func TestGitee_GetRepository(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	repo, err := c.GetRepository(GiteeUserName, GiteeTestRepo)
	if err != nil {
		t.Skipf("get gitee Repositorie err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(repo)
	t.Log(string(j))

	repo, err = c.GetRepository(GiteeOrgName, GiteeTestRepo)
	if err != nil {
		t.Skipf("get gitee Repositorie err: %s", err.Error())
		return
	}
	j, _ = json.Marshal(repo)
	t.Log(string(j))
}

func TestGitee_CreateRepository(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	baseRepo := &Repository{
		Name:        github.String("test-create-repo"),
		Description: github.String("i am description."),
		Homepage:    github.String("https://www.xiexianbin.cn"),
		Private:     github.Bool(true),
	}

	//repo, err := c.CreateRepository(baseRepo, GITEE_ORG_NAME)
	repo, err := c.CreateRepository(baseRepo, GiteeUserName)
	if err != nil {
		t.Skipf("create gitee Repositories err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(repo)
	t.Log(string(j))
}

func TestGitee_UpdateRepository(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	baseRepo := &Repository{
		Description: github.String(fmt.Sprintf("i am description, date: %s.", time.Now().Format("2006-01-02 15:04:06"))),
	}

	//repo, err := c.UpdateRepository(GITEE_ORG_NAME, "test-create-repo", baseRepo)
	repo, err := c.UpdateRepository(GiteeUserName, "test-create-repo", baseRepo)
	if err != nil {
		t.Skipf("update gitee Repositories err: %s", err.Error())
		return
	}
	j, _ := json.Marshal(repo)
	t.Log(string(j))
}

func TestGitee_RepositoriesByOrg(t *testing.T) {
	accessToken := os.Getenv(GiteeTokenKey)
	c, err := NewGiteeAPI(accessToken)
	if err != nil {
		t.Skipf("init gitee api client err: %s", err.Error())
		return
	}

	repos, err := c.RepositoriesByOrg(GiteeOrgName)
	if err != nil {
		t.Skipf("get gitee RepositoriesByOrg %s err: %s", GiteeOrgName, err.Error())
		return
	}
	for _, repo := range repos {
		j, _ := json.Marshal(repo)
		t.Log(string(j))
	}
}
