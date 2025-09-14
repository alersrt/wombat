#!/usr/bin/env make

env.file ?=

ifdef env.file
	include $(env.file)
	export $(shell sed 's/=.*//' $(env.file))
endif

###############################
# Common defaults/definitions #
###############################

comma := ,

# Checks two given strings for equality.
eq = $(if $(or $(1),$(2)),$(and $(findstring $(1),$(2)),\
                                $(findstring $(2),$(1))),1)




###############
# Git Section #
###############

MAINLINE_BRANCH := dev
CURRENT_BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)

# Squash changes of the current Git branch onto another Git branch.
#
# WARNING: You must merge `onto` branch in the current branch before squash!
#
# Usage:
#	make squash [onto=] [del=(no|yes)]

onto ?= $(MAINLINE_BRANCH)
del ?= no
upstream ?= origin

squash:
ifeq ($(CURRENT_BRANCH),$(onto))
	@echo "--> Current branch is '$(onto)' already" && false
endif
	git checkout $(onto)
	git branch -m $(CURRENT_BRANCH) orig-$(CURRENT_BRANCH)
	git checkout -b $(CURRENT_BRANCH)
	git branch --set-upstream-to $(upstream)/$(CURRENT_BRANCH)
	git merge --squash orig-$(CURRENT_BRANCH)
ifeq ($(del),yes)
	git branch -d orig-$(CURRENT_BRANCH)
endif




###########
# Aliases #
###########

build:
	@make go.build

clean: go.clean

# Resolve all project dependencies.
#
# Usage:
#	make deps

deps: go.deps



###############
# Go commands #
###############

builddir ?= ./build/bin

go.build:
	mkdir -p ${builddir}
	GOARCH=amd64 GOOS=linux go build -a --ldflags='-w -s -extldflags="-static"' --trimpath -o ${builddir}/wombat-linux ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -a --ldflags='-w -s' --trimpath --mod=mod --buildmode=plugin -o ${builddir}/cel-plugin.so ./plugins/cel/main.go
	GOARCH=amd64 GOOS=linux go build -a --ldflags='-w -s' --trimpath --mod=mod --buildmode=plugin -o ${builddir}/imap-plugin.so ./plugins/imap/main.go
	GOARCH=amd64 GOOS=linux go build -a --ldflags='-w -s' --trimpath --mod=mod --buildmode=plugin -o ${builddir}/tg-plugin.so ./plugins/telegram/main.go

go.clean:
	go clean
	rm -rf ${builddir}

go.deps:
	go mod download




###################
# Docker commands #
###################

# Execute docker command with needed params.
#
# Usage:
docker.run:

# Stop project in Docker Compose development environment
# and remove all related containers.
#
# Usage:
#	make docker.down

docker.down:
	CURRENT_UID=$(shell id -u):$(shell id -g) docker compose down --rmi=local -v

# Run Docker Compose development environment.
#
# Usage:
#	make docker.up [rebuild=(yes|no)]
#	               [background=(no|yes)]

rebuild ?= yes
background ?= no

docker.up:
	CURRENT_UID=$(shell id -u):$(shell id -g) docker compose up \
		$(if $(call eq,$(rebuild),no),,--build) \
		$(if $(call eq,$(background),yes),-d,--abort-on-container-exit)




.PHONY: squash \
		go.clean go.build go.deps \
		clean deps build up down \
		docker.up docker.down
