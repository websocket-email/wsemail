package websocketemail

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrTooManyConcurrentRequests            = errors.New("Too many concurrent requests")
	ErrUnauthorized                         = errors.New("Unauthorized, token invalid or expired")
	ErrServerAtCapacityOrDownForMaintenance = errors.New("Server is at capcity or down for maintenance")
)

type ParsedEmail struct {
	To          string
	From        string
	Subject     string
	ContentType string
	Body        string
}

// Connect to websocket.email over a secure connection and subscribe to emails received by address 'to'.
// 'forAddress' must be an email from the websocket.email domain or an error is returned.
//
// On success WaitForEmail returns a channel which will get parsed emails from the remote end, and
// a function which can be called to cleanup the worker goroutine. The returned channel
// shouldn't be closed by the caller.
//
// On failure may return ErrTooManyConcurrentRequests, ErrUnauthorized, ErrServerAtCapacityOrDownForMaintenance
// or a generic error.
func WaitForEmail(tok, forAddress string) (<-chan ParsedEmail, func(), error) {
	if !strings.HasSuffix(forAddress, "@websocket.email") {
		return nil, nil, errors.New("email must end with @websocket.email")
	}

	ch := make(chan ParsedEmail)
	done := make(chan struct{})

	u := url.URL{Scheme: "wss", Host: "api.websocket.email", Path: "/v1/subscribe"}

	hdr := make(http.Header)
	hdr.Add("Token", tok)
	hdr.Add("SubscribeTo", forAddress)

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), hdr)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusTooManyRequests:
				return nil, nil, ErrTooManyConcurrentRequests
			case http.StatusUnauthorized:
				return nil, nil, ErrUnauthorized
			case http.StatusServiceUnavailable:
				return nil, nil, ErrServerAtCapacityOrDownForMaintenance
			}
		}
		return nil, nil, err
	}

	go func() {
		defer close(ch)
		for {
			email := ParsedEmail{}
			err := c.ReadJSON(&email)
			if err != nil {
				return
			}
			select {
			case <-done:
				return
			case ch <- email:
			}
		}
	}()

	cleanup := func() {
		_ = c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Time{})
		_ = c.Close()
		close(done)
	}

	return ch, cleanup, nil
}
