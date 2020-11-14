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
	cfg.AddKey("crooks", "/home/crooks/.ssh/openwrt")
	cfg.AddKey("crooks", "/home/crooks/.ssh/id_nopass")
	var b bytes.Buffer
	for _, hostName := range hostNames {
		client, err := cfg.Auth(hostName)
		if err != nil {
			log.Println(err)
			continue
		}
		b, err = Cmd(client, "cat /etc/passwd")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		b, err = Cmd(client, "cat /etc/group")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		client.Close()
	}
}
