package main

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
	"github.com/upamune/go-esa/esa"
	"github.com/urfave/cli"
)

const (
	// Version
	Version string = "0.1.0"
	// ExitCodeOK ...
	ExitCodeOK int = 0
	// ExitCodeError ..
	ExitCodeError int = 1
	// DefaultConfigFileName...
	DefaultConfigFileName string = "config.toml"
	// DefaultBeforeDayNumber...
	DefaultBeforeDayNumber int = 1
)

// CLI ...
type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

// Config ...
type Config struct {
	AccessToken  string `toml:"access_token"`
	TeamName     string `toml:"team_name"`
	MyScreenName string `toml:"my_screen_name"`
	Path         string `toml:"path"`
}

// Run ...
func (c *CLI) Run(args []string) int {
	var configPath string
	var beforeDayNumber int

	app := cli.NewApp()
	app.Name = "esampo"
	app.Version = Version
	app.Usage = "esampo open"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configration from `FILE`",
			Destination: &configPath,
			Value:       defaultConfigPath(),
		},
		cli.IntFlag{
			Name:        "before-day-number, b",
			Usage:       "before day number",
			Destination: &beforeDayNumber,
			Value:       DefaultBeforeDayNumber,
		},
	}
	app.Action = func(c *cli.Context) error {
		cnf, err := loadConfig(configPath)
		if err != nil {
			return err
		}

		return open(c, cnf, beforeDayNumber)
	}

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func defaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/.config/esampo/%s", home, DefaultConfigFileName)
}

func loadConfig(path string) (*Config, error) {
	c := &Config{}
	if _, err := toml.DecodeFile(path, c); err != nil {
		return nil, err
	}
	return c, nil
}

func open(ctx *cli.Context, cnf *Config, beforeDayNumber int) error {
	client := esa.NewClient(cnf.AccessToken)

	q := url.Values{}
	q.Add("in", time.Now().AddDate(0, 0, beforeDayNumber*-1).Format(cnf.Path))
	res, err := client.Post.GetPosts(cnf.TeamName, q)
	if err != nil {
		return err
	}
	for _, p := range res.Posts {
		if p.CreatedBy.ScreenName == cnf.MyScreenName {
			continue
		}
		err := browser.OpenURL(p.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
