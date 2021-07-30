#!/bin/sh

if test -f .env; then
  . .env
else
  . .env.example
fi

json() {
  base64 -w 0 < "$1" | jq -Rn 'input | sub("[=]+$"; "") | { data: . }'
}

upload() {
  json $@ | curl -H 'content-type: application/json' -X POST "http://localhost:${PORT}/api/negative_image" -d @-
}

case "$1" in
  json|upload)
    $@
    ;;
  *)
    echo "Unknown command"
    exit 1
    ;;
esac
