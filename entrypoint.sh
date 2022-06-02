#!/bin/bash
set -e

DEBUG="${INPUT_DEBUG}"

if [[ X"$DEBUG" == X"true" ]]; then
  set -x
fi

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
GIT_HOST_ARRAY=(${INPUTE_SSH_KEYSCANS//,/ })
for host in ${GIT_HOST_ARRAY[@]}; do
  # ssh-keyscan -t rsa $host >> /root/.ssh/known_hosts
  # ssh-keyscan -t ecdsa $host >> /root/.ssh/known_hosts
  ssh-keyscan $host >> /root/.ssh/known_hosts
done

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
  --force-update "${INPUT_FORCE_UPDATE}" \
  --debug "${INPUT_DEBUG}" \
  --timeout "${INPUT_TIMEOUT}" \
  --mappings "${INPUT_MAPPINGS}"

echo "## Done. ##################"
