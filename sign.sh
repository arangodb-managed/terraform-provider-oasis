#!/bin/bash

gpg-agent --daemon --default-cache-ttl 7200 || true
echo "${GPG_PRIVATE_KEY}" | gpg --import --batch --no-tty
echo "hello world" > temp.txt
gpg --detach-sig --yes -v --output=/dev/null --pinentry-mode loopback --passphrase "${PASSPHRASE}" temp.txt
cd assets ; gpg --detach-sign ${PROJECT}-${VERSION}_SHA256SUMS
rm ../temp.txt
