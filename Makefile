GO = go
BIN = onsengo
VERSION != git describe --tags
WRKDIR= ./build

output = ${WRKDIR}/${BIN}-${VERSION}

cross-build: fbsd/amd64 mac/amd64 linux/amd64 win/amd64

bindist: fbsd/amd64_dist mac/amd64_dist linux/amd64_dist win/amd64_dist

fbsd/amd64:
	mkdir -p ${WRKDIR}
	GOOS=freebsd GOARCH=amd64 ${GO} build -o ${output}-freebsd-amd64

fbsd/amd64_dist: fbsd/amd64
	zstd -19 ${output}-freebsd-amd64

mac/amd64:
	mkdir -p ${WRKDIR}
	GOOS=darwin GOARCH=amd64 ${GO} build -o ${output}-darwin-amd64

mac/amd64_dist: mac/amd64
	zstd -19 ${output}-darwin-amd64

linux/amd64:
	mkdir -p ${WRKDIR}
	GOOS=linux GOARCH=amd64 ${GO} build -o ${output}-linux-amd64

linux/amd64_dist: linux/amd64
	zstd -19 ${output}-linux-amd64

win/amd64:
	mkdir -p ${WRKDIR}
	GOOS=windows GOARCH=amd64 ${GO} build -o ${output}-windows-amd64.exe

win/amd64_dist: win/amd64
	zip ${output}-windows-amd64.exe.zip ${output}-windows-amd64.exe

clean:
	rm -rf ${WRKDIR}
