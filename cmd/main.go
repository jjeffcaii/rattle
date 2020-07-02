package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jjeffcaii/rattle"
)

var hostname string

var (
	port     int
	seeds    string
	httpPort int
)

func init() {
	hostname, _ = os.Hostname()
	flag.IntVar(&port, "port", 5000, "discovery port")
	flag.StringVar(&seeds, "seeds", "127.0.0.1:4000,127.0.0.1:4001,127.0.0.1:4002", "seeds")
	flag.IntVar(&httpPort, "http-port", 9000, "http port")
	flag.Parse()
}

func main() {
	c := &rattle.Config{
		ID:   uuid.New().String(),
		Port: 10000 + port,
		Discovery: struct {
			Seeds []string
			Port  int
		}{
			Seeds: strings.Split(seeds, ","),
			Port:  port,
		},
	}
	r, err := rattle.NewRattle(c)
	if err != nil {
		panic(err)
	}
	log.Fatalln(r.Serve(context.Background()))
}
