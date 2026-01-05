package command

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"sort"
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

//go:embed templates/license/gnu-lesser-gpl-30.gotxt
var templateLicenseGNULesserGPL30 string

//go:embed templates/license/gnu-gpl-30.gotxt
var templateLicenseGNUGPL30 string

//go:embed templates/license/mozilla-public-20.gotxt
var templateLicenseMOZP20 string

//go:embed templates/license/apache-20.gotxt
var templateLicenseAPACHE20 string

//go:embed templates/license/bsl-10.gotxt
var templateLicenseBSL10 string

//go:embed templates/license/the-unlicense.gotxt
var templateLicenseTHEUNL string

//go:embed templates/bumpversion.txt
var templateBumpVersion string

type (
	licenseType  string
	licenseTypes map[licenseType]string

	licenseMITVariables struct {
		FullName string
		Year     int
	}

	licenseGNUGPL30Variables struct {
		FullName    string
		ProjectName string
		Year        int
	}

	licenseAPACHEVariables struct {
		FullName string
		Year     int
	}

	readmeVariables struct {
		FullName               string
		GitHubUsername         string
		ProjectName            string
		RepositoryName         string
		License                string
		LicenseDescription     string
		AddLicense             bool
		AddForkInfo            bool
		AddCOC                 bool
		AddBumpVersion         bool
		AddCodeowners          bool
		AddFunding             bool
		AddPullRequestTemplate bool
		AddIssueTemplate       bool
		AddSecurity            bool
	}
	projectStyle  string
	projectStyles map[projectStyle]string
)

func (lt licenseType) String() string {
	return string(lt)
}

func (ps projectStyle) String() string {
	return string(ps)
}

const (
	licenseMIT              = licenseType("mit")
	licenseMITNoAttribution = licenseType("mit-na")
	licenseGNUAfferoGPL30   = licenseType("gnu-agpl30")
	licenseGNUGPL30         = licenseType("gnu-gpl30")
	licenseGNULesserGPL30   = licenseType("gnu-lgpl30")
	licenseMOZP20           = licenseType("moz-p20")
	licenseAPACHE20         = licenseType("apache-20")
	licenseBSL10            = licenseType("bsl-10")
	licenseTHEUNL           = licenseType("unli")

	projectStyleGo = projectStyle("go")

	fnReadme      = "README.md"
	fnCOC         = "CODE_OF_CONDUCT.md"
	fnLicense     = "LICENSE"
	fnBumpVersion = ".bumpversion.toml"

	// fnIssueTemplateFeatureRequest = "feature_request.md".
)

// sentinel errors.
var (
	ErrProjectNameRequired    = errors.New("project name required")
	ErrRepositoryNameRequired = errors.New("repository name required")
	ErrInvalidLicense         = errors.New("invalid licence option")
	ErrAlreadyFolderExists    = errors.New("folder already exists")
)

func availableLicenseTypes() licenseTypes {
	return licenseTypes{
		licenseMIT:              "MIT",
		licenseMITNoAttribution: "MIT No Attribution",
		licenseGNUAfferoGPL30:   "GNU Affero General Public License v3.0",
		licenseGNUGPL30:         "GNU General Public License v3.0",
		licenseGNULesserGPL30:   "GNU Lesser General Public License v3.0",
		licenseMOZP20:           "Mozilla Public License 2.0",
		licenseAPACHE20:         "Apache License 2.0",
		licenseBSL10:            "Boost Software License 1.0",
		licenseTHEUNL:           "The Unlicense",
	}
}

func availableProjectStyles() projectStyles {
	return projectStyles{
		projectStyleGo: `creates .github/workflows/, linter and tester actions, 
            golangci.yml, .pre-commit-config.yaml, dependabot.yml, .gitignore
            .pre-commit-config.yaml, .codecov.yml`,
	}
}

func (k *cmd) actions() func(*cli.Context) error {
	return func(c *cli.Context) error {
		wr := c.App.Writer

		if c.Bool("bash-completion") {
			fmt.Fprintf(wr, "%s\n", extrasBashCompletion())

			return nil
		}

		if c.Bool("list-licenses") {
			fmt.Fprintf(wr, "\n%s: %d\n\n", "available license(s)", len(availableLicenseTypes()))

			keys := make([]string, 0, len(availableLicenseTypes()))
			for k := range availableLicenseTypes() {
				keys = append(keys, k.String())
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(wr, "    - `%s`: for `%s` license\n", k, availableLicenseTypes()[licenseType(k)])
			}
			fmt.Fprintln(wr, "")

			return nil
		}

		if c.Bool("list-project-styles") {
			fmt.Fprintf(wr, "\n%s: %d\n\n", "available project style(s)", len(availableProjectStyles()))

			keys := make([]string, 0, len(availableProjectStyles()))
			for k := range availableProjectStyles() {
				keys = append(keys, k.String())
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(wr, "    - `%s`: %s\n", k, availableProjectStyles()[projectStyle(k)])
			}
			fmt.Fprintln(wr, "")

			return nil
		}

		if existingRepoPath, _ := k.runGITCommand("rev-parse", "--git-dir"); existingRepoPath != "" {
			return ErrAlreadyInAGitRepo
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
		argNoLicense := c.Bool("disable-license")
		if !argNoLicense {
			licenseAsType := licenseType(argLicense)
			if _, ok := availableLicenseTypes()[licenseAsType]; !ok {
				lkeys := make([]string, 0, len(availableLicenseTypes()))
				for k := range availableLicenseTypes() {
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
		argLicenseDescription := availableLicenseTypes()[licenseType(argLicense)]

		argDisableCodeowners := c.Bool("disable-codeowners")
		argDisableFunding := c.Bool("disable-funding")
		argDisablePullRequestTemplate := c.Bool("disable-pull-request-template")
		argDisableSecurity := c.Bool("disable-security")
		argDisableIssueTemplate := c.Bool("disable-issue-template")

		readmeVars := readmeVariables{
			FullName:           argFullName,
			GitHubUsername:     argUserName,
			ProjectName:        argProjectName,
			RepositoryName:     argRepositoryName,
			License:            argLicense,
			LicenseDescription: argLicenseDescription,
			AddLicense:         !argNoLicense,
			AddForkInfo:        !argDisableFork,
			AddCOC:             !argDisableCOC,
			AddBumpVersion:     !argDisableBumpVersion,

			AddCodeowners:          !argDisableCodeowners,
			AddFunding:             !argDisableFunding,
			AddPullRequestTemplate: !argDisablePullRequestTemplate,
			AddSecurity:            !argDisableSecurity,
			AddIssueTemplate:       !argDisableIssueTemplate,
		}

		var createGitHubFolder bool
		if readmeVars.AddCodeowners || readmeVars.AddFunding || readmeVars.AddPullRequestTemplate ||
			readmeVars.AddIssueTemplate {
			createGitHubFolder = true
		}

		if createGitHubFolder {
			fmt.Println("createGitHubFolder", createGitHubFolder)
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
			case licenseTHEUNL.String():
				if err := k.GenerateTextFromTemplate(licenseFilePath, nil, templateLicenseTHEUNL); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseBSL10.String():
				if err := k.GenerateTextFromTemplate(licenseFilePath, nil, templateLicenseBSL10); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseAPACHE20.String():
				licenseParams := licenseAPACHEVariables{
					FullName: argFullName,
					Year:     now.Year(),
				}

				if err := k.GenerateTextFromTemplate(
					licenseFilePath,
					&licenseParams, templateLicenseAPACHE20); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseMOZP20.String():
				if err := k.GenerateTextFromTemplate(licenseFilePath, nil, templateLicenseMOZP20); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseGNULesserGPL30.String():
				if err := k.GenerateTextFromTemplate(licenseFilePath, nil, templateLicenseGNULesserGPL30); err != nil {
					return fmt.Errorf("could not generate %s file, %w", fnLicense, err)
				}

			case licenseGNUAfferoGPL30.String(), licenseGNUGPL30.String():
				licenseParams := licenseGNUGPL30Variables{
					FullName:    argFullName,
					ProjectName: argProjectName,
					Year:        now.Year(),
				}
				ltemp := templateLicenseGNUAfferoGPL30

				if readmeVars.License == licenseGNUGPL30.String() {
					ltemp = templateLicenseGNUGPL30
				}

				if err := k.GenerateTextFromTemplate(licenseFilePath, &licenseParams, ltemp); err != nil {
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
