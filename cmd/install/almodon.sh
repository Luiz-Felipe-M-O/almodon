#!/bin/bash

set -eu

log() {
	printf '%s\n' "$*" >&2
}

fatal() {
	printf 'fatal: %s\n' "$*" >&2
	exit 1
}

need() {
	command -v "$1" >/dev/null 2>&1 || fatal "$2 '$1' is not available in PATH"
}

machine64() {
	bits=$(getconf LONG_BIT 2>/dev/null || printf '')
	if [ "$bits" = "64" ]; then
		return 0
	fi

	arch=$(uname -m 2>/dev/null || printf '')
	if [[ "$arch" =~ (64|^(aarch64|arm64|amd64|x86_64|ppc64|ppc64le|s390x|riscv64))$ ]]; then
		return 0
	fi

	return 1
}

version() {
	re_version="go([0-9]+)\.([0-9]+)\.([0-9]+)"

	if [[ "$1" =~ $re_version ]]; then
		major="${BASH_REMATCH[1]}"
		minor="${BASH_REMATCH[2]}"
		patch="${BASH_REMATCH[3]}"

		if [[ "$2" =~ $re_version ]]; then
			reqmajor="${BASH_REMATCH[1]}"
			reqminor="${BASH_REMATCH[2]}"
			reqpatch="${BASH_REMATCH[3]}"

			if [ "$major" -lt "$reqmajor" ]; then
				echo 'major'
				return 1
			fi
			if [ "$minor" -lt "$reqminor" ]; then
				echo 'minor'
				return 1
			fi
			if [ "$patch" -lt "$reqpatch" ]; then
				echo 'patch'
				return 1
			fi

			return 0
		fi
	fi

	return 1
}

ccompiler() {
	need $1 'C Compiler'
	target=$("$1" -dumpmachine 2>/dev/null || printf '')
	if [[ "$target" =~ ^(x86_64|aarch64|arm64|ppc64|ppc64le|s390x|riscv64) ]]; then
		return 0
	fi
	
	return 1
}

log       'checking target machine architecture...'
machine64 || fatal 'target machine must be 64-bit'

go_version="$(go version | awk '{print $3}')"
log     'checking Go toolchain...'
need    go 'Go compiler'
version "$go_version" 'go1.26.0' || fatal "go1.26.0 or newer is required, found $go_version"

cc="$(go env CC)"
log       "checking C compiler for cgo..."
ccompiler "$cc" || fatal "C compiler '$cc' does not appear to target 64-bit architectures"

log  "checking for git..."
need git 'git'

repo_dir='almodon'
log  "checking that ./$repo_dir does not already exist..."
test ! -e "$repo_dir" || fatal "folder '$repo_dir' already exists in the current directory"

repo_url='https://github.com/alan-b-lima/almodon'
log "cloning $repo_url..."
git clone "$repo_url" "$repo_dir" || fatal "failed to clone $repo_url"
cd  "$repo_dir" || fatal "failed to enter $repo_dir"

binary_path='./bin/almodon'
log   "building $binary_path from ./cmd/main.go..."
export CGO_ENABLED=1
mkdir  -p ./bin || fatal 'failed to create bin directory'
go     build -o "$binary_path" -trimpath -ldflags='-s -w -linkmode external -extldflags "-static"' ./cmd/main.go || fatal 'static build failed'

log "installation completed successfully!"

