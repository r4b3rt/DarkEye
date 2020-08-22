go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_mac() {
    ${GOPATH}/bin/qtdeploy build darwin
}

build_win() {
    ${GOPATH}/bin/rsrc -manifest DarkEye.manifest -ico qml/logo.ico -arch=386 -o DarkEye_windows.syso
    #docker pull therecipe/qt:windows_32_static
    qtdeploy -docker build windows_32_static

    if [[ ! -e deploy/windows ]]; then
       echo "Build Failed"
       return
    else
        echo "Build Success"
    fi
}

clean() {
    rm -f moc*
    rm -f rcc*
    rm -f dark_eye.cfg*
    rm -f *.cpp
    rm -rf darwin
    rm -rf windows
}

clean
case "$1" in
    "mac")
        build_mac
        ;;
    "win")
        build_win
        ;;
      *)
        echo "./build.sh [mac|win]"
        ;;
esac
clean

