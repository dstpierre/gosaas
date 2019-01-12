#!/bin/sh

LATEST_STABLE_SUPPORTED_GO_VERSION="1.11"

main() {
  if local_go_version_is_latest_stable
  then
    run_gofmt
    run_golint
    run_govet
  fi
}

local_go_version_is_latest_stable() {
  go version | grep -q $LATEST_STABLE_SUPPORTED_GO_VERSION
}

log_error() {
  echo "$*" 1>&2
}

run_gofmt() {
  GOFILES_NOVENDOR=$(find . -type f -name '*.go' -not -path "./vendor/*")
  GOFMT_FILES=$(gofmt -l $GOFILES_NOVENDOR)
  if [ -n "$GOFMT_FILES" ]
  then
    log_error "gofmt failed for the following files:
$GOFMT_FILES

please run 'gofmt -w .' on your changes before committing."
    exit 1
  fi
}

run_golint() {
  GOFILES_NOVENDOR=$(find . -type f -name '*.go' -not -path "./vendor/*")
  GOLINT_ERRORS=$(golint $GOFILES_NOVENDOR)
  if [ -n "$GOLINT_ERRORS" ]
  then
    log_error "golint failed for the following reasons:
$GOLINT_ERRORS

please run 'golint ./...' on your changes before committing."
    exit 1
  fi
}

run_govet() {
  GOVET_ERRORS=$(go tool vet ./*.go 2>&1)
  if [ -n "$GOVET_ERRORS" ]
  then
    log_error "go vet failed for the following reasons:
$GOVET_ERRORS

please run 'go tool vet ./*.go' on your changes before committing."
    exit 1
  fi
}

main