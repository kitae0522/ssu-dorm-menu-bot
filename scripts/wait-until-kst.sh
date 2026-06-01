#!/usr/bin/env bash
set -euo pipefail

target="${1:?usage: wait-until-kst.sh HH:MM}"

while true; do
  now="$(TZ=Asia/Seoul date +%H:%M)"
  if [[ "$now" > "$target" || "$now" == "$target" ]]; then
    echo "KST target reached: now=${now}, target=${target}"
    exit 0
  fi
  echo "Waiting for ${target} KST. now=${now}"
  sleep 10
done
