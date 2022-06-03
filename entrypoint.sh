#!/bin/bash
set -e

DEBUG="${INPUT_DEBUG}"

if [[ X"$DEBUG" == X"true" ]]; then
  set -x
  DEBUG="true"
else
  DEBUG="false"
fi

FORCE_UPDATE="${INPUT_FORCE_UPDATE}"
if [[ X"$FORCE_UPDATE" == X"true" ]]; then
  FORCE_UPDATE="true"
else
  FORCE_UPDATE="false"
fi

echo "## Check User ##################"
whoami

echo "## Check Package Version ##################"
bash --version
git version
git lfs version
git-mirrors -v

echo "## Init Git Config ##################"
git config --global --add safe.directory /github/workspace/${PUBLISH_DIR}

echo "## Setup Deploy keys ##################"
[ -d /root/.ssh ] || mkdir /root/.ssh
if [ X"$INPUTE_SSH_KEYSCANS" = X"" ]; then
  INPUTE_SSH_KEYSCANS="github.com,gitee.com"
fi

KNOWN_HOSTS="/root/.ssh/known_hosts"
GIT_HOST_ARRAY=(${INPUTE_SSH_KEYSCANS//,/ })
for host in ${GIT_HOST_ARRAY[@]}; do
  # ssh-keyscan -t rsa $host >> ${KNOWN_HOSTS}
  # ssh-keyscan -t ecdsa $host >> ${KNOWN_HOSTS}
  ssh-keyscan $host >> ${KNOWN_HOSTS}
done
cat /root/.ssh/known_hosts
export SSH_KNOWN_HOSTS="${KNOWN_HOSTS}"

DST_KEY=""
if [ X"$INPUT_DST_KEY" = X"" ]; then
  echo "## Skip ssh key deploy ##################"
else
  DST_KEY="/root/.ssh/git_key"
  echo "${INPUT_DST_KEY}" > ${DST_KEY}
  chmod 400 ${DST_KEY}
  ls -lhart ${DST_KEY}
fi

echo "## begin sync ##################"

git-mirrors \
  --src "${INPUT_SRC}" \
  --src-token "${INPUT_SRC_TOKEN}" \
  --dst "${INPUT_DST}" \
  --dst-key "${DST_KEY}" \
  --dst-token "${INPUT_DST_TOKEN}" \
  --account-type "${INPUT_ACCOUNT_TYPE}" \
  --clone-style "${INPUT_CLONE_STYLE}" \
  --cache-path "${INPUT_CACHE_PATH}" \
  --black-list "${INPUT_BLACK_LIST}" \
  --white-list "${INPUT_WHITE_LIST}" \
  --force-update="${FORCE_UPDATE}" \
  --debug="${DEBUG}" \
  --timeout "${INPUT_TIMEOUT}" \
  --mappings "${INPUT_MAPPINGS}"

echo "## Done. ##################"
