export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on

Version="v1.0.1"
ldFlag="-s -w -X main.Version=${Version}"
AppName=forensics
topDir=$(pwd)

build_win() {
    rsrc -manifest qml/forensics.manifest -ico qml/logo.ico -o forensics_windows.syso
    GOOS=windows GOARCH=386    go build -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_windows_386.exe
    GOOS=windows GOARCH=amd64  go build -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_windows_amd64.exe
}

build_mac() {
    go build -ldflags="${ldflag}" -o "${topDir}"/dist/"${AppName}"_darwin_amd64
}

build_linux() {
    GOOS=linux GOARCH=amd64  go build  -ldflags "${ldflag}" -o "${topDir}"/dist/"${AppName}"_linux_amd64
}

build_mac
build_win
build_linux