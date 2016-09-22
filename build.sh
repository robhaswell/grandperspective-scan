#!/bin/bash

compile-gpscan() {
    GOOS=$1 GOARCH=$2 go build -o gpscan && zip build/gpscan-$1-$2.zip gpscan
}

for arch in arm; do {
    compile-gpscan android $arch
} done

for arch in amd64 arm arm64; do {
    compile-gpscan darwin $arch
} done

for arch in amd64; do {
    compile-gpscan dragonfly $arch
} done

for arch in 386 amd64 arm; do {
    compile-gpscan freebsd $arch
} done

for arch in 386 amd64 arm arm64 ppc64 ppc64li mips64 mips64le; do {
    compile-gpscan linux $arch
} done

for arch in 386 amd64 arm; do {
    compile-gpscan netbsd $arch
} done

for arch in 386 amd64 arm; do {
    compile-gpscan openbsd $arch
} done

for arch in 386 amd64; do {
    compile-gpscan plan9 $arch
} done

for arch in amd64; do {
    compile-gpscan solaris $arch
} done

for arch in 386 amd64; do {
    compile-gpscan windows $arch
} done

rm gpscan
