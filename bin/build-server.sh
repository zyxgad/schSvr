#!/bin/sh

cd `dirname $(dirname $0)`

SERVERS_DIR="./bin/server-bin"

function build_linux_go(){
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build "$@"
}

function build_go_file(){
	_dir="${1}"
	_file="${2}"
	_infile="servers/${_dir}/${_file}"
	_outfile="${SERVERS_DIR}/${_dir}"

	echo "\n---Building '${_infile}' out '${_outfile}'"
	echo '=================================================================='
	build_linux_go -o "$_outfile" "$_infile"
	_ecode=$?
	if [ "$_ecode" -ne "0" ]; then
	echo 'Build failed'
		echo '=================================================================='
		echo "exit code: ${_ecode}"
		exit $_ecode
	fi
	echo 'Build succeed'
	echo '=================================================================='
	return 0
}

function build_go_dir(){
	_dirs=$(ls servers)
	for d in $_dirs; do
		[ -d "servers/$d" ] || continue
		build_go_file "$d" "main.go"
	done
}

echo '=================================Building================================='

mkdir -p "${SERVERS_DIR}"

build_go_dir

echo '=================================endBuild================================='


