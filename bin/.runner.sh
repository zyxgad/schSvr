#!/bin/bash

RUN_DIR="./bin/server-bin"

respath="${PWD}/mount"

_ohost=''
_oport=''
function read_config(){
	file="${respath}/.config/servers/${1}.txt"
	if [[ '1' == '1' ]]; then
		read _ohost
		read _oport
	fi <"${file}"
}

function run_exe(){
	echo "================================================"

	_name=$1
	_exe="${RUN_DIR}/${_name}"
	shift

	if ! [ -f "${_exe}" ]; then
		echo "File not exist '${_exe}'"
		return -1
	fi

	_dkname="schsvr_${_name}"
	read_config "${_name}"
	echo "Running ${_exe} -> '${_dkname}' (${_ohost}:${_oport})"
	docker create \
	--rm -it \
	--name "$_dkname" \
	--env "SCH_SVR_MODE" \
	--network schbsdnet \
	--network-alias "${_ohost}" \
	--volume "${respath}":"/var/server" \
	"$@" \
	"ubuntu:latest" \
	"/root/main" \
	|| return $?
	docker cp "${_exe}" "${_dkname}":"/root/main"
	docker start "$_dkname"
}


# for f in `ls "${RUN_DIR}"`; do
# 	if [ -f "${RUN_DIR}/$f" ]; then
# 		run_exe "$f"
# 	fi
# done

echo "Run mode: ${SCH_SVR_MODE}"
echo

echo 'Initing network'
docker network create --subnet=172.20.0.0/16 schbsdnet
echo '================================'

log_dir="${PWD}/logs"
[ -d "${log_dir}" ] || mkdir "${log_dir}"
run_exe logger --volume "${log_dir}":"/var/server_logs"
read_config 'sql'
run_exe sql --publish "${_oport}":"${_oport}"
run_exe web
read_config 'proxy'
run_exe proxy --publish "${_oport}":"${_oport}"

