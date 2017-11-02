#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$( cd "$( dirname "${DIR}" )" && pwd )"
REV="$( env GIT_WORK_TREE=${ROOT} git describe --tags --always --dirty )"

make -C ${ROOT} container
docker run -it artefactual/rdss-siegfried-service-amd64:${REV}
