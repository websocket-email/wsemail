# wsemail

This repository is the official cli client for [websocket.email](https://websocket.email).
This command lets you subscribe to email addresses at the websocket.email domain and perform actions
when they arrive.

### Installing and using

The client requires an API token, which can be aquired for free from the website.

To build and install the development command line client into $GOPATH/bin run:

```
go get github.com/websocket-email/wsemail
```

To get a prebuilt version of the cli client download one from the [releases page](https://github.com/websocket-email/wsemail/releases).

## usage
```
A tool for interacting with https://websocket.email

Usage of wsemail:
  -api-token string
    	API token to authenticate with, can also be specified with the env variable WEBSOCKETEMAIL_TOKEN
  -for-address string
    	Subscribe to emails arriving at this email address
  -generate-address
    	Generate a random fake email, print to stdout and exit
  -n int
    	Wait for and print this many emails before exiting, less than or equal to zero waits forever (default 1)
  -timeout uint
    	Wait this many seconds for an email to arrive before giving up and terminating with an error, 0 for no timeout (default 60)
  -version
    	print the version then to stdout and exit

Examples:

  Generate a secure random email address:

    $ wsemail -generate-address

  Wait 10 seconds for a single email for john@websocket.email:

    $ wsemail -for-address john@websocket.email -timeout 10
```

