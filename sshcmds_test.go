package sshcmds

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	cfg := *NewConfig()
	hostNames := []string{"bingo.mixmin.net"}
	err := cfg.AddKey("crooks", "/home/crooks/.ssh/openwrt")
	if err != nil {
		panic(err)
	}
	err = cfg.AddKey("crooks", "/home/crooks/.ssh/id_nopass")
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	for _, hostName := range hostNames {
		client, err := cfg.Auth(hostName)
		if err != nil {
			log.Println(err)
			continue
		}
		b, err = cfg.Cmd(client, "cat /etc/passwd")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		b, err = cfg.Cmd(client, "cat /etc/group")
		if err != nil {
			panic(err)
		}
		fmt.Print(b.String())
		client.Close()
	}
}
