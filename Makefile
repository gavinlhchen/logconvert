BINARY=yjtosocserver
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
yjtosocserverPath=./cmd/yjtosocserver/

flags='-X "logconvert/version.GitVersion=${gitVersion}" \
-X "logconvert/version.BuildDate=${buildTime}" \
-X "logconvert/version.GitCommit=${gitHash}" \
-X "logconvert/version.GitTreeState=${gitTreeState}"'

build:
	go build -mod=mod -ldflags ${flags} -o ${BINARY} ${yjtosocserverPath}
	rm -rf ${DIR_BIN}
	mkdir -p ${DIR_BIN} && mv ${BINARY} ${DIR_BIN}

pkg:
	tar -zcvf ${BINARY}${DATETIME}.tar.gz ${DIR_BIN} ${DIR_CONF}

install:
	rm -rf ${DIR_INSTALL}
	mkdir -p ${DIR_INSTALL}
	cp -r ./configs ${DIR_INSTALL}
	go build -mod=mod -ldflags ${flags} -o ${DIR_INSTALL}/bin/${BINARY} ${yjtosocserverPath}

debug:
	go build -mod=mod -gcflags "-N -l" -ldflags ${flags} -o ${BINARY}

test:
	go test -v ./...
#Cleans our projects: deletes binaries
clean:
	rm -f *.tar.gz ${DIR_BIN}/*
.PHONY: clean install
