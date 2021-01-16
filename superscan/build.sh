#base
go env -w GOPROXY=https://goproxy.cn,direct
export GO111MODULE=off

ldflag="-s -w"

build_win() {
    GOOS=windows GOARCH=386 go build -ldflags="${ldflag}" -o ../dist/superscan_windows_386.exe
}

build_mac() {
    go build -ldflags="${ldflag}" -o ../dist/superscan_darwin_amd64
}

build_linux() {
     GOOS=linux GOARCH=amd64 go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_amd64
}

build_all() {
    GOOS=darwin GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/superscan_darwin_amd64
    GOOS=freebsd GOARCH=386 go build  -ldflags "${ldflag}" -o ../dist/superscan_freebsd_386
    GOOS=freebsd GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/superscan_freebsd_amd64
    GOOS=linux GOARCH=386 go build -ldflags "${ldflag}" -o ../dist/superscan_linux_386
    GOOS=linux GOARCH=amd64 go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_amd64
    GOOS=linux GOARCH=arm go build -ldflags "${ldflag}" -o ../dist/superscan_linux_arm
    GOOS=linux GOARCH=arm64 go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_arm64
    GOOS=windows GOARCH=386 go build  -ldflags "${ldflag}" -o ../dist/superscan_windows_386.exe
    GOOS=windows GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/superscan_windows_amd64.exe
    GOOS=linux GOARCH=mips64 go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_mips64
    GOOS=linux GOARCH=mips64le go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_mips64le
    GOOS=linux GOARCH=mips GOMIPS=softfloat go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_mips
    GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "${ldflag}" -o ../dist/superscan_linux_mipsle
}

clean() {
    rm -f dic/*.go
    rm -rf db_poc/*.go
}

clean

cd utils && go run dic.go && cd -
cd utils && go run poc.go && cd -

case "$1" in
    "mac")
        build_mac
        ;;
    "linux")
        build_linux
        ;;
    "win")
        build_win
        ;;
     "all")
        build_all
        ;;
     "clean")
        clean
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac

