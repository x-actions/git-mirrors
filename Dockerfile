FROM alpine:latest

ENV REFRESHED_AT 2022-05-25

LABEL "com.github.actions.name"="go-git-mirrors"
LABEL "com.github.actions.description"="Mirrors Code from github to gitee."
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="green"
LABEL "repository"="http://github.com/x-actions/git-mirrors"
LABEL "homepage"="http://github.com/x-actions/git-mirrors"
LABEL "maintainer"="xiexianbin<me@xiexianbin.cn>"

LABEL "Name"="git-mirrors"
LABEL "Version"="v1.1.0"

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8

RUN apk update && apk add --no-cache git git-lfs bash wget curl openssh-client tree && rm -rf /var/cache/apk/*

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

WORKDIR /github/workspace
ENTRYPOINT ["/entrypoint.sh"]
