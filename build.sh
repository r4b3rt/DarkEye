#base

#brew reinstall FiloSottile/musl-cross/musl-cross --with-x86_64  --with-mips --with-i486 --with-aarch64  --with-arm-hf  --with-mipsel --with-mips64 --with-mips64el
#brew install mingw-w64

go env -w GOPROXY=https://goproxy.cn,direct
export GO111MODULE=off

#Using go-sqlite3
export CGO_ENABLED=1
ldflag="-s -w"

support_rdp() {
    os=$1

    case $os in
        mac)
#build in mac
#cmake -D "CMAKE_OSX_ARCHITECTURES:STRING=x86_64" -DBUILD_SHARED_LIBS=OFF  -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/Volumes/dev/gosrc/src/github.com/zsdevX/FreeRDP/formac
#make
            path=${GOPATH}/src/github.com/zsdevX/freerdp_binary/formac
            export CGO_CFLAGS="-I${path}/include/freerdp3 -I${path}/include/winpr3 -DRDP_SUPPORT"
            export CGO_LDFLAGS="${path}/lib/libfreerdp3.a $path/lib/libwinpr3.a $path/lib/libfreerdp-client3.a $path/libcrypto.a $path/libssl.a"
            ;;
        linux)
#build in alpine:
#docker pull alpine:latest
#apk add git make cmake wget libstdc++ gcc g++ libuv-dev openssl-dev hwloc-dev linux-headers musl-dev
#cmake -DBUILD_SHARED_LIBS=OFF  -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/forlinux
#make install
            path=${GOPATH}/src/github.com/zsdevX/freerdp_binary/forlinux
            export CGO_CFLAGS="-I${path}/include/freerdp2 -I${path}/include/winpr2 -DRDP_SUPPORT -static"
            export CGO_LDFLAGS="${path}/lib64/libfreerdp2.a $path/lib64/libwinpr2.a $path/lib64/libfreerdp-client2.a $path/libssl.a $path/libcrypto.a"
            ;;
        win32)
            ;;
        win64)
#build in mac+mingw
#step1:  git clone https://github.com/vityank/FreeRDP.git -b mingw
#step2:  注释 winpr/include/winpr/spec.h:#define DECLSPEC_EXPORT 这个会和mingw的dlldecl冲突，太坑了。
#Step3:  openssl: /Configure mingw64  no-asm no-async no-engine --prefix=/Volumes/dev/gosrc/src/github.com/zsdevX/openssl-OpenSSL_1_1_1k/install CROSS_COMPILE=x86_64-w64-mingw32-
#Step4:  freerdp: cmake -DCMAKE_SYSTEM_NAME="Windows" -DCMAKE_C_COMPILER="x86_64-w64-mingw32-gcc" -DCMAKE_CXX_COMPILER="x86_64-w64-mingw32-g++" -DCMAKE_INSTALL_PREFIX=/Volumes/dev/gosrc/src/github.com/zsdevX/FreeRDP/forwin64  -DOPENSSL_ROOT_DIR=/Volumes/dev/gosrc/src/github.com/zsdevX/openssl-OpenSSL_1_1_1k -DCMAKE_EXE_LINKER_FLAGS="-L/usr/local/Cellar/mingw-w64/8.0.0/toolchain-x86_64/x86_64-w64-mingw32/lib" -DWIN32=1 -DUWP=1 -DBUILD_SHARED_LIBS=OFF ..
            path=${GOPATH}/src/github.com/zsdevX/freerdp_binary/forwin64
            export CGO_CFLAGS="-I${path}/include/freerdp3 -I${path}/include/winpr3 -DRDP_SUPPORT"
            export CGO_LDFLAGS="${path}/lib/libfreerdp3.a $path/lib/libwinpr3.a $path/lib/libfreerdp-client3.a $path/libcrypto.a $path/libssl.a -lws2_32 -limagehlp -lshlwapi"
            ;;

        *)
            export CGO_CFLAGS="-static"
            export CGO_LDFLAGS=
            echo ${os}":不支持rdp"
        ;;
    esac
}


build_win() {
    cd console
    support_rdp none
    GOOS=windows GOARCH=386 CC="i686-w64-mingw32-gcc"  go build -a  -ldflags "${ldflag}" -o ../dist/df_windows_386.exe
    GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc"  go build -a -ldflags "${ldflag}" -o ../dist/df_windows_amd64.exe
    cd -
}

build_mac() {
    support_rdp mac
    cd console
    go build -ldflags="${ldflag}" -o ../dist/df_darwin_amd64
    cd -
}

build_linux() {
    support_rdp linux
    cd console
    GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc  go build -a  -ldflags "${ldflag}" -o ../dist/df_linux_amd64
    cd -
}

build_all() {
    build_mac
    build_linux
    build_win

    support_rdp none
    cd console
#arm
    GOOS=linux GOARCH=arm  CC=arm-linux-musleabihf-gcc CGO_LDFLAGS="-static" go build -a -ldflags "${ldflag}" -o ../dist/df_linux_arm
    GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc CGO_LDFLAGS="-static" go build -a  -ldflags "${ldflag}" -o ../dist/df_linux_arm64
#mip[sel][64]
    GOOS=linux GOARCH=mips CC=mips-linux-musl-cc CGO_LDFLAGS="-static" GOMIPS=softfloat go build -a -ldflags "${ldflag}" -o ../dist/df_linux_mips
    GOOS=linux GOARCH=mipsel  CC=mipsel-linux-musl-cc CGO_LDFLAGS="-static" GOMIPS=softfloat go build -a -ldflags "${ldflag}" -o ../dist/df_linux_mipsel

    GOOS=linux GOARCH=mips64 CC=mips64-linux-musl-cc CGO_LDFLAGS="-static" go build -a  -ldflags "${ldflag}" -o ../dist/df_linux_mips64
    GOOS=linux GOARCH=mips64el CC=mips64el-linux-musl-cc CGO_LDFLAGS="-static" go build -a  -ldflags "${ldflag}" -o ../dist/df_linux_mips64el
    cd -
}

prepare() {
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
   # upx -9 df_darwin_*
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
        clean
        prepare
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

