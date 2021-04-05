#base
go env -w GOPROXY=https://goproxy.cn,direct
export GO111MODULE=off

ldflag="-s -w"

build_win() {
    cd console
    GOOS=windows GOARCH=386 go build -ldflags="${ldflag}" -o ../dist/df_windows_386.exe
    GOOS=windows GOARCH=amd64 go build -ldflags="${ldflag}" -o ../dist/df_windows_amd64.exe
    cd -
}

build_mac() {
    cd console
    go build -ldflags="${ldflag}" -o ../dist/df_darwin_amd64
    cd -
}

build_linux() {
    cd console
    GOOS=linux GOARCH=amd64 go build  -ldflags "${ldflag}" -o ../dist/superscan_linux_amd64
    cd -
}

build_all() {
    cd console

    GOOS=darwin GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/df_darwin_amd64
    GOOS=freebsd GOARCH=386 go build  -ldflags "${ldflag}" -o ../dist/df_freebsd_386
    GOOS=freebsd GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/df_freebsd_amd64
    GOOS=linux GOARCH=386 go build -ldflags "${ldflag}" -o ../dist/df_linux_386
    GOOS=linux GOARCH=amd64 go build  -ldflags "${ldflag}" -o ../dist/df_linux_amd64
    GOOS=linux GOARCH=arm go build -ldflags "${ldflag}" -o ../dist/df_linux_arm
    GOOS=linux GOARCH=arm64 go build  -ldflags "${ldflag}" -o ../dist/df_linux_arm64
    GOOS=windows GOARCH=386 go build  -ldflags "${ldflag}" -o ../dist/df_windows_386.exe
    GOOS=windows GOARCH=amd64 go build -ldflags "${ldflag}" -o ../dist/df_windows_amd64.exe
    GOOS=linux GOARCH=mips64 go build  -ldflags "${ldflag}" -o ../dist/df_linux_mips64
    GOOS=linux GOARCH=mips64le go build  -ldflags "${ldflag}" -o ../dist/df_linux_mips64le
    GOOS=linux GOARCH=mips GOMIPS=softfloat go build  -ldflags "${ldflag}" -o ../dist/df_linux_mips
    GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "${ldflag}" -o ../dist/df_linux_mipsle

    cd -
}

prepare() {
    clean

    cd superscan/utils

    go run dic.go

    cd -
}

clean() {
    cd superscan
    rm -f dic/*.go
    cd -

    rm -f console/console
    rm -f console/analysis.s3db
    rm -rf dist/*
}
compress() {
    cd  dist
    upx -9 df_windows_*
    upx -9 df_linux_*
    upx -9 df_darwin_*
    cd -
}

prepare
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
        compress
        ;;
     "clean")
        clean
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac

