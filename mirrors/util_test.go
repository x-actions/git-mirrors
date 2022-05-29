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
	"testing"

	"github.com/google/go-github/github"
)

func TestRemoveDuplicates(t *testing.T) {
	t.Logf("%#v", RemoveDuplicates([]string{"abc", "abc", "def"}))
}

func TestStringListToMap(t *testing.T) {
	t.Logf("%#v", StringListToMap([]string{"abc", "def"}))
}

func TestReposToMap(t *testing.T) {
	name := "abc"
	repos := []*Repository{
		{Name: github.String(name)},
	}
	result := ReposToMap(repos)
	t.Logf("%#v", *result[name].Name)
}
