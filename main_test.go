package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/vigo/git-init-githubrepo/version"
)

var (
	tmpDir    = strings.TrimRight(os.TempDir(), string(os.PathSeparator))
	tmpFolder = strings.Join([]string{tmpDir, "repo"}, string(os.PathSeparator))
)

func Test(t *testing.T) {
	tcs := []struct {
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
			name:  "run w/o required args",
			input: nil,
			want:  "",
			err:   errProjectNameRequired,
		},
		{
			name:  "run w/o repo name arg",
			input: []string{"--project-name", "test"},
			want:  "",
			err:   errRepoNameRequired,
		},
		{
			name: "run with wrong license",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
				"--license", "notexist",
			},
			want: "",
			err:  errInvalidLicense,
		},
		{
			name: "run with no coc",
			input: []string{
				"--project-name", "test",
				"--repository-name", "repo",
			},
			want: "",
			err:  errEmailRequired,
		},
		{
			name: "create under temp",
			input: []string{
				"--full-name", "full name",
				"--username", "username",
				"--email", "test@email.com",
				"--project-name", "test",
				"--repository-name", "repo",
			},
			want: "your new project is ready at " + tmpFolder,
			err:  nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			args := os.Args[0:1]
			args = append(args, tc.input...)
			out := new(bytes.Buffer)

			if err := run(args, out); !errors.Is(err, tc.err) {
				t.Errorf("want: %v, got: %v", tc.err, err)
			}

			got := string(bytes.TrimSpace(out.Bytes()))
			if tc.want != got {
				t.Errorf("want: %v, got: %v", tc.want, got)
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
	out := new(bytes.Buffer)

	if err := run(args, out); err != nil {
		t.Errorf("want: nil, got: %v", err)
	}

	checkFiles := []string{
		"CODE_OF_CONDUCT.md",
		"LICENSE",
		".bumpversion.cfg",
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
