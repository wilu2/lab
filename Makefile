GIT_TAG = $(shell git rev-parse --short=8 HEAD)

build_ko:
	export KO_DOCKER_REPO=registry.intsig.net/textin_gateway/data_panel && ko build ./cmd/api-server --bare --tags=$(GIT_TAG)
