export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on

Version="v5.0.1"
ldFlag="-s -w -X main.Version=${Version}"
topDir=$(pwd)
AppName=superscan

build_win() {
    GOOS=windows GOARCH=386    go build -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_windows_386.exe
    GOOS=windows GOARCH=amd64  go build -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_windows_amd64.exe
}

build_mac() {
    go build -ldflags="${ldflag}" -o "${topDir}"/dist/"${AppName}"_darwin_amd64
}

build_linux() {
    GOOS=linux GOARCH=amd64  go build  -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_linux_amd64
}

build_others() {
    GOOS=linux GOARCH=arm   go build -ldflags "${ldflag}" -o ../dist/df_linux_arm
    GOOS=linux GOARCH=arm64 go build -ldflags "${ldflag}" -o ../dist/df_linux_arm64

    GOOS=linux GOARCH=mips64  go build  -ldflags "${ldflag}" -o ../dist/df_linux_mips64
    GOOS=linux GOARCH=mips64le  go build  -ldflags "${ldflag}" -o ../dist/df_linux_mips64le
    GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "${ldflag}" -o ../dist/df_linux_mips
    GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "${ldflag}" -o ../dist/df_linux_mipsle
}

build_all() {
    build_mac
    build_linux
    build_win
}

build_dict() {
  rm -f "${topDir}"/dict/dict.go && cd "${topDir}"/dict && go-bindata -ignore .DS_Store -pkg dict -o "${topDir}"/dict/dict.go ./...
}

compress() {
    upx -9 df_windows_*
    upx -9 df_linux_*
   # upx -9 df_darwin_*
}

case "$1" in
    "dict")
        build_dict
        ;;
    "mac")
        build_mac
        ;;
    "linux")
        build_linux
        ;;
    "win")
        build_win
        ;;
    "win")
        build_others
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac


