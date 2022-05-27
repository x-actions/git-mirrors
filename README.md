# git-mirrors

[![build-test](https://github.com/xiexianbin/go-actions-demo/actions/workflows/workflow.yaml/badge.svg)](https://github.com/xiexianbin/go-actions-demo/actions/workflows/workflow.yaml)
[![GoDoc](https://godoc.org/github.com/xiexianbin/go-actions-demo?status.svg)](https://pkg.go.dev/github.com/xiexianbin/go-actions-demo)

a tools Mirrors Code from Github to Gitee.

## parameters

### Must

- `src` `github/<name>` name 可以是 user name 或 org name, eg: `github/xiexianbin`
- `src_token` 源的API tokens，仅支持[Github](https://gitee.com/profile/personal_access_tokens)
- `dst` `gitee/<name>` name 可以是 user name 或 org name, eg: `gitee/xiexianbin`
- `dst_key` 目的端和源端的 ssh private key
- `dst_token` 创建仓库的API tokens，仅支持[Gitee](https://gitee.com/profile/personal_access_tokens)

### optional

- `account_type` org(Organization) or user, default is user
- `clone_style` just support ssh, and `dst_key` must configure both github and gitee
- `cache_path` 默认为''，将代码缓存在指定目录，用于与 [actions/cache](https://github.com/actions/cache)配合以加速镜像过程。
- `black_list` 默认为''，配置后，黑名单中的repos将不会被同步，如“repo1,repo2,repo3”。
- `white_list` 默认为''，配置后，仅同步白名单中的repos，如“repo1,repo2,repo3”。
- `force_update` 默认为`false`, 配置后，启用git push -f强制同步，**注意：开启后，会强制覆盖目的端仓库**。
- `debug` 默认为`false`, 配置后，启用debug开关，会显示所有执行命令。
- `timeout` 默认为'30m', 用于设置每个git命令的超时时间，'600'=>600s, '30m'=>30 mins, '1h'=>1 hours
- `mappings` 源仓库映射规则，比如'A=>B, C=>CC', A会被映射为B，C会映射为CC，映射不具有传递性。主要用于源和目的仓库名不同的镜像。

## How to Use

- Github Action

```
```

- command line

```
# download
curl -Lfs -o git-mirrors https://github.com/xiexianbin/go-actions-demo/releases/latest/download/git-mirrors-{linux|darwin|windows}
chmod +x git-mirrors

# help
./git-mirrors -h
```

## ref

- 采用兼容 [Yikun/hub-mirror-action](https://github.com/Yikun/hub-mirror-action) 的配置参数，因此不可避免的参考其实现，在此表示感谢。
