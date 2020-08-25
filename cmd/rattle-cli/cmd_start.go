package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjeffcaii/rattle/server"
	"github.com/urfave/cli/v2"
)

func newStartCommand() *cli.Command {
	return &cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "id",
				Value: uuid.New().String(),
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "listen port",
				Value: 7878,
			},
			&cli.IntFlag{
				Name:  "cluster-port",
				Usage: "Cluster port",
				Value: 4000,
			},
			&cli.IntFlag{
				Name:  "http-port",
				Usage: "Http port",
				Value: 8080,
			},
			&cli.StringSliceFlag{
				Name:  "seeds",
				Usage: "seeds",
				Value: cli.NewStringSlice(
					"127.0.0.1:4000",
					"127.0.0.1:4001",
					"127.0.0.1:4002",
				),
			},
		},
		Action: func(ctx *cli.Context) error {
			var c server.Config
			c.Port = ctx.Int("port")
			c.ID = ctx.String("id")
			c.Discovery.Seeds = ctx.StringSlice("seeds")
			c.Discovery.Port = ctx.Int("cluster-port")
			r, err := server.NewRattle(&c)
			if err != nil {
				return err
			}
			return r.Serve(context.Background())
		},
	}
}
