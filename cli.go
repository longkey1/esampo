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
	Version string = "0.4.0"
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
	AccessToken   string `toml:"access_token"`
	TeamName      string `toml:"team_name"`
	MyScreenName  string `toml:"my_screen_name"`
	Path          string `toml:"path"`
}

// Run ...
func (c *CLI) Run(args []string) int {
	var configPath string
	var beforeDayNumber int
	var onlyUser string
	var startDate string
	var endDate string
	var month string

	app := cli.NewApp()
	app.Name = "esampo"
	app.Version = Version
	app.Usage = "esampo open"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configPath,
			Value:       defaultConfigPath(),
		},
		cli.IntFlag{
			Name:        "before-day-number, bd",
			Usage:       "Before day number",
			Destination: &beforeDayNumber,
			Value:       DefaultBeforeDayNumber,
		},
		cli.StringFlag{
			Name:        "user, u",
			Usage:       "Only user",
			Destination: &onlyUser,
			Value:      "",
		},
		cli.StringFlag{
			Name:        "start-date, s",
			Usage:       "Start date",
			Destination: &startDate,
			Value:      "",
		},
		cli.StringFlag{
			Name:        "end-date, e",
			Usage:       "End date",
			Destination: &endDate,
			Value:      "",
		},
		cli.StringFlag{
			Name:        "month, m",
			Usage:       "Target month",
			Destination: &month,
			Value:      "",
		},
	}
	app.Action = func(c *cli.Context) error {
		cnf, err := loadConfig(configPath)
		if err != nil {
			return err
		}

		return run(cnf, beforeDayNumber, onlyUser, startDate, endDate, month)
	}

	err := app.Run(args)
	if err != nil {
		_, _ = fmt.Fprintln(c.errStream, err)
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

func run(cnf *Config, beforeDayNumber int, onlyUserString string, startDateString string, endDateString string, monthString string) error {
	if (startDateString == "" || endDateString == "") && monthString == "" {
		targetDate :=time.Now().AddDate(0, 0, beforeDayNumber*-1)
		return open(cnf, targetDate, onlyUserString)
	}

	var start, end time.Time
	var err error
	if monthString != "" {
		start, err = time.Parse("2006-01-02", monthString + "-01")
		if err != nil {
			return err
		}
		end = start.AddDate(0,1,0).Add(-1)
	} else {
		start, err = time.Parse("2006-01-02", startDateString)
		if err != nil {
			return err
		}
		end, err = time.Parse("2006-01-02", endDateString)
		if err != nil {
			return err
		}
	}

	for target := start; target.After(end) == false; target = target.AddDate(0, 0, 1) {
		err = open(cnf, target, onlyUserString)
		if err != nil {
			return err
		}
	}
	return nil
}

func open(cnf *Config, targetDate time.Time, onlyUser string) error {
	client := esa.NewClient(cnf.AccessToken)

	q := url.Values{}
	q.Add("in", targetDate.Format(cnf.Path))
	if onlyUser != "" {
		q.Add("user", onlyUser)
	} else {
		q.Add("-user", cnf.MyScreenName)
	}
	res, err := client.Post.GetPosts(cnf.TeamName, q)
	if err != nil {
		return err
	}
	for _, p := range res.Posts {
		if onlyUser != "" && p.CreatedBy.ScreenName != onlyUser {
			continue
		}
		err := browser.OpenURL(p.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
