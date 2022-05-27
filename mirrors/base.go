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

import "fmt"

type IGit interface {
	Organizations(user string) ([]*Organization, error)
	GetOrganization(orgName string) (*Organization, error)
	Repositories(user string) ([]*Repository, error)
	GetRepository(orgName, repoName string) (*Repository, error)
	CreateRepository(baseRepo *Repository, orgName string) (*Repository, error)
	UpdateRepository(orgName, repoName string, baseRepo *Repository) (*Repository, error)
	RepositoriesByOrg(orgName string) ([]*Repository, error)
}

type User struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

type Organization struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"` // github: Organization or User; gitee: personal or group
}

// Repository represents a Common repository.
type Repository struct {
	Owner        *User         `json:"owner,omitempty"`
	Name         *string       `json:"name,omitempty"`
	FullName     *string       `json:"full_name,omitempty"`
	Description  *string       `json:"description,omitempty"`
	HTMLURL      *string       `json:"html_url,omitempty"`
	CloneURL     *string       `json:"clone_url,omitempty"`
	GitURL       *string       `json:"git_url,omitempty"`
	SSHURL       *string       `json:"ssh_url,omitempty"`
	Homepage     *string       `json:"homepage,omitempty"`
	Fork         *bool         `json:"fork,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
	Topics       []string      `json:"topics,omitempty"`

	// Additional mutable fields when creating and editing a repository
	Private  *bool `json:"private,omitempty"`
	Archived *bool `json:"archived,omitempty"`
}

func ErrNotFound(resource, name string) error {
	return fmt.Errorf("resource %s %s not found", resource, name)
}
