#!/bin/bash

RUN_EXE="./bin/server-bin"

function _help(){
	echo "Help page"
	echo "-B [--rebuild]: rebuild bin/server-main"
	echo "-m [--run-mode]: 'debug' or 'release'"
	exit 0
}

REBUILD=false
RESTART=false
MODE=debug

while [[ $1 != "" ]]; do
	case $1 in
		-B | --rebuild )
			REBUILD=true
			;;
		-h | --help )
			_help
			;;
		-m | --run-mode )
			shift
			MODE=$1
			;;
		-r | --restart )
			RESTART=true
			;;
	esac
	shift
done

cd `dirname $(dirname $0)`

[ -d "${RUN_EXE}" ] || REBUILD=true

if [[ "$REBUILD" == 'true' ]]; then
	./bin/build-server.sh || exit $?
fi


export SCH_SVR_MODE=RELEASE
if [[ "$MODE" == "debug" ]]; then
	export SCH_SVR_MODE=DEBUG
fi

exec "./bin/.runner.sh"

