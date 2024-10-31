package command

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

//go:embed templates/readme.gotxt
var templateREADME string

//go:embed templates/coc.gotxt
var templateCOC string

//go:embed templates/license/mit.gotxt
var templateLicenseMIT string

//go:embed templates/license/mit-na.gotxt
var templateLicenseMITNA string

//go:embed templates/license/gnu-affero-gpl-30.gotxt
var templateLicenseGNUAfferoGPL30 string

//go:embed templates/bumpversion.txt
var templateBumpVersion string

type (
	licenseType  string
	licenseTypes map[licenseType]string

	licenseMITVariables struct {
		FullName string
		Year     int
	}

	licenseGNUAfferoGPL30Variables struct {
		FullName    string
		ProjectName string
		Year        int
	}

	readmeVariables struct {
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
)

func (lt licenseType) String() string {
	return string(lt)
}

const (
	licenseMIT              = licenseType("mit")
	licenseMITNoAttribution = licenseType("mit-na")
	licenseGNUAfferoGPL30   = licenseType("gnu-agpl30")

	fnReadme      = "README.md"
	fnCOC         = "CODE_OF_CONDUCT.md"
	fnLicense     = "LICENSE"
	fnBumpVersion = ".bumpversion.toml"
)

// sentinel errors.
var (
	ErrProjectNameRequired    = errors.New("project name required")
	ErrRepositoryNameRequired = errors.New("repository name required")
	ErrInvalidLicense         = errors.New("invalid licence option")
	ErrAlreadyFolderExists    = errors.New("folder already exists")
)

var availableLicenseTypes = licenseTypes{
	licenseMIT:              "MIT",
	licenseMITNoAttribution: "MIT No Attribution",
	licenseGNUAfferoGPL30:   "GNU Affero General Public License v3.0",
}

func (k *cmd) actions() func(*cli.Context) error {
	return func(c *cli.Context) error {
		wr := c.App.Writer

		if c.Bool("bash-completion") {
			fmt.Fprintf(wr, "%s\n", extrasBashCompletion)

			return nil
		}

		if c.Bool("list-licenses") {
			fmt.Fprintf(wr, "\n%s: %d\n\n", "available license(s)", len(availableLicenseTypes))
			for k, v := range availableLicenseTypes {
				fmt.Fprintf(wr, "    - `%s`: for `%s` license\n", k, v)
			}
			fmt.Fprintln(wr, "")

			return nil
		}

		argProjectName := c.String("project-name")
		if argProjectName == "" {
			return ErrProjectNameRequired
		}

		argRepositoryName := c.String("repository-name")
		if argRepositoryName == "" {
			return ErrRepositoryNameRequired
		}

		argLicense := c.String("license")
		argNoLicense := c.Bool("no-license")
		if !argNoLicense {
			licenseAsType := licenseType(argLicense)
			if _, ok := availableLicenseTypes[licenseAsType]; !ok {
				lkeys := make([]string, 0, len(availableLicenseTypes))
				for k := range availableLicenseTypes {
					lkeys = append(lkeys, "`"+string(k)+"`")
				}

				return fmt.Errorf(
					"%w `%s`. valid license arguments are: %s",
					ErrInvalidLicense,
					argLicense,
					strings.Join(lkeys, ", "),
				)
			}
		}

		targetFolder := strings.Join(
			[]string{k.cwd, argRepositoryName},
			string(os.PathSeparator),
		)

		if _, err := os.Stat(targetFolder); !os.IsNotExist(err) {
			return fmt.Errorf("%s: %w", targetFolder, ErrAlreadyFolderExists)
		}

		if _, err := k.runGITCommand("init", targetFolder); err != nil {
			return fmt.Errorf("could not initialize git repository at: %s, %w", targetFolder, err)
		}

		argFullName := c.String("full-name")
		argEmail := c.String("email")
		argUserName := c.String("username")
		argDisableFork := c.Bool("disable-fork")
		argDisableCOC := c.Bool("disable-coc")
		argDisableBumpVersion := c.Bool("disable-bumpversion")

		readmeVars := readmeVariables{
			FullName:       argFullName,
			GitHubUsername: argUserName,
			ProjectName:    argProjectName,
			RepositoryName: argRepositoryName,
			License:        argLicense,
			AddLicense:     !argNoLicense,
			AddForkInfo:    !argDisableFork,
			AddCOC:         !argDisableCOC,
			AddBumpVersion: !argDisableBumpVersion,
		}

		readmeFilePath := strings.Join(
			[]string{targetFolder, fnReadme},
			string(os.PathSeparator),
		)

		if err := k.GenerateTextFromTemplate(readmeFilePath, &readmeVars, templateREADME); err != nil {
			return fmt.Errorf("could not generate %s file, %w", fnReadme, err)
		}

		if readmeVars.AddCOC {
			cocFilePath := strings.Join(
				[]string{targetFolder, fnCOC},
				string(os.PathSeparator),
			)

			codeOfConductVars := struct{ Email string }{argEmail}

			if err := k.GenerateTextFromTemplate(cocFilePath, &codeOfConductVars, templateCOC); err != nil {
				return fmt.Errorf("could not generate %s file, %w", fnCOC, err)
			}
		}

		if readmeVars.AddLicense {
			licenseFilePath := strings.Join(
				[]string{targetFolder, fnLicense},
				string(os.PathSeparator),
			)

			now := time.Now()

			switch readmeVars.License {
			case licenseGNUAfferoGPL30.String():
				licenseParams := licenseGNUAfferoGPL30Variables{
					FullName:    argFullName,
					ProjectName: argProjectName,
					Year:        now.Year(),
				}
				if err := k.GenerateTextFromTemplate(licenseFilePath, &licenseParams, templateLicenseGNUAfferoGPL30); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseMIT.String(), licenseMITNoAttribution.String():
				licenseParams := licenseMITVariables{
					FullName: argFullName,
					Year:     now.Year(),
				}

				ltemp := templateLicenseMIT
				if readmeVars.License == licenseMITNoAttribution.String() {
					ltemp = templateLicenseMITNA
				}

				if err := k.GenerateTextFromTemplate(licenseFilePath, &licenseParams, ltemp); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}
			}
		}

		if readmeVars.AddBumpVersion {
			bumpVersionFilePath := strings.Join(
				[]string{targetFolder, fnBumpVersion},
				string(os.PathSeparator),
			)

			if err := k.GenerateTextFromTemplate(bumpVersionFilePath, nil, templateBumpVersion); err != nil {
				return fmt.Errorf("could not generate %s file, %w", fnBumpVersion, err)
			}
		}
		fmt.Fprintf(wr, "your new project is ready at %s\n", targetFolder)

		return nil
	}
}
