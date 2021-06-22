NAME ?= hms-bmc-networkprotocol 
VERSION ?= $(shell cat .version)

all : image unittest coverage

image:
		docker build --pull ${DOCKER_ARGS} --tag '${NAME}:${VERSION}' .

unittest: buildbase
		docker build -t cray/hms-bmc-networkprotocol-testing -f Dockerfile.testing .

coverage: buildbase
		docker build -t cray/hms-bmc-networkprotocol-coverage -f Dockerfile.coverage .

buildbase: 
		docker build -t cray/hms-bmc-networkprotocol-build-base -f Dockerfile.build-base .
		
