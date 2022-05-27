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
	"github.com/xiexianbin/golib/logger"
)

var (
	src          string
	srcToken     string
	dst          string
	dstKey       string
	dstToken     string
	accountType  string
	cloneStyle   string
	cachePath    string
	blackList    []string
	blackListStr string
	whiteList    []string
	whiteListStr string
	forceUpdate  bool
	debug        bool
	timeout      string
	mappings     string

	help    bool
	version bool
	verbose bool
)

func init() {
	flag.StringVar(&src, "src", "", "Source name. Such as `github/xiexianbin`")
	flag.StringVar(&srcToken, "src-token", "", "The app token which is used to list repo in source hub")
	flag.StringVar(&dst, "dst", "", "Destination name. Such as `gitee/xiexianbin`")
	flag.StringVar(&dstKey, "dst-key", "", "The private SSH key which is used to to push code in destination hub")
	flag.StringVar(&dstToken, "dst-token", "", "The app token which is used to create repo in destination hub")
	flag.StringVar(&accountType, "account-type", "user", "The account type. Such as org, user")
	flag.StringVar(&cloneStyle, "clone-style", "ssh", "The git clone style, https or ssh")
	flag.StringVar(&cachePath, "cache-path", "/github/workspace/git-mirrors-cache", "The path to cache the source repos code")
	flag.StringVar(&blackListStr, "black-list", "", "Height priority, the back list of mirror repo. like 'repo1,repo2,repo3'")
	flag.StringVar(&whiteListStr, "white-list", "", "Low priority, the white list of mirror repo. like 'repo1,repo2,repo3'")
	flag.BoolVar(&forceUpdate, "force-update", false, "Force to update the destination repo, use '-f' flag do 'git push'")
	flag.BoolVar(&debug, "debug", false, "Enable the debug flag to show detail log")
	flag.StringVar(&timeout, "timeout", "30m", "Set the timeout for every git command, eg. '600'=>600s, '30m'=>30 minute, '2h'=>2 hours")
	flag.StringVar(&mappings, "mappings", "", "The source repos mappings, such as 'A=>B, C=>CC', source repo name would be mapped follow the rule: A to B, C to CC. Mapping is not transitive")

	flag.BoolVar(&help, "h", false, "print this help")
	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&verbose, "V", false, "be verbose, debug model")

	flag.Parse()

	flag.Usage = func() {
		fmt.Println("Run: git-mirrors xxx\n" +
			"Usage:")
		flag.PrintDefaults()
	}
}

func showVersion() {
	logger.Print("v0.1.0")
}

func main() {
	if help == true {
		flag.Usage()
		return
	}

	if version == true {
		showVersion()
	}

	if verbose == true || debug == true {
		logger.SetLogLevel(logger.DEBUG)
	}

}
