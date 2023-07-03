export NAME := ip-monitor

build:
	goreleaser build --snapshot --clean

run:
	goreleaser build --single-target --snapshot --clean
	dist/${NAME}*/${NAME} -c .config.json
