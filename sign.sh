#!/bin/bash
set -e

gpg-agent --daemon --default-cache-ttl 7200
# Multiline key in an envvar gets mangled.
echo "$GPG_PRIVATE_KEY_BASE64" | base64 -d | gpg --import --batch --no-tty
echo "hello world" > temp.txt
gpg --detach-sig --yes -v --output=/dev/null --pinentry-mode loopback --passphrase "${PASSPHRASE}" temp.txt
cd assets ; gpg --detach-sign ${PROJECT}-${VERSION}_SHA256SUMS
rm ../temp.txt
