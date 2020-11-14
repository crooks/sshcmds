package sshcmds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type config struct {
	timeout    time.Duration
	sshConfigs []*ssh.ClientConfig
}

func newConfig() *config {
	return &config{
		timeout: setTimeout("10s"),
	}
}

func setTimeout(duration string) time.Duration {
	// Expect a string like 10s and convert to a time.Duration.
	timeout, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	return timeout
}

// publicKeyFile creates an SSH authentication method from a text file.
func publicKey(file string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(signer)
}

// To authenticate with the remote server you must pass at least one
// implementation of AuthMethod via the Auth field in ClientConfig,
// and provide a HostKeyCallback.
func (c *config) makeSSHConfig(userName, keyFile string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			publicKey(keyFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}
}

// addKey appends an SSH private key to the list of keys
func (c *config) addKey(userName, keyFile string) {
	c.sshConfigs = append(c.sshConfigs, c.makeSSHConfig(userName, keyFile))
}

// sshClient creates an SSH session and issues a single command.  The output
// from the command is returned as raw bytes.
func (c *config) sshClient(hostname string) (client *ssh.Client, err error) {
	hostport := fmt.Sprintf("%s:22", hostname)
	for _, sshConfig := range c.sshConfigs {
		client, err = ssh.Dial("tcp", hostport, sshConfig)
		// If err is nil, we successfully dialed with the sshKey.
		// We can stop iterating over keys and break out of the loop.
		if err == nil {
			log.Printf("Successfully dialed %s", hostname)
			break
		} else {
			log.Printf("Failed to dial %s: %s", hostname, err)
		}
	}
	return
}

func sshCmd(client *ssh.Client, cmd string) (b bytes.Buffer, err error) {
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
