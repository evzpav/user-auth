#!/bin/bash

docker run --rm -d \
		--name ${1} \
		--network ${2} \
		mongo:3.6 \
		--smallfiles --noprealloc --nojournal
