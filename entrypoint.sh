#!/bin/bash
set -e

DEBUG="${INPUT_DEBUG}"

if [[ "$DEBUG" == "true" ]]; then
  set -x
fi

echo "## Check Package Version ##################"
bash --version
git version
git lfs version

echo "## Init Git Config ##################"
git config --global --add safe.directory /github/workspace/${PUBLISH_DIR}

echo "## Setup Deploy keys ##################"
mkdir /root/.ssh && ssh-keyscan -t rsa github.com > /root/.ssh/known_hosts

if [ X"$INPUT_DST_KEY" = X"" ]; then
  echo "## Skip ssh key deploy ##################"
else
  echo ${INPUT_DST_KEY} > /root/.ssh/id_rsa
  chmod 400 /root/.ssh/id_rsa && \
  ls -lhart /root/.ssh/id_rsa
fi

echo "## begin sync ##################"

git-mirrors \
  --src "${INPUT_SRC}" \
  --src-token "${INPUT_SRC_TOKEN}" \
  --dst "${INPUT_DST}" \
  --dst-key "${INPUT_DST_KEY}" \
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
