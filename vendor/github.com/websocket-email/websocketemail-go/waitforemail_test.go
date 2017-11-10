package websocketemail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"testing"
	"time"
)

func GetTestToken(t *testing.T) string {
	tok := os.Getenv("WEBSOCKETEMAIL_TOKEN")
	if tok == "" {
		t.Fatal("please set WEBSOCKETEMAIL_TOKEN env variable")
	}
	return tok
}

func SendEmail(t *testing.T, from, to, subject, content string) {

	c, err := smtp.Dial("websocket.email:25")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	host := "websocket.email"
	config := &tls.Config{ServerName: host}
	err = c.StartTLS(config)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Mail(from)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Rcpt(to)
	if err != nil {
		t.Fatal(err)
	}
	w, err := c.Data()
	if err != nil {
		t.Fatal(err)
	}

	mail := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, content)
	_, err = w.Write([]byte(mail))
	if err != nil {
		w.Close()
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Quit()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWaitForEmail(t *testing.T) {
	tok := GetTestToken(t)
	from := MustGenerateEmailAddress()
	forAddress := MustGenerateEmailAddress()

	ch, cleanup, err := WaitForEmail(tok, forAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	SendEmail(t, from, forAddress, "hi", "hello world")

	select {
	case email, ok := <-ch:
		if !ok {
			t.Fatal("unable to wait for email!")
		}
		if email.From != from {
			t.Fatalf("bad email from: %#v", email)
		}
		if email.To != forAddress {
			t.Fatalf("bad email to: %#v", email)
		}
		if email.Subject != "hi" {
			t.Fatalf("bad email subject: %#v", email)
		}
		if strings.Trim(email.Body, " \n") != "hello world" {
			t.Fatalf("bad email body: %#v", email)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("email was not recieved after 20 seconds!")
	}
}

// Wait 20 seconds for an email and then give up.
func ExampleWaitForEmail() {
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
}
