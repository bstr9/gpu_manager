# vim:set noet

default: build

lint:
	prototool lint gpu_manager.proto

build: build-pb

build-pb:
	protoc --go_out=plugins=grpc:proto gpu_api.proto
