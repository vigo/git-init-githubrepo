package command_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/vigo/git-init-githubrepo/internal/command"
	"github.com/vigo/git-init-githubrepo/internal/version"
)

var (
	tmpDir    = strings.TrimRight(os.TempDir(), string(os.PathSeparator))
	tmpFolder = strings.Join([]string{tmpDir, "repo"}, string(os.PathSeparator))
)

func TestBashCompletion(t *testing.T) {
	input := []string{
		"--bash-completion",
	}

	args := os.Args[0:1]
	args = append(args, input...)

	out := new(bytes.Buffer)
	cmd, err := command.New(
		command.WithWriter(out),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Run(args); err != nil {
		t.Errorf("want: nil, got: %v", err)
	}

	got := string(bytes.TrimSpace(out.Bytes()))
	if !strings.Contains(got, "_git_init_githubrepo") {
		t.Errorf("want: contains _git_init_githubrepo, got: %v", got)
	}
}

func TestListLicences(t *testing.T) {
	input := []string{
		"--list-licenses",
	}

	args := os.Args[0:1]
	args = append(args, input...)

	out := new(bytes.Buffer)
	cmd, err := command.New(
		command.WithWriter(out),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Run(args); err != nil {
		t.Errorf("want: nil, got: %v", err)
	}

	got := string(bytes.TrimSpace(out.Bytes()))
	if !strings.Contains(got, "mit") {
		t.Errorf("want: mit, got: %v", got)
	}
	if !strings.Contains(got, "mit-na") {
		t.Errorf("want: mit-na, got: %v", got)
	}
}

func Test(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  string
		err   error
	}{
		{
			name:  "check version",
			input: []string{"-v"},
			want:  version.Version,
			err:   nil,
		},
		{
			name:  "run w/o project name",
			input: nil,
			want:  "",
			err:   command.ErrProjectNameRequired,
		},
		{
			name:  "run w/o repo name arg",
			input: []string{"--project-name", "test"},
			want:  "",
			err:   command.ErrRepositoryNameRequired,
		},
		{
			name: "run with wrong license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "notexist",
			},
			want: "",
			err:  command.ErrInvalidLicense,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			args := os.Args[0:1]
			args = append(args, testCase.input...)

			out := new(bytes.Buffer)
			cmd, err := command.New(
				command.WithWriter(out),
			)
			if err != nil {
				t.Fatal(err)
			}

			if err := cmd.Run(args); !errors.Is(err, testCase.err) {
				t.Errorf("want: %v, got: %v", testCase.err, err)
			}

			got := string(bytes.TrimSpace(out.Bytes()))
			if testCase.want != got {
				t.Errorf("want: %v, got: %v", testCase.want, got)
			}

			if err := os.RemoveAll(tmpFolder); err != nil {
				t.Errorf("can not delete temp folder: %v", err)
			}
		})
	}
}

func TestCreateAll(t *testing.T) {
	testCases := []struct {
		name            string
		input           []string
		lookupInLicense string
		checkFiles      []string
		err             error
	}{
		{
			name: "create with default (mit) license",
			input: []string{
				"--full-name", "Uğur Özyılmazel",
				"--username", "vigo",
				"--email", "ugurozyilmazel@gmail.com",
				"--project-name", "test",
				"--repository-name", "repo",
			},
			lookupInLicense: "The MIT License",
			checkFiles: []string{
				"CODE_OF_CONDUCT.md",
				"LICENSE",
				".bumpversion.toml",
				"README.md",
			},
			err: nil,
		},
		{
			name: "create with apache-20 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "apache-20",
			},
			lookupInLicense: "Apache License",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with bsl-10 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "bsl-10",
			},
			lookupInLicense: "Boost Software License",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with gnu-agpl30 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "gnu-agpl30",
			},
			lookupInLicense: "GNU AFFERO GENERAL PUBLIC LICENSE",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with gnu-gpl30 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "gnu-gpl30",
			},
			lookupInLicense: "GNU GENERAL PUBLIC LICENSE",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with gnu-lgpl30 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "gnu-lgpl30",
			},
			lookupInLicense: "GNU LESSER GENERAL PUBLIC LICENSE",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with mit-na license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "mit-na",
			},
			lookupInLicense: "MIT No Attribution",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with moz-p20 license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "moz-p20",
			},
			lookupInLicense: "Mozilla Public License Version 2.0",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
		{
			name: "create with unli license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "unli",
			},
			lookupInLicense: "This is free and unencumbered",
			checkFiles: []string{
				"LICENSE",
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			args := os.Args[0:1]
			args = append(args, testCase.input...)

			cmd, err := command.New()
			if err != nil {
				t.Fatal(err)
			}

			if err := cmd.Run(args); err != nil {
				t.Errorf("want: nil, got: %v", err)
			}

			if testCase.checkFiles != nil {
				for _, file := range testCase.checkFiles {
					filePath := strings.Join([]string{tmpFolder, file}, string(os.PathSeparator))
					if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
						t.Errorf("%s does not exists %v", filePath, err)
					}
				}
			}

			if testCase.lookupInLicense != "" {
				filePath := strings.Join([]string{tmpFolder, "LICENSE"}, string(os.PathSeparator))

				data, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("can not open file: %v", err)
				}

				if !strings.Contains(string(data), testCase.lookupInLicense) {
					t.Errorf("LICENSE does not contain: %s", testCase.lookupInLicense)
				}
			}

			if err := os.RemoveAll(tmpFolder); err != nil {
				t.Errorf("can not delete temp folder: %v", err)
			}
		})
	}
}
