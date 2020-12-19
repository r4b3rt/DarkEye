go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off
#set -x
build_win() {
    GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ../dist/superscan_windows_386.exe
}

build_mac() {
    go build -ldflags="-s -w" -o ../dist/superscan_darwin_amd64
}

build_linux() {
     GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_arm64
}

build_all() {
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_darwin_amd64
    CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_freebsd_386
    CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_freebsd_amd64
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_386
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_amd64
    CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_arm
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_arm64
    CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_windows_386.exe
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_windows_amd64.exe
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_mips64
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_mips64le
    CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_mips
    CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -trimpath -ldflags "-s -w" -o ../dist/superscan_linux_mipsle
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

