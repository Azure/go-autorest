#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

set -euo pipefail
set -x

(
cd "${DIR}"
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' .
az storage blob upload \
	--account-name=colemickmsitest0 \
	--container=testapp \
	--name=example-msi-app \
	--file=./example-msi-app
)
