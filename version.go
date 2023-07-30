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

package main

import (
	"fmt"
	"runtime"

	"github.com/xiexianbin/golib/logger"
)

var (
	version      = "v0.0.0"               // value from VERSION file
	buildDate    = "1970-01-01T00:00:00Z" // output from `TZ=UTC-8 date +'%Y-%m-%dT%H:%M:%SZ+08:00'`
	gitCommit    = ""                     // output from `git rev-parse HEAD`
	gitTag       = ""                     // output from `git describe --exact-match --tags HEAD` (if clean tree state)
	gitTreeState = ""                     // determined from `git status --porcelain`. either 'clean' or 'dirty'
)

type Version struct {
	Version      string `json:"version"`
	BuildDate    string `json:"buildDate"`
	GitCommit    string `json:"gitCommit"`
	GitTag       string `json:"gitTag"`
	GitTreeState string `json:"gitTreeState"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// GetVersion returns the version information
func GetVersion() Version {
	var versionStr string
	if gitCommit != "" && gitTag != "" && gitTreeState == "clean" {
		// if we have a clean tree state and the current commit is tagged,
		// this is an official release.
		versionStr = gitTag
	} else {
		// otherwise formulate a version string based on as much metadata
		// information we have available.
		versionStr = version
		if len(gitCommit) >= 7 {
			versionStr += "+" + gitCommit[0:7]
			if gitTreeState != "clean" {
				versionStr += ".dirty"
			}
		} else {
			versionStr += "+unknown"
		}
	}
	return Version{
		Version:      versionStr,
		BuildDate:    buildDate,
		GitCommit:    gitCommit,
		GitTag:       gitTag,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func PrintVersion(cliName string, version Version, short bool) {
	logger.Printf("%s: %s", cliName, version.Version)
	if short {
		return
	}
	logger.Printf("  BuildDate: %s", version.BuildDate)
	logger.Printf("  GitCommit: %s", version.GitCommit)
	logger.Printf("  GitTreeState: %s", version.GitTreeState)
	if version.GitTag != "" {
		logger.Printf("  GitTag: %s", version.GitTag)
	}
	logger.Printf("  GoVersion: %s", version.GoVersion)
	logger.Printf("  Compiler: %s", version.Compiler)
	logger.Printf("  Platform: %s", version.Platform)
}
