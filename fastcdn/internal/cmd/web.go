package cmd

import (
	"github.com/urfave/cli"

	"fastcdn/internal/app"
	"fastcdn/internal/conf"
	"fastcdn/internal/db"
	"fastcdn/internal/log"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "this command start web services",
	Description: `start web services`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {

	conf.InitConf(c.String("config"))

	log.Init()
	log.RewriteStderrFile()

	if conf.Security.InstallLock {
		db.InitDb()
	}

	app.Run()
	return nil
}
