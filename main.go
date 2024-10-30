package main

// import (
// 	"bytes"
// 	_ "embed"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// 	"text/template"
// 	"time"
//
// 	"github.com/urfave/cli/v2"
// 	"github.com/vigo/git-init-githubrepo/version"
// )
//
// const filePerm = 0o0600
//
// //go:embed templates/readme.gotxt
// var templateReadme string
//
// //go:embed templates/coc.gotxt
// var templateCOC string
//
// //go:embed templates/license/mit.gotxt
// var templateLicenseMIT string
//
// //go:embed templates/license/mit-na.gotxt
// var templateLicenseMITNA string
//
// //go:embed templates/bumpversion.cfg
// var templateBumpVersion string
//
// // ReadmePlaceholder holds required params for README.md.
// type ReadmePlaceholder struct {
// 	FullName       string
// 	GitHubUsername string
// 	ProjectName    string
// 	RepositoryName string
// 	License        string
// 	AddLicense     bool
// 	AddForkInfo    bool
// 	AddCOC         bool
// 	AddBumpVersion bool
// }
//
// // LicenseMITParams fields hold required data LICENSE.
// type LicenseMITParams struct {
// 	FullName string
// 	Year     int
// }
//
// var templateFilters = template.FuncMap{
// 	"Upper": strings.ToUpper,
// }
//
// var (
// 	errInvalidLicense      = errors.New("invalid licence option")
// 	errEmailRequired       = errors.New("email required")
// 	errGitRequired         = errors.New("you need to install git")
// 	errInGitRepo           = errors.New("you are now in a git repo")
// 	errFolderExists        = errors.New("folder already exists")
// 	errProjectNameRequired = errors.New("project name required")
// 	errRepoNameRequired    = errors.New("repository name required")
// )
//
// type licenseTypes map[string]string
//
// func (l licenseTypes) String() string {
// 	ks := make([]string, 0, len(l))
//
// 	for key := range l {
// 		ks = append(ks, key)
// 	}
//
// 	return strings.Join(ks, ",")
// }
//
// var availableLicenses = licenseTypes{
// 	"mit":    "MIT",
// 	"mit-na": "MIT No Attribution",
// }
//
// func getFromGitConfig(configName string) string {
// 	if flag.Lookup("test.v") != nil {
// 		return ""
// 	}
//
// 	buff := &bytes.Buffer{}
//
// 	cmd := exec.Command("git", "config", configName)
// 	cmd.Stdout = buff
//
// 	if err := cmd.Run(); err != nil {
// 		return ""
// 	}
//
// 	return string(bytes.TrimSpace(buff.Bytes()))
// }
//
// func commandExists(exe string) error {
// 	_, err := exec.LookPath(exe)
//
// 	return fmt.Errorf("commandExists err: %w", err)
// }
//
// func inGITRepo() error {
// 	if flag.Lookup("test.v") != nil {
// 		_ = os.Chdir(os.TempDir())
// 	}
//
// 	return fmt.Errorf("git rev-parse err: %w", exec.Command("git", "rev-parse", "--git-dir").Run())
// }
//
// func gitInit(path string) error {
// 	return fmt.Errorf("git init err: %w", exec.Command("git", "init", path).Run())
// }
//
// func createFile(content any, fileName string, ts string) error {
// 	tmpl, err := template.New(fileName).Funcs(templateFilters).Parse(ts)
// 	if err != nil {
// 		return fmt.Errorf("%w", err)
// 	}
//
// 	file, err := os.OpenFile(filepath.Clean(fileName), os.O_RDWR|os.O_CREATE, filePerm)
// 	if err != nil {
// 		return fmt.Errorf("%w", err)
// 	}
//
// 	if err = tmpl.Execute(file, content); err != nil {
// 		return fmt.Errorf("%w", err)
// 	}
//
// 	if err = file.Close(); err != nil {
// 		return fmt.Errorf("%w", err)
// 	}
//
// 	return nil
// }
//
// func getCWD() (string, error) {
// 	if flag.Lookup("test.v") != nil {
// 		tmpDir := strings.TrimRight(os.TempDir(), string(os.PathSeparator))
//
// 		return tmpDir, nil
// 	}
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return "", fmt.Errorf("%w", err)
// 	}
//
// 	return cwd, nil
// }
//
// func run(args []string, wr io.Writer) error {
// 	cli.VersionFlag = &cli.BoolFlag{
// 		Name:    "version",
// 		Aliases: []string{"v"},
// 		Usage:   "version information",
// 	}
// 	cli.VersionPrinter = func(c *cli.Context) {
// 		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
// 	}
// 	cli.AppHelpTemplate = fmt.Sprintf("%s%s\n", cli.AppHelpTemplate, helpExtras)
//
// 	app := &cli.App{
// 		EnableBashCompletion: true,
// 		Version:              version.Version,
// 		Usage:                "create git repository with built-in README, LICENSE and more...",
// 		Compiled:             time.Now(),
// 		Authors: []*cli.Author{
// 			{
// 				Name:  "Uğur \"vigo\" Özyılmazel",
// 				Email: "ugurozyilmazel@gmail.com",
// 			},
// 		},
// 		Writer: wr,
// 	}
//
// 	flags := []cli.Flag{
// 		&cli.BoolFlag{
// 			Name:  "bash-completion",
// 			Usage: "generate bash-completion code",
// 		},
// 		&cli.StringFlag{
// 			Name:    "full-name",
// 			Aliases: []string{"f"},
// 			Usage:   "your `FULLNAME`",
// 			Value:   getFromGitConfig("user.name"),
// 		},
// 		&cli.StringFlag{
// 			Name:    "username",
// 			Aliases: []string{"u"},
// 			Usage:   "your GitHub `USERNAME`",
// 			Value:   getFromGitConfig("github.user"),
// 		},
// 		&cli.StringFlag{
// 			Name:    "email",
// 			Aliases: []string{"e"},
// 			Usage:   "your contact `EMAIL`",
// 			Value:   getFromGitConfig("user.email"),
// 		},
// 		&cli.StringFlag{
// 			Name:    "project-name",
// 			Aliases: []string{"p"},
// 			Usage:   "`NAME` of your project",
// 		},
// 		&cli.StringFlag{
// 			Name:    "repository-name",
// 			Aliases: []string{"r"},
// 			Usage:   "`NAME` of your GitHub repository",
// 		},
// 		&cli.StringFlag{
// 			Name:    "license",
// 			Aliases: []string{"l"},
// 			Usage:   "add `LICENSE`",
// 			Value:   "mit",
// 		},
// 		&cli.BoolFlag{
// 			Name:    "list-licenses",
// 			Aliases: []string{"ll"},
// 			Usage:   "list licenses",
// 		},
// 		&cli.BoolFlag{
// 			Name:  "no-license",
// 			Usage: "do not add LICENSE file",
// 		},
// 		&cli.BoolFlag{
// 			Name:  "disable-fork",
// 			Usage: "do not add fork information to README",
// 		},
// 		&cli.BoolFlag{
// 			Name:  "disable-bumpversion",
// 			Usage: "do not create .bumpversion.cfg and badge to README",
// 		},
// 		&cli.BoolFlag{
// 			Name:  "disable-coc",
// 			Usage: "do not add CODE_OF_CONDUCT",
// 		},
// 	}
//
// 	app.Flags = flags
//
// 	app.Action = func(c *cli.Context) error {
// 		if c.Bool("bash-completion") {
// 			fmt.Fprintf(c.App.Writer, "%s\n", bashCompletion)
//
// 			return nil
// 		}
//
// 		if c.Bool("list-licenses") {
// 			t := []string{"available licence types are:"}
// 			for k, v := range availableLicenses {
// 				t = append(t, "  - `"+k+"`: "+v)
// 			}
// 			fmt.Fprintln(c.App.Writer, strings.Join(t, "\n"))
//
// 			return nil
// 		}
//
// 		if err := commandExists("git"); err != nil {
// 			return errGitRequired
// 		}
//
// 		if inGITRepo() == nil {
// 			return errInGitRepo
// 		}
//
// 		if c.String("project-name") == "" {
// 			return errProjectNameRequired
// 		}
//
// 		if c.String("repository-name") == "" {
// 			return errRepoNameRequired
// 		}
//
// 		if !c.Bool("no-license") {
// 			_, ok := availableLicenses[c.String("license")]
// 			if !ok {
// 				return errInvalidLicense
// 			}
// 		}
//
// 		if !c.Bool("disable-coc") && c.String("email") == "" {
// 			return errEmailRequired
// 		}
//
// 		cwd, err := getCWD()
// 		if err != nil {
// 			return fmt.Errorf("%w", err)
// 		}
//
// 		targetFolder := strings.Join([]string{
// 			cwd,
// 			c.String("repository-name"),
// 		}, string(os.PathSeparator))
//
// 		if _, err = os.Stat(targetFolder); !os.IsNotExist(err) {
// 			return errFolderExists
// 		}
//
// 		if err = gitInit(c.String("repository-name")); err != nil {
// 			return fmt.Errorf("%w", err)
// 		}
//
// 		placeholder := ReadmePlaceholder{
// 			FullName:       c.String("full-name"),
// 			GitHubUsername: c.String("username"),
// 			ProjectName:    c.String("project-name"),
// 			RepositoryName: c.String("repository-name"),
// 			License:        c.String("license"),
// 			AddLicense:     !c.Bool("no-license"),
// 			AddForkInfo:    !c.Bool("disable-fork"),
// 			AddCOC:         !c.Bool("disable-coc"),
// 			AddBumpVersion: !c.Bool("disable-bumpversion"),
// 		}
// 		readmeFilePath := strings.Join([]string{
// 			targetFolder,
// 			"README.md",
// 		}, string(os.PathSeparator))
//
// 		if err = createFile(&placeholder, readmeFilePath, templateReadme); err != nil {
// 			return fmt.Errorf("%w", err)
// 		}
//
// 		if placeholder.AddCOC {
// 			cocFilePath := strings.Join([]string{
// 				targetFolder,
// 				"CODE_OF_CONDUCT.md",
// 			}, string(os.PathSeparator))
//
// 			if err = createFile(struct{ Email string }{c.String("email")}, cocFilePath, templateCOC); err != nil {
// 				return fmt.Errorf("%w", err)
// 			}
// 		}
//
// 		if placeholder.AddLicense {
// 			licenceFilePath := strings.Join([]string{
// 				targetFolder,
// 				"LICENSE",
// 			}, string(os.PathSeparator))
//
// 			switch placeholder.License {
// 			case "mit-na":
// 				now := time.Now()
// 				licenseParams := LicenseMITParams{
// 					Year:     now.Year(),
// 					FullName: placeholder.FullName,
// 				}
// 				if err = createFile(&licenseParams, licenceFilePath, templateLicenseMITNA); err != nil {
// 					return fmt.Errorf("%w", err)
// 				}
// 			case "mit":
// 				now := time.Now()
// 				licenseParams := LicenseMITParams{
// 					Year:     now.Year(),
// 					FullName: placeholder.FullName,
// 				}
// 				if err = createFile(&licenseParams, licenceFilePath, templateLicenseMIT); err != nil {
// 					return fmt.Errorf("%w", err)
// 				}
// 			}
// 		}
// 		if placeholder.AddBumpVersion {
// 			bumpconfigFilePath := strings.Join([]string{
// 				targetFolder,
// 				".bumpversion.cfg",
// 			}, string(os.PathSeparator))
//
// 			if err = createFile(struct{}{}, bumpconfigFilePath, templateBumpVersion); err != nil {
// 				return fmt.Errorf("%w", err)
// 			}
// 		}
//
// 		fmt.Fprintf(c.App.Writer, "your new project is ready at %s\n", targetFolder)
//
// 		return nil
// 	}
//
// 	return fmt.Errorf("%w", app.Run(args))
// }
//
// func main() {
// 	if err := run(os.Args, nil); err != nil {
// 		fmt.Fprintln(os.Stderr, err.Error())
// 		os.Exit(1)
// 	}
// }
