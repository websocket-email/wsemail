package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/websocket-email/websocketemail-go"
)

var (
	generateAddress = flag.Bool("generate-address", false, "Generate a random fake email, print to stdout and exit")
	apiToken        = flag.String("api-token", "", "API token to authenticate with, can also be specified with the env variable WEBSOCKETEMAIL_TOKEN")
	fromAddress     = flag.String("for-address", "", "Subscribe to emails from this address")
	numEmails       = flag.Int64("n", 1, "Wait for and print this many emails before exiting, less than or equal to zero waits forever")
	timeoutSeconds  = flag.Uint64("timeout", 60, "Wait this many seconds for an email to arrive before giving up and terminating with an error, 0 for no timeout")
)

func usage() {
	fmt.Fprintln(os.Stderr, "A tool for interacting with https://websocket.email")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Generate a secure random email address:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "    $ %s -generate-address\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Wait 10 seconds for a single email for john@websocket.email:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "    $ %s -for-address john@websocket.email -timeout 10\n", os.Args[0])

	os.Exit(1)
}

func main() {

	flag.Usage = usage

	flag.Parse()

	if *generateAddress {
		addr := websocketemail.MustGenerateEmailAddress()
		_, err := fmt.Println(addr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing output: %s\n", err)
			os.Exit(1)
		}
		return
	}

	if *fromAddress == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-from-address or -generate-address required, see -help for details")
		os.Exit(1)
	}

	if *apiToken == "" {
		*apiToken = os.Getenv("WEBSOCKETEMAIL_TOKEN")
		if *apiToken == "" {
			_, _ = fmt.Fprintln(os.Stderr, "-api-token or env variable WEBSOCKETEMAIL_TOKEN required, see -help for details")
			os.Exit(1)
		}
	}

	ch, cleanup, err := websocketemail.WaitForEmail(*apiToken, *fromAddress)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error subscribing to email: %s\n", err)
		os.Exit(1)
	}
	defer cleanup()

	for {
		timer := time.NewTimer(time.Duration(*timeoutSeconds) * time.Second)
		timeoutChan := timer.C
		if *timeoutSeconds == 0 {
			timeoutChan = make(chan time.Time)
		}
		select {
		case email, ok := <-ch:
			if !ok {
				_, _ = fmt.Fprintf(os.Stderr, "an error occured while waiting for email\n")
				os.Exit(1)
			}
			buf, err := json.Marshal(&email)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error marshalling output: %s\n", err)
				os.Exit(1)
			}
			_, err = fmt.Println(string(buf))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error writing output: %s\n", err)
				os.Exit(1)
			}
		case <-timeoutChan:
			_, _ = fmt.Fprintln(os.Stderr, "no emails arrived before timeout")
			os.Exit(2)
		}
		timer.Stop()

		*numEmails -= 1
		if *numEmails == 0 {
			break
		}
	}
}
