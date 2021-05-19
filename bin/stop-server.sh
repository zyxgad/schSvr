#!/bin/bash


for name in `docker ps --filter name=schsvr_ --format='{{.Names}}'`; do
	echo "Stoping '${name}'"
	docker stop "${name}"
	echo '================================================'
done


