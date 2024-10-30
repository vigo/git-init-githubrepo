package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vigo/git-init-githubrepo/internal/version"
)

var templateFilters = template.FuncMap{
	"Upper": strings.ToUpper,
}

const filePerm = 0o0644

type cmd struct {
	writer io.Writer
	app    *cli.App

	cwd             string
	gitPath         string
	gitUserFullName string
	gitUserEmail    string
	gitHubUserName  string
}

func (k *cmd) Run(args []string) error {
	if err := k.app.Run(args); err != nil {
		return fmt.Errorf("could not run app: %w", err)
	}

	return nil
}

func (k *cmd) GenerateTextFromTemplate(fileName string, content any, templateString string) error {
	tmpl, err := template.New(fileName).Funcs(templateFilters).Parse(templateString)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	var wr io.Writer
	if k.writer == nil {
		var file *os.File
		file, err = os.OpenFile(filepath.Clean(fileName), os.O_RDWR|os.O_CREATE, filePerm)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer func() { _ = file.Close() }()
		wr = file
	}

	if err = tmpl.Execute(wr, content); err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	return nil
}

// Option represents option function type for functional options.
type Option func(*cmd)

// WithWriter sets writer.
func WithWriter(wr io.Writer) Option {
	return func(k *cmd) {
		k.writer = wr
	}
}

// New instantiates new gircmd instance.
func New(options ...Option) (*cmd, error) { //nolint:revive
	kommand := &cmd{}

	for _, opt := range options {
		opt(kommand)
	}

	if err := kommand.checkDefaults(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	extrasAvailableLicenses := make([]string, 0, len(availableLicenseTypes))
	for k, v := range availableLicenseTypes {
		extrasAvailableLicenses = append(extrasAvailableLicenses, fmt.Sprintf("  - `%s`: %s", k, v))
	}

	extrasHelpFormatted := fmt.Sprintf(extrasHelp, strings.Join(extrasAvailableLicenses, "\n"))
	cli.AppHelpTemplate = fmt.Sprintf("%s%s\n", cli.AppHelpTemplate, extrasHelpFormatted)
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}

	app := &cli.App{
		EnableBashCompletion: true,
		Version:              version.Version,
		Writer:               kommand.writer,
		Usage:                extrasAppUsage,
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Uğur “vigo” Özyılmazel",
				Email: "ugurozyilmazel@gmail.com",
			},
		},
		// Before:               commandBeforeAction,
		Flags:  kommand.getFlags(),
		Action: kommand.actions(),
	}
	kommand.app = app

	return kommand, nil
}
