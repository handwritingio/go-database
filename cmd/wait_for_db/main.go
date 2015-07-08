package main

// This script waits for the database container to begin accepting connections,
// because we don't want the tests to fail just because the database hadn't
// fully started yet.
//
// https://github.com/docker/compose/issues/374 shows a lot of people dealing
// with the same problem. The solution used here is based on the one shown in
// https://github.com/aanand/docker-wait. Unfortunately, the author explains that
// it doesn't work with docker-compose:
// https://github.com/docker/compose/issues/374#issuecomment-69212755

import (
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

const timeout = 5 * time.Second

func main() {
	databaseUrl, err := url.Parse(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	address := databaseUrl.Host

	log.Printf("Waiting for TCP connection to %s...\n", address)

	_, err = net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("ok")
}
