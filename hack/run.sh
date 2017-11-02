#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$( cd "$( dirname "${DIR}" )" && pwd )"

cd ${DIR}
go install .
rdss-siegfried-service -home=${DIR}/siegfried-data -sf=$(which sf)
