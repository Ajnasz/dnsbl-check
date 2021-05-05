#!/bin/sh

set -eu

build() {
	VERSION=$(git describe --tags)
	BUILD=$(date +%FT%T%z)
	ARCHLIST="$1"
	OSLIST="$2"
	BUILD_DIR="$3"
	FILE_NAME="dnsbl-check"

	echo "$VERSION"
	echo "$BUILD"

	mkdir -p "$BUILD_DIR"

	for os in $OSLIST
	do
		for arch in $ARCHLIST
		do
			echo "building $os.$arch"
			if go tool dist list | grep -q "^${os}/${arch}$"
			then
				GOOS="$os" GOARCH="$arch" go build -ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}" -o "$BUILD_DIR/$FILE_NAME.$os.$arch"
			fi
		done
	done
}

remove() {
  BUILD_DIR="$1"
	for os in $OSLIST
	do
		for arch in $ARCHLIST
		do
			if [ -f "$BUILD_DIR/$FILE_NAME.$os.$arch" ]
			then
				echo "removing $os.$arch"
				rm "$BUILD_DIR/$FILE_NAME.$os.$arch"
			fi
		done
	done
}

REMOVE=0
OSLIST="linux darwin"
ARCHLIST="amd64 arm64 arm"
BUILD_DIR="./build"
BUILD=0

if [ $# -lt 1 ];then
	echo "Command required, available commands are \"test\" and \"build\"" >&2
	exit 1
fi

subcommand="$1"
shift
case "$subcommand" in
	"test")
		go test -v -tags test
		return
		;;
	"build")
		BUILD=1;
		;;
	*)
		echo "Invalid command, available commands are \"test\" and \"build\"" >&2
		exit 1
		;;
esac

while getopts "ra:o:s:b:" opt
do
	case "$opt" in
		"r")
			REMOVE=1
			;;
		"a")
			ARCHLIST="$OPTARG"
			;;
		"o")
			OSLIST="$OPTARG"
			;;
		"b")
			BUILD_DIR="$OPTARG"
			;;
		[?])
			exit 1
			;;
	esac
done

if [ $REMOVE -eq 1 ];then
	remove "$BUILD_DIR"
else
	build "$ARCHLIST" "$OSLIST" "$BUILD_DIR"
fi
