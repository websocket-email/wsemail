# websocketemail-go

This repository is the official go library for [websocket.email](https://websocket.email), 
a tool that makes testing email based workflows easy. This code lets you subscribe to email addresses at the websocket.email domain from go code, and then do actions when they arrive.

Some API calls require an API token, which can be aquired for free from the website.

## Example

```
apiToken := "YOU_API_TOKEN"
forAddress := MustGenerateEmailAddress()

ch, cleanup, err := WaitForEmail(apiToken, forAddress)
if err != nil {
    panic(err)
}
defer cleanup()

select {
case email, ok := <-ch:
    if !ok {
        panic("unable to wait for email!")
    }
    fmt.Printf("Got email: %#v", email)
case <-time.After(20 * time.Second):
    panic("email was not recieved after 20 seconds!")
}
```

## Installing and using

To use the go library, follow the example provided here [this file](https://github.com/websocket-email/wsemail)

Library documentation can be found [here](https://godoc.org/github.com/websocket-email/websocketemail-go).

## Running the tests

Get a valid API token from [websocket.email](https://websocket.email), change to the project directoryin your go path and run:

```
export WEBSOCKETEMAIL_TOKEN="$YOUR_TOKEN_HERE"
go test
```

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## License

See [LICENSE.md](LICENSE) file for details

