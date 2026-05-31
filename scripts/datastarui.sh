#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
DATASTARUI_ROOT=$(cd -- "$SCRIPT_DIR/.." && pwd)
BIN_DIR="$DATASTARUI_ROOT/bin"
LAUNCHER="$BIN_DIR/datastarui"
LAUNCHER_HASH_FILE="$BIN_DIR/datastarui-launcher.hash"

launcher_hash() {
  (
    cd "$DATASTARUI_ROOT"
    {
      find cmd/datastarui-launcher -type f -name '*.go' -print
      printf '%s\n' go.mod go.sum
    } | sort | while IFS= read -r path; do
      if [[ -f "$path" ]]; then
        sha256sum "$path"
      fi
    done | sha256sum | awk '{print $1}'
  )
}

current_hash=$(launcher_hash)
previous_hash=""
if [[ -f "$LAUNCHER_HASH_FILE" ]]; then
  previous_hash=$(cat "$LAUNCHER_HASH_FILE")
fi

if [[ ! -x "$LAUNCHER" || "$current_hash" != "$previous_hash" ]]; then
  mkdir -p "$BIN_DIR"
  echo "Building DatastarUI launcher: $LAUNCHER" >&2
  (
    cd "$DATASTARUI_ROOT"
    go build -o "$LAUNCHER.tmp" ./cmd/datastarui-launcher
  )
  mv "$LAUNCHER.tmp" "$LAUNCHER"
  printf '%s\n' "$current_hash" > "$LAUNCHER_HASH_FILE"
fi

exec "$LAUNCHER" "$@"
