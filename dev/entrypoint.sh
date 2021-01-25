#!/usr/bin/env bash
set -e

if [[ "$1" = 'run' ]]; then
    chmod +x ./scripts/build.sh

    exec CompileDaemon \
    -build="./scripts/build.sh polygnosics" \
    -command="polygnosics ${@:2}" \
    -polling \
    -graceful-kill=true \
    -log-prefix=false
fi

exec "$@"