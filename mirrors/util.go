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

func RemoveDuplicates(strs []string) []string {
	keys := make(map[string]struct{}, len(strs))
	d := 0
	for i, s := range strs {
		if _, ok := keys[s]; ok {
			strs[d], strs[i] = strs[i], strs[d]
			d++
		} else {
			keys[s] = struct{}{}
		}
	}
	return strs[d:]
}

// StringListToMap convert strings list to map
func StringListToMap(strs []string) map[string]string {
	result := make(map[string]string, len(strs))
	for _, str := range strs {
		result[str] = str
	}

	return result
}

// ReposToMap convert Repository List to map
func ReposToMap(repos []*Repository) map[string]*Repository {
	result := make(map[string]*Repository, len(repos))
	for _, repo := range repos {
		result[*repo.Name] = repo
	}

	return result
}

// StringsEqual compare two string point; if value equal return true; else false
func StringsEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && *a == "" && b == nil {
		return true
	}
	if a == nil && b != nil && *b == "" {
		return true
	}
	if a != nil && b != nil && *a == *b {
		return true
	}

	return false
}
