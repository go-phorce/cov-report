#!/usr/bin/env bash
# A simple wrapper around the golint command line reference in the Makefile as a dependency.
#
# The tool itself will eventually filter out files tagged with 'DO NOT EDIT' as tracked by
# the issue here: https://github.com/golang/lint/issues/239
#
# Usage: golint.sh <golint-command-line>
# Example: golint.sh .tools/bin/golint -set_exit_status salesforce.com/... testonly/...
#
# The intent is for this script to be transient in nature until golint actual sorts out how
# to either ignore swaths of files or by-default ignores generated files (or alternatively
# if the code-gen tool stops writing 'bad go').
#
# See: Makefile: lint:
#
bin_cut="`\which cut`"
bin_grep="`\which grep`"
grep_opts="--extended-regexp"
pattern_generated='^[\s]*//.*DO NOT EDIT'

function _analyze() {
  if [ ${ret} -eq 0 ]; then
    # golint passed; do not analyze
    return
  else
    # if you set IFS to nil; you'll keep token characters in your lines
    failed=0
    generated=0
    while IFS= read -r l; do
      f=$(printf ${l} |${bin_cut} -d: -f1)
      if [ -f "${f}" ]; then
        ${bin_grep} ${grep_opts} "${pattern_generated}" ${f} >/dev/null 2>&1
        grep_ret=$?
        if [ 0 -eq ${grep_ret} ]; then
          # match found; generated code
          (( generated +=1 ))
        else
          (( failed += 1 ))
          printf "${l}\n" # stdout
        fi
      else
        printf "${l}\n" 1>&2  # stderr
      fi
      export failed generated
    done <<< "${out_err}"
    printf "Found %d false-positives from generated code" "${generated}" 2>&1 # stderr
    if [ 0 -lt ${failed} ]; then
      printf ".\n" 2>&1 # stderr
      # fail
      ret=1
    else
      printf "; unfailing.\n" 2>&1 # stderr
      ret=0
    fi
  fi
}

function _golint() {
  cli="${@}"
  printf "%s\n" "${cli}"
  out_err=$(eval ${cli} 2>&1)
  ret=$?
}

function main() {
  _golint ${@}
  _analyze
  exit ${ret}
}

main ${@}