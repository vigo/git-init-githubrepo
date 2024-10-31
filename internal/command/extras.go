package command

var (
	extrasBashCompletion = `_git_init_githubrepo() {
	  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
	    local cur opts base
	    COMPREPLY=()
	    cur="${COMP_WORDS[COMP_CWORD]}"
	    if [[ "$cur" == "-"* ]]; then
	      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
	    else
	      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
	    fi
	    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
	    return 0
	  fi
	}
	`

	extrasHelp = `
AVALILABLE LICENSES (%d):

%s

EXAMPLES:

  $ git init-githubrepo -p "My Awesome Project" -r "hello-world"
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc --no-license
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --license gnu-agpl30
  $ git init-githubrepo -p "My Awesome Project" -r "hello-world" --license moz-p20

`

	extrasAppUsage = `create GitHub friendly git repository with built-in README, LICENSE and more...`
)
