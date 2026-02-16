#!/bin/sh

set -e

# Inject SIXTINEAI_ environment variables into index.html
JSON_STRING='window.configs = {'
FIRST=true

for var in $(env | grep '^SIXTINEAI_' | sed 's/=.*//'); do
  value=$(printenv "$var" | sed 's/"/\\"/g')
  if [ "$FIRST" = true ]; then
    FIRST=false
  else
    JSON_STRING="${JSON_STRING},"
  fi
  JSON_STRING="${JSON_STRING}\"${var}\":\"${value}\""
done

JSON_STRING="${JSON_STRING}}"

sed -i "s@// CONFIGURATIONS_PLACEHOLDER@${JSON_STRING}@" /usr/share/nginx/html/index.html

nginx -g "daemon off;"
