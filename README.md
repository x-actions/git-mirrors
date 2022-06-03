# git-mirrors

[![build-test](https://github.com/x-actions/git-mirrors/actions/workflows/workflow.yaml/badge.svg)](https://github.com/x-actions/git-mirrors/actions/workflows/workflow.yaml)
[![GoDoc](https://godoc.org/github.com/x-actions/git-mirrors?status.svg)](https://pkg.go.dev/github.com/x-actions/git-mirrors)
[![Go Report Card](https://goreportcard.com/badge/github.com/x-actions/git-mirrors)](https://goreportcard.com/report/github.com/x-actions/git-mirrors)

a tools Mirrors Code from Github to Gitee.

## Feature

- Support Private/Public Organization/User 's Repos sync
- Support The Sample Git Provides sync, like `github.com/xiexianbin/test` to `github.com/x-actions/test`

## parameters

兼容 `Yikun/hub-mirror-action`

### Must

- `src` `github/<name>` name 可以是 user name 或 org name, eg: `github/xiexianbin`
- `src_token` 源的 API tokens（扩展参数：为支持同步私有仓库），支持 [Gitee](https://gitee.com/profile/personal_access_tokens)、[Github](https://github.com/settings/tokens),配置 `${{ secrets.GITHUB_TOKEN }}` 在 Github Action 中自动注入 token。
- `dst` `gitee/<name>` name 可以是 user name 或 org name, eg: `gitee/xiexianbin`
- `dst_key` 目的端和源端的 ssh private key
- `dst_token` 创建仓库的API tokens，支持[Gitee](https://gitee.com/profile/personal_access_tokens)、[Github](https://github.com/settings/tokens)

### optional

- `account_type` org(Organization) or user, default is user
- `src_account_type` 默认为account_type，源账户类型，可以设置为org（组织）或者user（用户）。
- `dst_account_type` 默认为account_type，目的账户类型，可以设置为org（组织）或者user（用户）。
- `clone_style` just support ssh, and `dst_key` must configure both github and gitee
- `cache_path` 默认为''，将代码缓存在指定目录，用于与 [actions/cache](https://github.com/actions/cache)配合以加速镜像过程。
- `black_list` 默认为''，配置后，黑名单中的repos将不会被同步，如“repo1,repo2,repo3”。
- `white_list` 默认为''，配置后，仅同步白名单中的repos，如“repo1,repo2,repo3”。
- `force_update` 默认为`false`, 配置后，启用`git push -f`强制同步，**注意：开启后，会强制覆盖目的端仓库**。
- `debug` 默认为`false`, 配置后，启用debug开关，会显示所有执行命令。
- `timeout` 默认为'30m', 用于设置每个git命令的超时时间，'600'=>600s, '30m'=>30 mins, '1h'=>1 hours
- `mappings` 源仓库映射规则，比如'A=>B, C=>CC', A会被映射为B，C会映射为CC，映射不具有传递性。主要用于源和目的仓库名不同的镜像。
- `ssh_keyscans` 默认为 `github.com,gitee.com`（扩展参数）

## How to Use

- Github Action

Sample Use

```
      - name: git mirror
        uses: x-actions/git-mirrors@main
        with:
          src: github/${{ matrix.github }}
          src_token: ${{ secrets.GITHUB_TOKEN }}
          dst: gitee/${{ matrix.gitee }}
          dst_key: ${{ secrets.GITEE_PRIVATE_KEY }}
          dst_token: ${{ secrets.GITEE_TOKEN }}
          account_type: user
          cache_path: "/github/workspace/git-mirrors-cache"
          black_list: "openbilibili,test1"
          clone_style: ssh
```

all Params

```
      - name: git mirror
        uses: x-actions/git-mirrors@main
        with:
          src: github/${{ matrix.github }}
          src_token: ${{ secrets.GITHUB_TOKEN }}
          dst: gitee/${{ matrix.gitee }}
          dst_key: ${{ secrets.GITEE_PRIVATE_KEY }}
          dst_token: ${{ secrets.GITEE_TOKEN }}
          account_type: user
          # src_account_type: org
          # dst_account_type: org
          cache_path: "/github/workspace/git-mirrors-cache"
          black_list: "openbilibili,test1"
          white_list: "w1,w2"
          clone_style: ssh
          force_update: false
          debug: true
          timeout: 30m
          mappings: "A=>B, C=>CC"
```

- command line

```
# download
curl -Lfs -o git-mirrors https://github.com/x-actions/git-mirrors/releases/latest/download/git-mirrors-{linux|darwin|windows}
chmod +x git-mirrors

# help
./git-mirrors -h

# demo
git-mirrors \
  --src "github/estack" \
  --src-token "${GITHUB_TOKEN}" \
  --dst "gitee/e-stack" \
  --dst-key "" \
  --dst-token "${GITEE_TOKEN}" \
  --account-type "user" \
  --clone-style "ssh" \
  --cache-path "./temp/" \
  --black-list "" \
  --white-list "" \
  --force-update=true \
  --debug=true \
  --timeout "10m"
```

## FaQ

- ssh key err

```
clone git@github.com:xx/xx.git err: unknown error: ERROR: You're using an RSA key with SHA-1, which is no longer allowed. Please use a newer client or a different key type.
```

regenerate ssh key:

```
$ ssh-keygen -t ed25519 -C "your_email@example.com"

# or
$ ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

more info go to https://github.blog/2021-09-01-improving-git-protocol-security-github/

## ref

- 采用兼容 [Yikun/hub-mirror-action](https://github.com/Yikun/hub-mirror-action) 的配置参数，因此不可避免的参考其实现，在此表示感谢。
