package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
)

const (
	defaultPort = 3000
)

func main() {
	var listenOnPort int
	var upstream string

	flag.IntVar(&listenOnPort, "port", defaultPort, "Port that HTTP server should listen on")
	flag.StringVar(&upstream, "upstream", "", "The URL of upstream Apisonator server")
	flag.Parse()

	if err := validateFlags(upstream); err != nil {
		log.Fatal(err)
	}
}

func validateFlags(upstream string) error {
	if upstream == "" {
		return errors.New("invalid input. upstream must be provided")
	}

	if _, err := url.ParseRequestURI(upstream); err != nil {
		return fmt.Errorf("invalid input, failed to parse provided upstream %s - %v", upstream, err)
	}
	return nil
}
