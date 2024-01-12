package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
	"github.com/storvik/goshrt/postgres"
	"github.com/storvik/goshrt/token"
	"github.com/storvik/goshrt/version"
	"github.com/urfave/cli/v2"
)

// AppConfig represents application configuration.
type AppConfig struct {
	Server struct {
		Key        string `toml:"key"`         // Key is the serever key used to create tokens and validate API calls
		Port       string `toml:"port"`        // Port is the server port where the API should be served
		SlugLength int    `toml:"slug_length"` // SlugLength is the default slug length when generating slugs
	} `toml:"server"`
	Database struct {
		DB       string `toml:"db"`       // database name
		User     string `toml:"user"`     // database username
		Password string `toml:"password"` // database password
		Address  string `toml:"address"`  // database address
		Schema   string `toml:"schema"`   // database schema, public if not used
		Port     string `toml:"port"`     // database port
	} `toml:"database"`
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	cfg *AppConfig
}

func main() {
	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	a := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	app := &cli.App{
		Name:        "goshrt",
		Usage:       "Self hosted URL shortener",
		Description: "Self hosted URL shortener written in Go",
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
				cfg, err = xdg.SearchConfigFile("goshrt/server.toml")
				if err != nil {
					return err
				}
			}

			if buf, err := os.ReadFile(cfg); err != nil {
				return err
			} else if err := toml.Unmarshal(buf, appcfg); err != nil {
				return err
			}
			a.cfg = appcfg
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
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "options for handling jwt token",
				Subcommands: []*cli.Command{
					{
						Name:  "generate",
						Usage: "[client name] generate new valid jwt token to be used by clients",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("invalid number of arguments")
							}
							auth := token.NewAuth(a.cfg.Server.Key)
							toknStr, err := auth.Create(c.Args().First())
							if err != nil {
								return err
							}
							fmt.Println(toknStr)
							return nil
						},
					},
					{
						Name:  "validate",
						Usage: "validate jwt token",
						Action: func(c *cli.Context) error {
							auth := token.NewAuth(a.cfg.Server.Key)

							switch c.NArg() {
							case 0:
								fmt.Printf("JWT Token: ")
								var toknStr string
								fmt.Scanln(&toknStr)
								valid, err := auth.Validate(toknStr)
								if err != nil {
									return err
								}
								if valid {
									a.infoLog.Println("Token is valid")
								} else {
									a.infoLog.Println("Token is NOT valid")
								}

							case 1:
								valid, err := auth.Validate(c.Args().First())
								if err != nil {
									return err
								}
								if valid {
									a.infoLog.Println("Token is valid")
								} else {
									a.infoLog.Println("Token is NOT valid")
								}

							default:
								return errors.New("invalid number of arguments")
							}
							return nil
						},
					},
				},
			},

			{
				Name:  "database",
				Usage: "database operations",
				Subcommands: []*cli.Command{
					{
						Name:  "migrate",
						Usage: "migrate database, no other instance of goshrt must be running",
						Action: func(c *cli.Context) error {
							db := postgres.NewClient(a.cfg.Database.DB, a.cfg.Database.User, a.cfg.Database.Password, a.cfg.Database.Address, a.cfg.Database.Schema)
							if err := db.Open(); err != nil {
								return err
							}
							if err := db.Migrate(); err != nil {
								return err
							}
							if err := db.Close(); err != nil {
								return err
							}
							return nil
						},
					},
				},
			},
		},
		Action: func(c *cli.Context) error {
			return a.Serve()
		},
	}

	// Run application
	err := app.Run(os.Args)
	if err != nil {
		a.errorLog.Fatal(err)
	}
}
