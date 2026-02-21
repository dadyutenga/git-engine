package ssh

import (
	"fmt"
	"io"
	"os"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

// Config holds SSH connection configuration.
type Config struct {
	Host           string
	User           string
	Port           int
	Password       string
	PrivateKeyPath string
}

// Client wraps an SSH client connection.
type Client struct {
	client *gossh.Client
}

// New creates and connects an SSH client.
func New(cfg Config) (*Client, error) {
	auths := []gossh.AuthMethod{}
	if cfg.PrivateKeyPath != "" {
		key, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read private key: %w", err)
		}
		signer, err := gossh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		auths = append(auths, gossh.PublicKeys(signer))
	}
	if cfg.Password != "" {
		auths = append(auths, gossh.Password(cfg.Password))
	}

	if len(auths) == 0 {
		return nil, fmt.Errorf("no SSH authentication method provided")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	clientConfig := &gossh.ClientConfig{
		User:            cfg.User,
		Auth:            auths,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), // nolint:gosec
		Timeout:         15 * time.Second,
	}

	cli, err := gossh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("dial ssh: %w", err)
	}
	return &Client{client: cli}, nil
}

// Run executes a command and returns the combined output.
func (c *Client) Run(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

// RunStream executes a command streaming stdout to writer.
func (c *Client) RunStream(command string, writer io.Writer) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = writer
	session.Stderr = writer
	return session.Run(command)
}

// Close terminates the SSH connection.
func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}
