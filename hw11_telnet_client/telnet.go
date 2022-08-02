package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

// TelnetClient represents simple telnet client behaviour.
type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

// NewTelnetClient creates new instance of TelnetClient.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Connect open connection to specified host and port.
func (t *telnetClient) Connect() error {
	var err error

	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return errors.Wrapf(err, "cannot connect to %s", t.address)
	}

	fmt.Fprintf(os.Stderr, "... connected to %s\n", t.address)

	return nil
}

// Send starts sending data to given "out" param.
func (t *telnetClient) Send() error {
	var err error

	scanner := bufio.NewScanner(t.in)
	for {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "... eof\n")
			break
		}

		_, err = fmt.Fprintf(t.conn, "%s\n", scanner.Text())
		if err != nil {
			break
		}
	}

	return errors.Wrapf(err, "get error when sending")
}

// Receive starts receiving data from "in" param.
func (t *telnetClient) Receive() error {
	var err error

	scanner := bufio.NewScanner(t.conn)
	for {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "... connection was closed by peer\n")
			break
		}

		_, err = fmt.Fprintf(t.out, "%s\n", scanner.Text())
		if err != nil {
			break
		}
	}

	return errors.Wrapf(err, "get error when receiving")
}

// Close client.
func (t *telnetClient) Close() error {
	if t.conn == nil {
		return nil
	}

	err := t.conn.Close()
	if err != nil {
		return errors.Wrapf(err, "get error when closing connect")
	}

	t.conn = nil

	return nil
}
