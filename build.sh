go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_mac() {
    ${GOPATH}/bin/qtdeploy build darwin
    if [[ ! -e deploy/darwin ]]; then
       echo "Build Failed"
       return
    else
        mv deploy/darwin/DarkEye.app dist/
        echo "Build Success"
    fi
    cd portscan && ./build.sh mac && mv portscan ../dist/portscan.mac
}

build_linux() {
    #docker pull therecipe/qt:linux
    #${GOPATH}/bin/qtdeploy build linux
    cd portscan && ./build.sh linux && mv portscan ../dist/portscan.linux64
}

build_win() {
    ${GOPATH}/bin/rsrc -manifest DarkEye.manifest -ico qml/logo.ico -arch=386 -o DarkEye_windows.syso
    #docker pull therecipe/qt:windows_32_static
    qtdeploy -docker build windows_32_static

    if [[ ! -e deploy/windows ]]; then
       echo "Build Failed"
       return
    else
        mv deploy/windows/DarkEye.exe dist/
        echo "Build Success"
    fi
    cd portscan && ./build.sh win && mv portscan.exe ../dist/
}

clean() {
    rm -f moc*
    rm -f rcc*
    rm -f dark_eye.cfg*
    rm -f *.cpp
    rm -rf darwin
    rm -rf windows
    rm -rf portscan/portscan
    rm -rf portscan/tmp
}

clean
case "$1" in
    "mac")
        build_mac
        ;;
    "win")
        build_win
        ;;
    "linux")
        build_linux
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac
clean

