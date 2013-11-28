package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type ServerQuery struct {
	Hostname   string `json:"hostname,omitempty"`
	Address    string `json:"address,omitempty"`
	Mod        string `json:"mod,omitempty"`
	Map        string `json:"map,omitempty"`
	Gametype   string `json:"gametype,omitempty"`
	Online     string `json:"players,omitempty"`
	MaxPlayers string `json:"max_players,omitempty"`
}

func purify_host(a string) string {
	a = strings.Trim(a, "\x00")
	a = strings.Replace(a, "\xE1", "", -1)
	return a
}

func server_query(address string) *ServerQuery {
	conn, err := net.DialTimeout("udp", address, time.Second * 3)
	if err != nil {
		panic(err)
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "\xFF\xFF\xFF\xFFgetinfo\x00")

	var packet [4096]byte
	n, err := conn.Read(packet[:])
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return nil
	}

	received := packet[:n]

	data := strings.Split(string(received), "\\")[1:]

	if len(data) % 2 != 0 {
		fmt.Printf("error: invalid packet returned\n")
	}

	var merged map[string]string
	merged = make(map[string]string)

	for i := 0; i < len(data); i = i + 2 {
		key := data[i]
		value := data[i + 1]
		merged[key] = value
	}

	query := new(ServerQuery)
	query.Hostname = purify_host(merged["hostname"])
	query.Address = address
	query.Mod = merged["fs_game"]
	query.Gametype = merged["gametype"]
	query.Online = merged["clients"]
	query.MaxPlayers = merged["sv_maxclients"]

	return query
}