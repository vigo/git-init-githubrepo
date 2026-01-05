![Version](https://img.shields.io/badge/version-0.4.0-orange.svg)
[![golangci-lint](https://github.com/vigo/git-init-githubrepo/actions/workflows/golang-lint.yml/badge.svg)](https://github.com/vigo/git-init-githubrepo/actions/workflows/golang-lint.yml)
[![build and test](https://github.com/vigo/git-init-githubrepo/actions/workflows/go.yml/badge.svg)](https://github.com/vigo/git-init-githubrepo/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/vigo/git-init-githubrepo/branch/main/graph/badge.svg?token=QFA1S8DT00)](https://codecov.io/gh/vigo/git-init-githubrepo)
![Powered by Rake](https://img.shields.io/badge/powered_by-rake-blue?logo=ruby)


# GitHub Friendly Repo Creator/Initializer

Create git repository for GitHub style:

- `README.md` (as seen here!)
- `LICENSE`
- `CODE_OF_CONDUCT.md` (optional)
- `.bumpversion.toml` (optional)
- `SECURITY.md` (optional)
- `.github/CODEOWNERS` (optional)
- `.github/FUNDING.yml` (optional)
- `.github/pull_request_template.md` (optional)

According to `--project-style` (currently only `go` available)

- `.github/workflows/go-test.yml`
- `.github/workflows/go-lint.yml`
- `.github/dependabot.yml`
- `.golangci.yml`

---

## Installation

Install from source;

```bash
go install github.com/vigo/git-init-githubrepo/cmd/git-init-githubrepo@latest
```

or

```bash
brew install vigo/git-init-githubrepo/git-init-githubrepo
```

---

## Usage

You can use with standard git command. `-h`, `--help` or `help` will display
help :)

```bash
$ git init-githubrepo -h

@wip
```

Command fetches some variables from git configuration as default.

- `--full-name`: default is your `git config user.name` if exists
- `--username`: default is your `git config github.user` if exists
- `--email`: default is your `git config user.email` if exists. Email will be used for `CODE_OF_CONDUCT` file.
- `--license`: default license type is `mit`.
- `--disable-license` do not add license information to `README` and do not create `LICENSE` file
- `--disable-fork`: do not add fork information to `README`
- `--disable-bumpversion`: do not create `.bumpversion.cfg` file
- `--disable-coc`: do not create add code of conduct information `README` and do not create `CODE_OF_CONDUCT` file

Required flags are:

- `--project-name`: Name of your project (*title of your project*)
- `--repository-name`: The name you gave when creating the project on GitHub
  (*ex: github.com/USERNAME/REPOSITORYNAME*)

Let’s start a new project. Let’s `cd` to `/tmp`:

```bash
$ git init-githubrepo -p "My Awesome Project" -r "hello-world"
your new project is ready at /tmp/hello-world
$ ls -al /tmp/hello-world/
total 16K
drwxr-xr-x  7 vigo wheel  224 Jun 14 13:15 .
drwxrwxrwt 23 root wheel  736 Jun 14 13:15 ..
drwxr-xr-x  9 vigo wheel  288 Jun 14 13:15 .git
-rwxr-xr-x  1 vigo wheel  182 Jun 14 13:15 .bumpversion.toml
-rwxr-xr-x  1 vigo wheel 3.2K Jun 14 13:15 CODE_OF_CONDUCT.md
-rwxr-xr-x  1 vigo wheel 1.1K Jun 14 13:15 LICENSE.md
-rwxr-xr-x  1 vigo wheel  942 Jun 14 13:15 README.md
```

For bash-completion add:

```bash
eval "$(git-init-githubrepo --bash-completion)"
```

to your bash profile! (*bash completion automatically shipped with brew tap!*)

---

## Contributor(s)

* [Uğur Özyılmazel](https://github.com/vigo) - Creator, maintainer
* [Yiğithan Karabulut](https://github.com/yigithankarabulut) - Contributor

---

## Contribute

All PR’s are welcome!

1. `fork` (https://github.com/vigo/git-init-githubrepo/fork)
1. Create your `branch` (`git checkout -b my-feature`)
1. `commit` yours (`git commit -am 'add some functionality'`)
1. `push` your `branch` (`git push origin my-feature`)
1. Than create a new **Pull Request**!

---

## License

This project is licensed under MIT

---

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].

[coc]: https://github.com/vigo/git-init-githubrepo/blob/main/CODE_OF_CONDUCT.md
