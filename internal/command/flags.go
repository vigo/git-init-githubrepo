package command

import "github.com/urfave/cli/v2"

func (c *cmd) getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "bash-completion",
			Usage: "generate bash-completion code",
		},
		&cli.StringFlag{
			Name:    "full-name",
			Aliases: []string{"f"},
			Usage:   "your `FULLNAME`",
			Value:   c.gitUserFullName,
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "your GitHub `USERNAME`",
			Value:   c.gitHubUserName,
		},
		&cli.StringFlag{
			Name:    "email",
			Aliases: []string{"e"},
			Usage:   "your contact `EMAIL`",
			Value:   c.gitUserEmail,
		},
		&cli.StringFlag{
			Name:    "project-name",
			Aliases: []string{"p"},
			Usage:   "`NAME` of your project",
		},
		&cli.StringFlag{
			Name:    "project-type",
			Aliases: []string{"pt"},
			Usage:   "type of your project",
		},
		&cli.StringFlag{
			Name:    "repository-name",
			Aliases: []string{"r"},
			Usage:   "`NAME` of your GitHub repository",
		},

		&cli.StringFlag{
			Name:    "license",
			Aliases: []string{"l"},
			Usage:   "add `LICENSE`",
			Value:   licenseMIT.String(),
		},

		&cli.BoolFlag{
			Name:    "list-licenses",
			Aliases: []string{"ll"},
			Usage:   "list licenses",
		},

		&cli.BoolFlag{
			Name:  "no-license",
			Usage: "do not add LICENSE file",
		},

		&cli.BoolFlag{
			Name:  "disable-fork",
			Usage: "do not add fork information to README",
		},

		&cli.BoolFlag{
			Name:  "disable-bumpversion",
			Usage: "do not create .bumpversion.cfg and badge to README",
		},

		&cli.BoolFlag{
			Name:  "disable-coc",
			Usage: "do not add CODE_OF_CONDUCT",
		},
	}
}
