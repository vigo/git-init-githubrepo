package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/urfave/cli/v2"
	"github.com/vigo/git-init-githubrepo/version"
)

//go:embed templates/readme.gotxt
var templateReadme string

//go:embed templates/coc.gotxt
var templateCOC string

//go:embed templates/license/mit.gotxt
var templateLicenseMIT string

//go:embed templates/bumpversion.cfg
var templateBumpVersion string

// ReadmeParams fields hold required data for README.md
type ReadmeParams struct {
	FullName       string
	GitHubUsername string
	ProjectName    string
	RepositoryName string
	License        string
	AddLicense     bool
	AddForkInfo    bool
	AddCOC         bool
	AddBumpVersion bool
}

// LicenseMITParams fields hold required data LICENSE
type LicenseMITParams struct {
	Year     int
	FullName string
}

var templateFilters = template.FuncMap{
	"Upper": strings.ToUpper,
}

func getFromGitConfig(configName string) string {
	buff := &bytes.Buffer{}

	cmd := exec.Command("git", "config", configName)
	cmd.Stdout = buff

	if err := cmd.Run(); err != nil {
		return ""
	}

	return string(bytes.TrimSpace(buff.Bytes()))
}

func commandExists(exe string) error {
	return exec.Command("command", "-v", exe).Run()
}

func inGITRepo() error {
	return exec.Command("git", "rev-parse", "--git-dir").Run()
}

func gitInit(path string) error {
	return exec.Command("git", "init", path).Run()
}

type licenseTypes []string

func (l licenseTypes) String() string {
	o := ""
	for i, typeName := range l {
		o = o + typeName
		if i+1 < len(l) {
			o = o + ","
		}
	}
	return o
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "version information",
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}

	cli.AppHelpTemplate = fmt.Sprintf("%s%s\n", cli.AppHelpTemplate, helpExtras)

	availableLicenses := licenseTypes{"mit"}

	app := &cli.App{
		EnableBashCompletion: true,
		Version:              version.Version,
		Usage:                "create git repository with built-in README, LICENSE and more...",
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Uğur \"vigo\" Özyılmazel",
				Email: "ugurozyilmazel@gmail.com",
			},
		},
	}

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:  "bash-completion",
			Usage: "generate bash-completion code",
		},
		&cli.StringFlag{
			Name:    "full-name",
			Aliases: []string{"f"},
			Usage:   "your `FULLNAME`",
			Value:   getFromGitConfig("user.name"),
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "your GitHub `USERNAME`",
			Value:   getFromGitConfig("github.user"),
		},
		&cli.StringFlag{
			Name:    "email",
			Aliases: []string{"e"},
			Usage:   "your contact `EMAIL`",
			Value:   getFromGitConfig("user.email"),
		},
		&cli.StringFlag{
			Name:    "project-name",
			Aliases: []string{"p"},
			Usage:   "`NAME` of your project",
		},
		&cli.StringFlag{
			Name:    "repository-name",
			Aliases: []string{"r"},
			Usage:   "`NAME` of your GitHub repository",
		},
		&cli.StringFlag{
			Name:    "license",
			Aliases: []string{"l"},
			Usage:   fmt.Sprintf("add `LICENSE`. available license(s): %s", availableLicenses),
			Value:   "mit",
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

	app.Flags = flags

	app.Action = func(c *cli.Context) error {
		if c.Bool("bash-completion") {
			fmt.Fprintf(c.App.Writer, "%s\n", bashCompletion)
			return nil
		}

		if c.String("project-name") == "" {
			return cli.Exit("project name required", 1)
		}
		if c.String("repository-name") == "" {
			return cli.Exit("repository name required", 1)
		}
		if c.String("full-name") == "" {
			return cli.Exit("full name required", 1)
		}
		if c.String("username") == "" {
			return cli.Exit("username required", 1)
		}
		if !c.Bool("no-license") {
			licenseValid := false
			for _, licenseType := range availableLicenses {
				if c.String("license") == licenseType {
					licenseValid = true
					break
				}
			}
			if !licenseValid {
				return cli.Exit(fmt.Sprintf("invalid license type: %s", c.String("license")), 1)
			}
		}
		if !c.Bool("disable-coc") {
			if c.String("email") == "" {
				return cli.Exit("you need to provide email due to code of conduct choise!", 1)
			}
		}
		if err := commandExists("git"); err != nil {
			fmt.Fprintf(os.Stderr, "you need to instal %q to continue...", "git")
			os.Exit(1)
		}
		if inGITRepo() == nil {
			fmt.Fprintln(os.Stderr, "you are now in a git repository!")
			os.Exit(1)
		}

		// check folder exists?
		cwd, err := os.Getwd()
		if err != nil {
			return cli.Exit(err, 1)
		}
		targetFolder := cwd + "/" + c.String("repository-name")
		if _, err := os.Stat(targetFolder); !os.IsNotExist(err) {
			return cli.Exit(fmt.Sprintf("folder %q already exists!", targetFolder), 1)
		}

		// git init
		if err := gitInit(c.String("repository-name")); err != nil {
			return cli.Exit(err, 1)
		}

		// create files under folder
		readmeParams := &ReadmeParams{
			FullName:       c.String("full-name"),
			GitHubUsername: c.String("username"),
			ProjectName:    c.String("project-name"),
			RepositoryName: c.String("repository-name"),
			License:        c.String("license"),
			AddLicense:     !c.Bool("no-license"),
			AddForkInfo:    !c.Bool("disable-fork"),
			AddCOC:         !c.Bool("disable-coc"),
			AddBumpVersion: !c.Bool("disable-bumpversion"),
		}
		// create README
		if err := createFile(readmeParams, targetFolder+"/README.md", templateReadme); err != nil {
			return cli.Exit(err, 1)
		}

		// create CODE_OF_CONDUCT
		if readmeParams.AddCOC {
			if err := createFile(struct{ Email string }{c.String("email")}, targetFolder+"/CODE_OF_CONDUCT.md", templateCOC); err != nil {
				return cli.Exit(err, 1)
			}
		}

		// create LICENSE
		if readmeParams.AddLicense {
			switch readmeParams.License {
			case "mit":
				now := time.Now()
				licenseParams := &LicenseMITParams{
					Year:     now.Year(),
					FullName: readmeParams.FullName,
				}
				if err := createFile(licenseParams, targetFolder+"/LICENSE.md", templateLicenseMIT); err != nil {
					return cli.Exit(err, 1)
				}
			}
		}

		// create .bumpversion.cfg
		if readmeParams.AddBumpVersion {
			if err := createFile(struct{}{}, targetFolder+"/.bumpversion.cfg", templateBumpVersion); err != nil {
				return cli.Exit(err, 1)
			}
		}

		fmt.Fprintf(c.App.Writer, "your new project is ready at %s\n", targetFolder)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}

func createFile(s interface{}, fileName string, ts string) error {
	tmpl, err := template.New(fileName).Funcs(templateFilters).Parse(ts)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, s)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
