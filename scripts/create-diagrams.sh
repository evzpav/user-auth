#!/bin/sh

for format in png; do \
	plantuml \
		-nometadata \
		-r \
		-t$format\
		-v \
		"/data/docs/diagrams/**/*.puml"; \
done
