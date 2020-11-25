TAG=`git describe --tags`
VERSION ?= `git describe --tags`
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

build = echo "\n\nBuilding $(1)-$(2)" && GOOS=$(1) GOARCH=$(2) go build ${LDFLAGS} -o dist/upload2dav_${VERSION}_$(1)_$(2) \
	&& bzip2 dist/upload2dav_${VERSION}_$(1)_$(2)

upload2dav: *.go
	go build ${LDFLAGS} -o upload2dav
	rm -rf /tmp/go-*

clean:
	rm -f upload2dav

release:
	mkdir -p dist
	rm -f dist/upload2dav_${VERSION}_*
	$(call build,linux,amd64)
	$(call build,linux,386)
	$(call build,linux,arm)
	$(call build,linux,arm64)
	$(call build,darwin,amd64)
