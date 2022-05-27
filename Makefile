DIR_BIN=./bin/
CONF=./configs/
DIR_INSTALL=../../YujianInstaller/services/logcollect/logconvert
goVersion=$(shell go version)
gitVersion=$(shell git version)
buildTime=$(shell date +"%Y-%m-%d %H:%m:%S")
gitHash=$(shell git rev-parse HEAD)
gitTreeState=clean
ifneq ($(strip ${shell git status --porcelain}),)
        gitTreeState=dirty
endif

flags='-X "logconvert/version.GitVersion=${gitVersion}" \
-X "logconvert/version.BuildDate=${buildTime}" \
-X "logconvert/version.GitCommit=${gitHash}" \
-X "logconvert/version.GitTreeState=${gitTreeState}"'

build: build_logcollecttool build_yjtosocserver

build_yjtosocserver: mkdir_bin
	go build -mod=mod -ldflags ${flags} -o yjtosocserver ./cmd/yjtosocserver/
	mv yjtosocserver ${DIR_BIN}

build_logcollecttool: mkdir_bin
	go build -mod=mod -ldflags ${flags} -o logcollecttool ./cmd/logcollecttool/
	mv logcollecttool ${DIR_BIN}

mkdir_bin:
	rm -rf ${DIR_BIN}
	mkdir -p ${DIR_BIN}

pkg:
	tar -zcvf logcollect.tar.gz ${DIR_BIN} ${DIR_CONF}

install: install_logcollecttool install_yjtosocserver

install_yjtosocserver:install_mkdir
	go build -mod=mod -ldflags ${flags} -o ${DIR_INSTALL}/bin/yjtosocserver ./cmd/yjtosocserver/

install_logcollecttool:install_mkdir
	go build -mod=mod -ldflags ${flags} -o ${DIR_INSTALL}/bin/logcollecttool ./cmd/logcollecttool/

install_mkdir:
	rm -rf ${DIR_INSTALL}
	mkdir -p ${DIR_INSTALL}
	cp -r ./configs ${DIR_INSTALL}

test:
	go test -v ./...
#Cleans our projects: deletes binaries
clean:
	rm -f *.tar.gz ${DIR_BIN}/*
.PHONY: clean install
