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
LABEL "Version"="v0.2.1"

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
# change VERSION when make a release, v1.0.0
ENV VERSION "v0.2.1"

RUN apk update && \
    apk add --no-cache git git-lfs bash wget curl openssh-client tree && \
    rm -rf /var/cache/apk/* && \
    cd /tmp/ && \
    wget https://github.com/x-actions/git-mirrors/releases/download/${VERSION}/git-mirrors-linux && \
    chmod +x /tmp/git-mirrors-linux && \
    mv /tmp/git-mirrors-linux /usr/local/bin/git-mirrors

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

WORKDIR /github/workspace
ENTRYPOINT ["/entrypoint.sh"]
