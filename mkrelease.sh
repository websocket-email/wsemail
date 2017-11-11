#! /bin/sh

set -e
set -u

version=`git describe --always --dirty`

for os in linux darwin freebsd openbsd windows 
do
	for arch in amd64
	do
		build=wsemail-$version-$os-$arch
		builddir=./release/$build
		mkdir -p $builddir
		export GOOS=$os
		export GOARCH=$arch
		export CGO_ENABLED=0
		go build -ldflags "-X main.VersionString=$version" -a -o $builddir/wsemail
		cd release
		tar czf $build.tar.gz $build
		zip -q -r $build.zip $build
		cd ..
		rm -rf $builddir
	done
done


cd release
sha256sum *.tar.gz *.zip > SHA256