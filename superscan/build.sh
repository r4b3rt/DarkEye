go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_linux() {
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
}

build_win() {
    GOOS=windows GOARCH=386 go build -ldflags="-s -w"
}

build_mac() {
    go build -ldflags="-s -w"
}

clean() {
    rm -f superscan
    rm -f superscan.exe
    rm -f *.csv
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
     "clean")
        clean
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac

