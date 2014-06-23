VERSION=0.1
REVISION=1

ARCH=i386

default: all

all: compile test

compile:
	go build -v

test: 
	go test -v ./...

install:
	go install

package: deb rpm

deb:
	fpm -s dir -t deb \
	--prefix /usr/bin \
	--name stopgap-metrics \
	--version ${VERSION} \
	--iteration ${REVISION} \
	--architecture amd64 \
	--maintainer pat.downey@gmail.com \
	--description 'stopgap metrics graphite json' \
	stopgap

rpm:
	fpm -s dir -t rpm \
	--prefix /usr/bin \
	--name stopgap-metrics \
	--version ${VERSION} \
	--iteration ${REVISION} \
	--epoch 0 \
	--architecture i386 \
	--rpm-user 0 \
	--rpm-group 0 \
	--rpm-os linux \
	--maintainer pat.downey@gmail.com \
	--description 'stopgap metrics graphite json' \
	stopgap

