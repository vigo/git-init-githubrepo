package main

var (
	bashCompletion = `#! /bin/bash

: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
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

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete $PROG
unset PROG
`

	helpExtras = `
EXAMPLES:
	
	$ git init-githubrepo -p "My Awesome Project" -r "hello-world"
	$ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork
	$ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion
	$ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc
	$ git init-githubrepo -p "My Awesome Project" -r "hello-world" --disable-fork --disable-bumpversion --disable-coc --no-license

NOTES:
	
	Currently, MIT license is available, more to come soon!

	`
)
