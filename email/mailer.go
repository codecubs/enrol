package email

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"os"
)

type Emailer interface {
	Message() string
	FromAddr() string
	ToAddr() string
}

func Send(b Emailer) error {
	// Connect to the SMTP Server
	servername := os.Getenv("smtpserver")

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", os.Getenv("adminemail"), os.Getenv("password"), host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	defer c.Quit()
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(b.FromAddr()); err != nil {
		return err
	}

	if err = c.Rcpt(b.ToAddr()); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(b.Message()))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
