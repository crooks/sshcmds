package sshcmds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

// Config holds details of each key config
type Config struct {
	timeout    time.Duration
	sshConfigs []*ssh.ClientConfig
}

// NewConfig creates a new instance of the Config struct
func NewConfig() *Config {
	return &Config{
		timeout: setDefaultTimeout("10s"),
	}
}

func setDefaultTimeout(duration string) time.Duration {
	// Expect a string like 10s and convert to a time.Duration.
	timeout, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	return timeout
}

// publicKeyFile creates an SSH authentication method from a text file.
func publicKey(file string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// To authenticate with the remote server you must pass at least one
// implementation of AuthMethod via the Auth field in ClientConfig,
// and provide a HostKeyCallback.
func (c *Config) makeSSHConfig(userName, keyFile string) (*ssh.ClientConfig, error) {
	pk, err := publicKey(keyFile)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			pk,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}, nil
}

// AddKey appends an SSH private key to the list of keys
func (c *Config) AddKey(userName, keyFile string) (err error) {
	newConfig, err := c.makeSSHConfig(userName, keyFile)
	if err != nil {
		return
	}
	c.sshConfigs = append(c.sshConfigs, newConfig)
	return
}

// Auth returns an ssh.Client struct after successfully authenticating with a key.
func (c *Config) Auth(hostname string) (client *ssh.Client, err error) {
	hostport := fmt.Sprintf("%s:22", hostname)
	for _, sshConfig := range c.sshConfigs {
		client, err = ssh.Dial("tcp", hostport, sshConfig)
		// If err is nil, we successfully dialed with the sshKey.
		// We can stop iterating over keys and break out of the loop.
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("Failed to authenticate with %d keys", len(c.sshConfigs))
	return
}

// Cmd runs a single command against a previously authenticated session and returns the output as a Byte buffer.
func (c *Config) Cmd(client *ssh.Client, cmd string) (b bytes.Buffer, err error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		err = fmt.Errorf("Failed to create session: %s", err)
		return
	}

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	session.Stdout = &b
	if err = session.Run(cmd); err != nil {
		err = fmt.Errorf("Failed to run: %s", err.Error())
	}
	return
}
