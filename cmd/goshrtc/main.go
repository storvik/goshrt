package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"github.com/storvik/goshrt"
	"github.com/storvik/goshrt/http"
	"github.com/storvik/goshrt/version"
	"github.com/urfave/cli/v2"
)

// TODO: Better error descriptions in client

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
	shrtFlags := []cli.Flag{
		&cli.IntFlag{Name: "id"},
		&cli.StringFlag{Name: "domain"},
		&cli.StringFlag{Name: "slug"},
		&cli.StringFlag{Name: "dest"},
		&cli.StringFlag{Name: "expiry"},
	}

	var a *AppConfig
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
			// TODO: Should check config file better, ex validate server address
			a = appcfg
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
						Name:     "add",
						Category: "shrt",
						Usage:    "add shrt, needs domain and destination and optional slug",
						Flags:    shrtFlags,
						Action: func(c *cli.Context) error {
							var t time.Time
							var err error
							if c.String("expiry") != "" {
								t, err = time.Parse(c.String("expiry"), "2006-02-01")
								if err != nil {
									return err
								}
							}
							shrt := &goshrt.Shrt{
								ID:     c.Int("id"),
								Domain: c.String("domain"),
								Slug:   c.String("slug"),
								Dest:   c.String("dest"),
								Expiry: t,
							}
							if shrt.Domain == "" || shrt.Dest == "" {
								return goshrt.ErrInvalid
							}
							client := &http.Client{
								Address: a.Server.Address,
								Key:     a.Client.Key,
							}
							err = client.ShrtAdd(shrt)
							if err != nil {
								return err
							}
							fmt.Printf("Shrt successfully added!\n\n")
							shrt.Printp()
							return nil
						},
					},
					{
						Name:     "get",
						Category: "shrt",
						Usage:    "get shrt details, needs either id or domain and slug ",
						Flags:    shrtFlags,
						Action: func(c *cli.Context) error {
							var t time.Time
							var err error
							if c.String("expiry") != "" {
								t, err = time.Parse(c.String("expiry"), "2006-02-01")
								if err != nil {
									return err
								}
							}
							shrt := &goshrt.Shrt{
								ID:     c.Int("id"),
								Domain: c.String("domain"),
								Slug:   c.String("slug"),
								Dest:   c.String("dest"),
								Expiry: t,
							}
							if shrt.ID == 0 {
								if shrt.Domain == "" || shrt.Slug == "" {
									return goshrt.ErrInvalid
								}
							}
							client := &http.Client{
								Address: a.Server.Address,
								Key:     a.Client.Key,
							}
							err = client.ShrtGet(shrt)
							if err != nil {
								return err
							}
							shrt.Printp()
							return nil
						},
					},
					{
						Name:     "delete",
						Category: "shrt",
						Usage:    "deletes shrt by id",
						Flags:    []cli.Flag{&cli.IntFlag{Name: "id"}},
						Action: func(c *cli.Context) error {
							shrt := &goshrt.Shrt{
								ID: c.Int("id"),
							}
							if shrt.ID == 0 {
								return goshrt.ErrInvalid
							}
							client := &http.Client{
								Address: a.Server.Address,
								Key:     a.Client.Key,
							}
							err := client.ShrtDelete(shrt)
							if err != nil {
								return err
							}
							fmt.Printf("Shrt successfully deleted!\n\n")
							shrt.Printp()
							return nil
						},
					},
					{
						Name:     "list",
						Category: "shrt",
						Usage:    "list all shrts or for given domain if domain is set",
						Flags:    []cli.Flag{&cli.StringFlag{Name: "domain"}},
						Action: func(c *cli.Context) error {
							client := &http.Client{
								Address: a.Server.Address,
								Key:     a.Client.Key,
							}
							shrts, err := client.ShrtGetList(c.String("domain"))
							if err != nil {
								return err
							}
							printList(shrts)
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

func printList(shrts []*goshrt.Shrt) {
	for _, i := range shrts {
		t := i.Expiry.Format("2006.01.02")
		if i.Expiry.IsZero() {
			t = "          "
		}
		fmt.Printf("%3d\t%s\t%-25s\t%-40s\t%s\n", i.ID, t, i.Domain, i.Slug, i.Dest)
	}
}
