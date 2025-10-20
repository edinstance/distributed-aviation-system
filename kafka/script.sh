#!/bin/sh
set -e

echo "Waiting for Schema Registry..."
until curl -sf http://schema-registry:8081/subjects >/dev/null; do
  sleep 2
done
echo "Registry is ready"

if ! ls /schemas/*.avsc 1> /dev/null 2>&1; then
  echo "No .avsc files found in /schemas"
  exit 0
fi

for f in /schemas/*.avsc; do
  subject=$(basename "$f" .avsc)

  subject=$(basename "$f" .avsc)
  echo "Registering $subject"

  # Compact and escape the file contents
  ESCAPED=$(jq -c . < "$f" | sed 's/"/\\"/g')

  response=$(curl -sf -w "%{http_code}" -o /tmp/resp.txt \
    -X POST \
    -H "Content-Type: application/vnd.schemaregistry.v1+json" \
    --data "{\"schema\":\"$ESCAPED\"}" \
    "http://schema-registry:8081/subjects/${subject}/versions") || true

  if [ "$response" = "200" ] || [ "$response" = "201" ]; then
    echo "
    $subject registered"
  else
    echo "Failed ($response): $(cat /tmp/resp.txt)"
  fi
done