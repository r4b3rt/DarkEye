#首次使用
```
brew install protobuf
go get -u github.com/golang/protobuf/protoc-gen-go
```
#编译
```
export PATH=$PATH:$GOPATH/bin
protoc --go_out=. *.proto
```

