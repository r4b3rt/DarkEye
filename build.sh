#go env -w GOPROXY=https://goproxy.cn,direct
#GOOS=windows GOARCH=386 go build -ldflags="-s -w"
export GO111MODULE=off

build_super() {
    cd superscan && ./build.sh all
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
}

build_linux() {
     os=`uname`
     if [[ "$os" == "Linux" ]]; then
            qtdeploy  -uic=false build linux
            mv deploy/linux/DarkEye dist/
     else
        echo "Note:"
        echo "#macbook下编译linux QT界面比较复杂，但是按照如下操作是成功的:
            #加载docker方式加载qt编译环境（注意映射代码到docker中）： docker run -idt -v /Volumes/dev/gosrc:/Volumes/dev/src therecipe/qt:linux
            #进入容器：docker exec -it qt-linux-docker bash
            #配置环境：export PATH=$PATH:/opt/Qt5.13.0/5.13.0/gcc_64/bin
            #配置环境：export GOPATH=/Volumes/dev/src
            #编译：qtdeploy  -uic=false build linux
            "
     fi
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
}

clean() {
    rm -f moc*
    rm -f rcc*
    rm -f dark_eye.cfg*
    rm -f *.cpp
    rm -rf darwin
    rm -rf windows
    cd superscan && ./build.sh clean && cd -
    find . -type f -name "*.go"|xargs gofmt -s -w
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
    "super")
        build_super
        ;;
    "clean")
        clean
        ;;
      *)
        echo "./build.sh [mac|win|linux|super]"
        ;;
esac


