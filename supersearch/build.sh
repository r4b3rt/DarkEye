#!/bin/bash

build_win() {
    ${GOPATH}/bin/rsrc -manifest DarkEye.manifest -ico qml/logo.ico -arch=386 -o DarkEye_windows.syso
    #docker pull therecipe/qt:windows_32_static
    qtdeploy  -uic=false -docker build windows_32_static

    if [[ ! -e deploy/windows ]]; then
       echo "Build Failed"
       return
    else
        echo "Build Success"
    fi
}

build_mac() {
    ${GOPATH}/bin/qtdeploy  -uic=false build darwin
    if [[ ! -e deploy/darwin ]]; then
       echo "Build Failed"
       return
    else
        echo "Build Success"
    fi
}

clean() {
  rm -f moc*
  rm -f rcc*
  rm -f *.qrc
}

case "$1" in
    "mac")
        build_mac
        ;;
    "win")
        build_win
        ;;
      *)
        echo "./build.sh [mac|win|linux|super]"
        ;;
esac
clean