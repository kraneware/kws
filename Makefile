SHELL := /bin/bash

TEST_PACKAGES = $(shell find . -name "*_test.go" | sort | rev | cut -d'/' -f2- | rev | uniq)
CURDIR = $(shell pwd)

.DEFAULT_GOAL := package

.PHONY: clean init displavars test coverage coverage-checks

clean:
	@rm -Rf target

init: clean
	@mkdir target
	@mkdir target/testing
	@mkdir target/bin
	@mkdir target/deploy
	@mkdir target/tools

deps: init
	go get -d && go mod tidy

displayvars:
	@for package in $(TEST_PACKAGES); do \
		echo $${package:2}; \
	done

cleanup:
	gofmt -w .
	$(GOPATH)/bin/goimports -w .

test: init
	@for package in $(TEST_PACKAGES); do \
	  echo Testing package $$package ; \
	  cd $(CURDIR)/$$package ; \
	  mkdir -p ${CURDIR}/target/testing/$$package ; \
	  go test -v -race -covermode=atomic -coverprofile ${CURDIR}/target/testing/$$package/coverage.out | tee ${CURDIR}/target/testing/$$package/target.txt ; \
	  if [ "$${PIPESTATUS[0]}" -ne "0" ]; then exit 1; fi; \
	  grep "FAIL!" ${CURDIR}/target/testing/$$package/target.txt ; \
	  if [ "$$?" -ne "1" ]; then exit 1; fi; \
	  cat ${CURDIR}/target/testing/$$package/coverage.out >> ${CURDIR}/target/coverage_profile.out ; \
	done

coverage: test
	@for package in ${TEST_PACKAGES}; do \
	  export MIN_COVERAGE=95 ; \
	  echo Generating coverage report for $$package ; \
	  cd $(CURDIR)/$$package ; \
	  if [ -f test.config ]; then source test.config; fi; \
	  go tool cover -html=${CURDIR}/target/testing/$$package/coverage.out -o ${CURDIR}/target/testing/$$package/coverage.html ; \
	done

coverage-checks: coverage
	@for package in ${TEST_PACKAGES}; do \
	  export MIN_COVERAGE=100 ; \
	  cd $(CURDIR)/$$package ; \
	  if [ -f test.config ]; then source ./test.config; fi; \
	  echo Checking coverage for $$package at $$MIN_COVERAGE% ; \
	  export COVERAGE_PCT=`grep "coverage: " ${CURDIR}/target/testing/$$package/target.txt | cut -d' ' -f2` ; \
	  export COVERAGE=`echo $$COVERAGE_PCT | cut -d'.' -f1` ; \
	  if [ "$$COVERAGE" -lt "$$MIN_COVERAGE" ]; then echo - Coverage not met at $$COVERAGE_PCT. ; exit 1; fi ; \
	  echo "  Coverage passed with $$COVERAGE_PCT" ; \
	done

build: coverage-checks
	@if [ -f lambda-deploy.json ]; then \
	  echo Building lambda target/bin/`cat lambda-deploy.json | python3 -c 'import json,sys;print(json.load(sys.stdin)["handler"])'` ... ; \
	  env GOOS=linux GOARCH=amd64 go build -o target/bin/`cat lambda-deploy.json | python3 -c 'import json,sys;print(json.load(sys.stdin)["handler"])'` ; \
	fi

package: build
	@if [ -f lambda-deploy.json ]; then \
	  echo Packaging lambda target/deploy/`cat lambda-deploy.json | python3 -c 'import json,sys;print(json.load(sys.stdin)["handler"])'`.zip ... ; \
	  zip -D target/deploy/`cat lambda-deploy.json | python3 -c 'import json,sys;print(json.load(sys.stdin)["handler"])'`.zip lambda-deploy.json ; \
	  cd target/bin && zip -Du ../deploy/`cat ../../lambda-deploy.json | python3 -c 'import json,sys;print(json.load(sys.stdin)["handler"])'`.zip * ; \
	fi