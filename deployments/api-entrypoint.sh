#!/bin/sh
set -eu

if [ "$#" -gt 0 ]; then
	exec "$@"
fi

mkdir -p "$CATENA_GIT_ROOT"
chown catena:catena "$CATENA_GIT_ROOT"

exec su-exec catena /usr/local/bin/catena-api
