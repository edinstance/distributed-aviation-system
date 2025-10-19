#!/bin/bash
set -euo pipefail

OS_URL="http://opensearch-node-1:9200"
DASH_URL="http://opensearch-dashboards:5601"

wait_for() {
  local URL=$1 NAME=$2
  echo "⏳ Waiting for $NAME..."
  until curl -fs "$URL" >/dev/null 2>&1; do
    sleep 5
    echo "  still waiting..."
  done
  echo "✅ $NAME ready."
}

wait_for "$OS_URL/_cluster/health" "OpenSearch"
wait_for "$DASH_URL/api/status" "Dashboards"

# ---------------- INDEX CREATION -----------------
for file in /init/indexes/*.json; do
  name=$(basename "$file" .json)
  echo "🗃️  Creating index '$name'..."
  curl -sf -X PUT "$OS_URL/$name" \
    -H "Content-Type: application/json" \
    -d @"$file" \
    || echo "⚠️  could not create index $name"
done
echo "✅ Index creation finished"

# ---------------- INDEX PATTERN -----------------
INDEX_PATTERN_ID="flights"
echo "📋 Creating index pattern '$INDEX_PATTERN_ID'..."

# Delete previous if exists
EXISTING_PATTERN=$(curl -s \
  "$DASH_URL/api/saved_objects/_find?type=index-pattern&search_fields=title&search=$INDEX_PATTERN_ID" \
  -H "osd-xsrf: true" | jq -r '.saved_objects[0].id // empty')

if [ -n "$EXISTING_PATTERN" ]; then
  echo "🗑️  Deleting existing index pattern: $EXISTING_PATTERN"
  curl -s -X DELETE \
    "$DASH_URL/api/saved_objects/index-pattern/$EXISTING_PATTERN" \
    -H "osd-xsrf: true" >/dev/null || true
fi

# Create new pattern with fixed ID "flights"
PATTERN_RESP=$(curl -s -X POST \
  "$DASH_URL/api/saved_objects/index-pattern/$INDEX_PATTERN_ID" \
  -H "Content-Type: application/json" \
  -H "osd-xsrf: true" \
  -d '{"attributes":{"title":"flights","timeFieldName":"indexedAt"}}')

if echo "$PATTERN_RESP" | jq -e 'has("id")' >/dev/null; then
  echo "✅ Index pattern created with fixed ID: $INDEX_PATTERN_ID"
else
  echo "⚠️  Failed to create index pattern. Response: $PATTERN_RESP"
  exit 1
fi

# Helper to create an object and print ID
create_object() {
  local TYPE="$1" FILE="$2"
  local RESP
  RESP="$(curl -s -X POST "$DASH_URL/api/saved_objects/$TYPE" \
    -H "Content-Type: application/json" -H "osd-xsrf: true" \
    -d @"$FILE")"
  echo "$RESP" | jq -r '.id // empty'
}

mkdir -p /tmp/os_ids

# ---------------- VISUALIZATIONS -----------------
echo "📊 Creating visualizations..."
for file in /init/vizualizations/*.json; do
  vis_name=$(basename "$file" .json)
  id=$(create_object visualization "$file")
  echo "$vis_name=$id" >>/tmp/os_ids/vis.list
  echo "  -> $vis_name : $id"
done

# ---------------- SEARCHES -----------------
echo "🔍 Creating searches..."
for file in /init/searches/*.json; do
  search_name=$(basename "$file" .json)
  id=$(create_object search "$file")
  echo "$search_name=$id" >>/tmp/os_ids/vis.list
  echo "  -> $search_name : $id"
done

# ---------------- DASHBOARD -----------------
template=/init/dashboards/flight-analytics.json.template
filled=/tmp/dashboard.json

cp "$template" "$filled"
while IFS= read -r line; do
  key=$(echo "$line" | cut -d= -f1)
  val=$(echo "$line" | cut -d= -f2)
  sed -i "s/{{${key^^}_ID}}/$val/g" "$filled"
done </tmp/os_ids/vis.list

# Replace index pattern placeholder (keeps name 'flights')
sed -i "s/{{INDEX_PATTERN_ID}}/$INDEX_PATTERN_ID/g" "$filled"

echo "🎨 Creating dashboard..."
dash_id=$(create_object dashboard "$filled")
echo "✅ Dashboard created with id $dash_id"
echo "🌐 http://localhost:5601/app/dashboards#/view/$dash_id"