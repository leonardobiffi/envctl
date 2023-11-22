package cli

import (
	"os"

	"github.com/urfave/cli"
)

// Info defines the basic information required for the CLI.
type Info struct {
	Name        string
	Version     string
	Description string
	AuthorName  string
	AuthorEmail string
}

// Initialize and bootstrap the CLI.
func Initialize(info *Info) error {
	var secretName string
	var env string
	var region string
	var profile string
	var envFile string
	var upper bool

	app := cli.NewApp()
	app.Name = info.Name
	app.Version = info.Version
	app.Usage = info.Description
	app.Authors = []cli.Author{
		{
			Name:  info.AuthorName,
			Email: info.AuthorEmail,
		},
	}

	flags := []cli.Flag{
		cli.StringFlag{
			Name:        "secret, s",
			Usage:       "Secret's Name to fetch environment from",
			Destination: &secretName,
		},
		cli.StringFlag{
			Name:        "env, e",
			Usage:       "Environment to use the secret name from",
			Destination: &env,
		},
		cli.StringFlag{
			Name:        "region, r",
			Usage:       "AWS Region",
			Destination: &region,
		},
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "Profile",
			Destination: &profile,
		},
		cli.StringFlag{
			Name:        "envfile, ef",
			Usage:       "Use .env file",
			Destination: &envFile,
		},
		cli.BoolFlag{
			Name:        "upper, u",
			Usage:       "Set Upper Case for all environment variables",
			Destination: &upper,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "setup",
			Usage: "Setup envctl configuration",
			Action: func(ctx *cli.Context) error {
				Setup()
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List environment variables stored in Secrets Manager",
			Flags: flags,
			Action: func(ctx *cli.Context) error {
				List(secretName, env, region, profile, envFile, upper)
				return nil
			},
		},
		{
			Name:  "update",
			Usage: "Update environment variables from env file to Secrets Manager",
			Flags: flags,
			Action: func(ctx *cli.Context) error {
				Update(secretName, env, region, profile, envFile)
				return nil
			},
		},
		{
			Name:      "run",
			Usage:     "Run a command with the injected env variables",
			ArgsUsage: "[command]",
			Flags:     flags,
			Action: func(ctx *cli.Context) error {
				Run(secretName, ctx.Args().Get(0), env, region, profile, envFile)
				return nil
			},
		},
	}

	return app.Run(os.Args)
}
