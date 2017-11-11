#! /bin/sh

set -e
set -u

version=`git describe --tags --always --dirty`

for os in linux darwin freebsd openbsd windows 
do
	for arch in amd64
	do
		build=wsemail-$version-$os-$arch
		mkdir -p ./release/$build
		export GOOS=$os
		export GOARCH=$arch
		export CGO_ENABLED=0
		go build -ldflags "-X main.VersionString=$version" -a -o ./release/$build/wsemail
		cd ./release/
		tar czf $build.tar.gz $build
		zip -q -r $build.zip $build
		cd ..
		rm -rf ./release/$build
	done
done

cd release
sha256sum *.tar.gz *.zip > SHA256