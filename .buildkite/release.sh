#!/usr/bin/env sh
set -eu

apk add --no-cache curl git

# install goreleaser
GORELEASER_VERSION=v2.8.2
curl -fsSL "https://github.com/goreleaser/goreleaser/releases/download/${GORELEASER_VERSION}/goreleaser_Linux_x86_64.tar.gz" \
  | tar -xz -C /usr/local/bin goreleaser

goreleaser release --clean

ORG=geqo
API="https://api.buildkite.com/v2/packages/organizations/${ORG}/registries"

upload() {
  registry=$1
  file=$2
  echo "Uploading ${file} to ${registry}..."
  curl -fsSL \
    -H "Authorization: Bearer ${BUILDKITE_API_TOKEN}" \
    -F "file=@${file}" \
    "${API}/${registry}/packages"
}

for f in dist/*.deb; do upload togram-deb "$f"; done
for f in dist/*.rpm; do upload togram-rpm "$f"; done
for f in dist/*.apk; do upload togram-apk "$f"; done
