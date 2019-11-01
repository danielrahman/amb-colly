#!/usr/bin/env bash
set -exuo pipefail
IFS=$'\n\t'

PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/go/bin:/usr/local/go/bin"

if ping -q -c 1 -W 1 8.8.8.8 >/dev/null; then
  go run /Users/danielrahman/Sites/craness/amb-colly/transactions/transactions.go
fi
