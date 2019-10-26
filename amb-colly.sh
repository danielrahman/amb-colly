#!/bin/bash
if ping -q -c 1 -W 1 8.8.8.8 >/dev/null; then
  go run /Users/danielrahman/Sites/craness/amb-colly/main.go
  go run /Users/danielrahman/Sites/craness/amb-colly/ambassadors/ambassadors.go
fi
