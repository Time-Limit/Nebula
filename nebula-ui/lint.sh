#!/usr/bin/env bash
if test "$BASH" = "" || "$BASH" -uc "a=();true \"\${a[@]}\"" 2>/dev/null; then
    # Bash 4.4, Zsh
    set -euo pipefail
else
    # Bash 4.3 and older chokes on empty arrays with set -u.
    set -eo pipefail
fi
cd "$(dirname "$0")"
yarn
find src -type f \( -name "*.js" -or -name "*.vue" \) | xargs ./node_modules/.bin/eslint
