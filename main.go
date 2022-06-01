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
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xiexianbin/golib/logger"

	"github.com/x-actions/git-mirrors/constants"
	"github.com/x-actions/git-mirrors/mirrors"
)

var (
	src            string
	srcToken       string
	dst            string
	dstKey         string
	dstToken       string
	accountType    string
	srcAccountType string
	dstAccountType string
	cloneStyle     string
	cachePath      string
	blackList      []string
	blackListStr   string
	whiteList      []string
	whiteListStr   string
	forceUpdate    bool
	debug          bool
	timeoutStr     string
	timeout        time.Duration
	mappingsStr    string
	mappings       map[string]string

	help    bool
	version bool
	verbose bool
)

var (
	srcGit string
	srcOrg string
	dstGit string
	dstOrg string
)

func init() {
	flag.StringVar(&src, "src", "", "Source name. Such as `github/xiexianbin`")
	flag.StringVar(&srcToken, "src-token", "", "The app token which is used to list repo in source hub")
	flag.StringVar(&dst, "dst", "", "Destination name. Such as `gitee/xiexianbin`")
	flag.StringVar(&dstKey, "dst-key", "", "The private SSH key which is used to to push code in destination hub")
	flag.StringVar(&dstToken, "dst-token", "", "The app token which is used to create repo in destination hub")
	flag.StringVar(&accountType, "account-type", "user", "The account type. Such as org, user")
	flag.StringVar(&srcAccountType, "src-account-type", "", "The src account type. Such as org, user")
	flag.StringVar(&dstAccountType, "dst-account-type", "", "The dst account type. Such as org, user")
	flag.StringVar(&cloneStyle, "clone-style", "ssh", "The git clone style, https or ssh")
	flag.StringVar(&cachePath, "cache-path", "/github/workspace/git-mirrors-cache", "The path to cache the source repos code")
	flag.StringVar(&blackListStr, "black-list", "", "Height priority, the back list of mirror repo. like 'repo1,repo2,repo3'")
	flag.StringVar(&whiteListStr, "white-list", "", "Low priority, the white list of mirror repo. like 'repo1,repo2,repo3'")
	flag.BoolVar(&forceUpdate, "force-update", false, "Force to update the destination repo, use '-f' flag do 'git push'")
	flag.BoolVar(&debug, "debug", false, "Enable the debug flag to show detail log")
	flag.StringVar(&timeoutStr, "timeout", "30m", "Set the timeout for every git command, eg. '600s'=>600s, '30m'=>30 minute, '2h'=>2 hours")
	flag.StringVar(&mappingsStr, "mappings", "", "The source repos mappings, such as 'A=>B, C=>CC', source repo name would be mapped follow the rule: A to B, C to CC. Mapping is not transitive")

	flag.BoolVar(&help, "h", false, "print this help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&verbose, "V", false, "be verbose, debug model")

	flag.Parse()

	flag.Usage = func() {
		fmt.Println("Run: git-mirrors -h\n\n" +
			"Usage:")
		flag.PrintDefaults()
	}
}

func parseParams() error {
	// parse black and white list
	if blackListStr == "" {
		blackList = []string{}
	} else {
		blackList = strings.Split(blackListStr, ",")
	}
	if whiteListStr == "" {
		whiteList = []string{}
	} else {
		whiteList = strings.Split(whiteListStr, ",")
	}

	// check source and destination git service
	checkGitSource := func(str string) ([]string, bool) {
		if gitInfo := strings.Split(str, "/"); len(gitInfo) == 2 {
			for _, g := range constants.SupportGit {
				if g == gitInfo[0] {
					return gitInfo, true
				}
			}
		}

		return nil, false
	}

	if gitInfo, ok := checkGitSource(src); ok == false {
		return fmt.Errorf("un-support git source %s", src)
	} else {
		srcGit, srcOrg = gitInfo[0], gitInfo[1]
	}
	if gitInfo, ok := checkGitSource(dst); ok == false {
		return fmt.Errorf("un-support git destination %s", src)
	} else {
		dstGit, dstOrg = gitInfo[0], gitInfo[1]
	}

	// account type
	if srcAccountType == "" {
		srcAccountType = accountType
	}
	if dstAccountType == "" {
		dstAccountType = accountType
	}

	// parse mappings
	maps := strings.Split(mappingsStr, ",")
	for _, m := range maps {
		if m == "" {
			continue
		}
		if r := strings.Split(m, "=>"); len(r) == 2 {
			mappings[r[0]] = r[1]
		} else {
			return fmt.Errorf("parse mappings: %s format invalied", m)
		}
	}

	// parse timeout
	var err error
	timeout, err = time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("parse timeout %s err: %s", timeoutStr, err.Error())
	}

	// token check
	if srcToken == "" {
		logger.Warn("un-configure srcToken, Only mirror Public Repos.")
	}

	return nil
}

func main() {
	if help == true {
		flag.Usage()
		return
	}

	if version == true {
		showVersion()
		return
	}

	if verbose == true || debug == true {
		logger.SetLogLevel(logger.DEBUG)
	}

	if err := parseParams(); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	mirror := mirrors.New(srcGit, srcOrg, srcToken, dstGit, dstOrg, dstKey, dstToken, srcAccountType, dstAccountType,
		cloneStyle, cachePath, blackList, whiteList, forceUpdate, timeout, mappings)
	err := mirror.Do()
	if err != nil {
		logger.Fatalf("%s", err.Error())
		os.Exit(1)
	}
}
