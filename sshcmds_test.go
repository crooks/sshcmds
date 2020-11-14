package sshcmds

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

var cfg config

func TestClient(t *testing.T) {
	cfg = *newConfig()
	hostNames := []string{"bingo.mixmin.net"}
	cfg.addKey("crooks", "/home/crooks/.ssh/openwrt")
	cfg.addKey("crooks", "/home/crooks/.ssh/id_nopass")
	var b bytes.Buffer
	for _, hostName := range hostNames {
		client, err := cfg.sshClient(hostName)
		if err != nil {
			log.Println(err)
			continue
		}
		b, err = sshCmd(client, "cat /etc/passwd")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		b, err = sshCmd(client, "cat /etc/group")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		client.Close()
	}
}
