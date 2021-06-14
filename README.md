![Version](https://img.shields.io/badge/version-0.2.0-orange.svg)

# GitHub Friendly Repo Creator/Initializer

Create git repository for GitHub style:

- `README.md` (as seen here!)
- `LICENSE` (optional, currently MIT only)
- `CODE_OF_CONDUCT.md` (optional)
- `.bumpversion.cfg` file injection (optional)

---

## Installation

```bash
go get -u github.com/vigo/git-init-githubrepo
```

---

## Usage

You can use with standard git command. `-h`, `--help` or `help` will display
help :)

```bash
$ git init-githubrepo -h

NAME:
   git-init-githubrepo - create git repository with built-in README, LICENSE and more...

USAGE:
   git-init-githubrepo [global options] command [command options] [arguments...]

VERSION:
   0.1.1

AUTHOR:
   Uğur "vigo" Özyılmazel <ugurozyilmazel@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --bash-completion                  generate bash-completion code (default: false)
   --full-name FULLNAME, -f FULLNAME  your FULLNAME (default: "Uğur Özyılmazel")
   --username USERNAME, -u USERNAME   your GitHub USERNAME (default: "vigo")
   --email EMAIL, -e EMAIL            your contact EMAIL (default: "ugurozyilmazel@gmail.com")
   --project-name NAME, -p NAME       NAME of your project
   --repository-name NAME, -r NAME    NAME of your GitHub repository
   --license LICENSE, -l LICENSE      add LICENSE. available license(s): mit (default: "mit")
   --no-license                       do not add LICENSE file (default: false)
   --disable-fork                     do not add fork information to README (default: false)
   --disable-bumpversion              do not create .bumpversion.cfg and badge to README (default: false)
   --disable-coc                      do not add CODE_OF_CONDUCT (default: false)
   --help, -h                         show help (default: false)
   --version, -v                      version information (default: false)

EXAMPLES:
  
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world"
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc --no-license

NOTES:
  
  Currently, MIT license is available, more to come soon!
```

Command fetches some variables from git configuration as default.

- `--full-name`: default is your `git config user.name` if exists
- `--username`: default is your `git config github.user` if exists
- `--email`: default is your `git config user.email` if exists. Email will be
  used if `--add-coc` is true. (*for code of conduct text*)

- `--add-license`: default is `mit`. If you don’t want to license your project
  you need to pass empty string: `--add-license ""`

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
-rwxr-xr-x  1 vigo wheel  182 Jun 14 13:15 .bumpversion.cfg
-rwxr-xr-x  1 vigo wheel 3.2K Jun 14 13:15 CODE_OF_CONDUCT.md
-rwxr-xr-x  1 vigo wheel 1.1K Jun 14 13:15 LICENSE.md
-rwxr-xr-x  1 vigo wheel  942 Jun 14 13:15 README.md
```

For bash-completion add:

```bash
eval "$(git-init-githubrepo --bash-completion)"
```

to your bash profile!

---

## Contributor(s)

* [Uğur Özyılmazel](https://github.com/vigo) - Creator, maintainer

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
