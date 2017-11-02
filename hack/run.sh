#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$( cd "$( dirname "${DIR}" )" && pwd )"

cd ${ROOT}
go install .
rdss-siegfried-service -home=${ROOT}/siegfried-data -sf=$(which sf)
