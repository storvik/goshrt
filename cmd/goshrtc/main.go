package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"github.com/storvik/goshrt/version"
	"github.com/urfave/cli/v2"
)

// AppConfig represents application configuration
type AppConfig struct {
	Client struct {
		Key string `toml:"key"` // Key is the secret key used to authinticate with server
	} `toml:"client"`
	Server struct {
		Address string `toml:"address"` // Address is server address
	} `toml:"Server"`
}

func main() {
	// var a *AppConfig
	app := &cli.App{
		Name:        "goshrtc",
		Usage:       "Goshrt client",
		Description: "Client for self hosted URL shortener, goshrt",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Before: func(c *cli.Context) error {
			var err error
			appcfg := new(AppConfig)
			cfg := c.String("config")
			if cfg == "" {
				cfg, err = xdg.SearchConfigFile("goshrt/client.toml")
				if err != nil {
					return err
				}
			}
			if buf, err := ioutil.ReadFile(cfg); err != nil {
				return err
			} else if err := toml.Unmarshal(buf, appcfg); err != nil {
				return err
			}
			// a = appcfg
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version",
				Action: func(c *cli.Context) error {
					version.Print()
					return nil
				},
			},
			{
				Name:    "shrt",
				Aliases: []string{"s"},
				Usage:   "handle shrts through the server api",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "[client name] generate new valid jwt token to be used by clients",
						Action: func(c *cli.Context) error {
							fmt.Println("Not implemented yet")
							return nil
						},
					},
				},
			},
		},
	}

	// Run application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
