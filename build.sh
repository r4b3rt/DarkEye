#go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_hackTools() {
    os=$1
    cd superscan && ./build.sh $os && mv superscan ../../dist/portscan.${os}
}

build_mac() {
    ${GOPATH}/bin/qtdeploy  -uic=false build darwin
    if [[ ! -e deploy/darwin ]]; then
       echo "Build Failed"
       return
    else
        rm -rf dist/DarkEye.app
        mv deploy/darwin/DarkEye.app dist/
        echo "Build Success"
    fi
    build_hackTools "mac"
}

build_linux() {
     os=`uname`
     if [[ "$os" == "Linux" ]]; then
            #docker run -idt -v /Volumes/dev/gosrc:/Volumes/dev/src therecipe/qt:linux
            #docker exec -it qt-linux-docker bash
            #export PATH=$PATH:/opt/Qt5.13.0/5.13.0/gcc_64/bin
            #export GOPATH=/Volumes/dev/src
            qtdeploy  -uic=false build linux
            mv deploy/linux/DarkEye dist/
     fi
     build_hackTools "linux"
}

build_win() {
    ${GOPATH}/bin/rsrc -manifest DarkEye.manifest -ico qml/logo.ico -arch=386 -o DarkEye_windows.syso
    #docker pull therecipe/qt:windows_32_static
    qtdeploy  -uic=false -docker build windows_32_static

    if [[ ! -e deploy/windows ]]; then
       echo "Build Failed"
       return
    else
        mv deploy/windows/DarkEye.exe dist/
        echo "Build Success"
    fi
     build_hackTools "windows"
}

clean() {
    rm -f moc*
    rm -f rcc*
    rm -f dark_eye.cfg*
    rm -f *.cpp
    rm -rf darwin
    rm -rf windows
    rm -rf hacktools/portscan/portscan
    rm -rf hacktools/portscan/tmp
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

