#!/usr/bin/env bash
set -euo pipefail

# Collect metadata
DATETIME_TZ=$(date '+%Y-%m-%d %H:%M:%S %Z')
FILENAME_TS=$(date '+%Y-%m-%d_%H-%M-%S')

if command -v git >/dev/null 2>&1 && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
	REPO_ROOT=$(git rev-parse --show-toplevel)
	REPO_NAME=$(basename "$REPO_ROOT")
	GIT_BRANCH=$(git branch --show-current 2>/dev/null || git rev-parse --abbrev-ref HEAD)
	GIT_COMMIT=$(git rev-parse HEAD)
	GIT_USERNAME=$(git config user.name 2>/dev/null || echo "unknown")
else
	REPO_ROOT=""
	REPO_NAME=""
	GIT_BRANCH=""
	GIT_COMMIT=""
	GIT_USERNAME=""
fi

# Print similar to the individual command outputs
echo "Current Date/Time (TZ): $DATETIME_TZ"
[ -n "$GIT_USERNAME" ] && echo "Git Username: $GIT_USERNAME place document in thoughts/$GIT_USERNAME (not in thoughts/shared)"
[ -n "$GIT_COMMIT" ] && echo "Current Git Commit Hash: $GIT_COMMIT"
[ -n "$GIT_BRANCH" ] && echo "Current Branch Name: $GIT_BRANCH"
[ -n "$REPO_NAME" ] && echo "Repository Name: $REPO_NAME"
echo "Timestamp For Filename: $FILENAME_TS"
