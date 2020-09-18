go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_mac() {
    ${GOPATH}/bin/qtdeploy build darwin
    cd lowSpeedPortScan && ./build.sh mac && mv lowSpeedPortScan ../dist/
}

build_linux() {
    #docker pull therecipe/qt:linux
    ${GOPATH}/bin/qtdeploy build linux
    cd lowSpeedPortScan && ./build.sh linux && mv lowSpeedPortScan ../dist/
}

build_win() {
    ${GOPATH}/bin/rsrc -manifest DarrmkEye.manifest -ico qml/logo.ico -arch=386 -o DarkEye_windows.syso
    #docker pull therecipe/qt:windows_32_static
    qtdeploy -docker build windows_32_static

    if [[ ! -e deploy/windows ]]; then
       echo "Build Failed"
       return
    else
        echo "Build Success"
    fi
    cd lowSpeedPortScan && ./build.sh win && mv lowSpeedPortScan.exe ../dist/
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
    "linux")
        build_linux
        ;;
      *)
        echo "./build.sh [mac|win|linux]"
        ;;
esac
clean

