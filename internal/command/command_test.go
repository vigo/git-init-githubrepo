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
	input := []string{
		"--full-name", "Uğur Özyılmazel",
		"--username", "vigo",
		"--email", "ugurozyilmazel@gmail.com",
		"--project-name", "test",
		"--repository-name", "repo",
	}

	args := os.Args[0:1]
	args = append(args, input...)

	cmd, err := command.New()
	if err != nil {
		t.Fatal(err)
	}

	if err := cmd.Run(args); err != nil {
		t.Errorf("want: nil, got: %v", err)
	}

	checkFiles := []string{
		"CODE_OF_CONDUCT.md",
		"LICENSE",
		".bumpversion.toml",
		"README.md",
	}

	for _, file := range checkFiles {
		filePath := strings.Join([]string{tmpFolder, file}, string(os.PathSeparator))
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			t.Errorf("%s does not exists %v", filePath, err)
		}
	}

	if err := os.RemoveAll(tmpFolder); err != nil {
		t.Errorf("can not delete temp folder: %v", err)
	}
}
