cd ..
rmdir /S /Q target
mkdir target
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
cp server/config.json target/datahub.cfg
go build -o target/datahub server/main.go
pause