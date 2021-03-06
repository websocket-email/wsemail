#! /bin/sh

set -e
set -u

cat <<EOF > README.md
# wsemail

This repository is the official cli client for [websocket.email](https://websocket.email).
This command lets you subscribe to email addresses at the websocket.email domain and perform actions
when they arrive.

### Installing and using

The client requires an API token, which can be aquired for free from the website.

To build and install the development command line client into \$GOPATH/bin run:

\`\`\`
go get github.com/websocket-email/wsemail
\`\`\`

To get a prebuilt version of the cli client download one from the [releases page](https://github.com/websocket-email/wsemail/releases).

EOF

echo "## usage" >> README.md
echo '```' >> README.md
go build
export PATH=`pwd`:$PATH
set +e
wsemail -help 2>> README.md
set -e
echo '```' >> README.md

echo <<EOF >> README.md
## License

See [LICENSE.md](LICENSE.md) file for details
EOF