package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vigo/git-init-githubrepo/internal/version"
)

func templateFilters() template.FuncMap {
	return template.FuncMap{
		"Upper": strings.ToUpper,
	}
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
	tmpl, err := template.New(fileName).Funcs(templateFilters()).Parse(templateString)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	var wr io.Writer

	wr = k.writer

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

	if kommand.gitHubUserName == "" {
		kommand.gitHubUserName = "your-github-username"
	}
	if kommand.gitUserFullName == "" {
		kommand.gitUserFullName = "Your Full Name"
	}
	if kommand.gitUserEmail == "" {
		kommand.gitUserEmail = "your@email"
	}

	licenseTypeKeys := make([]string, 0, len(availableLicenseTypes()))
	for k := range availableLicenseTypes() {
		licenseTypeKeys = append(licenseTypeKeys, k.String())
	}
	sort.Strings(licenseTypeKeys)

	extrasAvailableLicenses := make([]string, 0, len(licenseTypeKeys))
	for _, k := range licenseTypeKeys {
		extrasAvailableLicenses = append(
			extrasAvailableLicenses,
			fmt.Sprintf("  - `%s`: %s", k, availableLicenseTypes()[licenseType(k)]),
		)
	}

	projectStyleKeys := make([]string, 0, len(availableProjectStyles()))
	for k := range availableProjectStyles() {
		projectStyleKeys = append(projectStyleKeys, k.String())
	}
	sort.Strings(projectStyleKeys)

	extrasProjectStyles := make([]string, 0, len(projectStyleKeys))
	for _, k := range projectStyleKeys {
		extrasProjectStyles = append(
			extrasProjectStyles,
			fmt.Sprintf("  - `%s`", k),
		)
	}

	extrasHelpFormatted := fmt.Sprintf(
		extrasHelp(),
		len(licenseTypeKeys),
		strings.Join(extrasAvailableLicenses, "\n"),
		len(projectStyleKeys),
		strings.Join(extrasProjectStyles, "\n"),
	)
	cli.AppHelpTemplate = fmt.Sprintf("%s%s\n", cli.AppHelpTemplate, extrasHelpFormatted)
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}

	app := &cli.App{
		EnableBashCompletion: true,
		Version:              version.Version,
		Writer:               kommand.writer,
		Usage:                extrasAppUsage(),
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Uğur “vigo” Özyılmazel",
				Email: "ugurozyilmazel@gmail.com",
			},
		},
		Flags:  kommand.getFlags(),
		Action: kommand.actions(),
	}
	kommand.app = app

	return kommand, nil
}
