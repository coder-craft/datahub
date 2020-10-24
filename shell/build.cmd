cd ..
rmdir /S /Q target
mkdir target
cp server/datahub/config.json target/datahub.cfg
go build -o target/datahub.exe server/datahub/main.go
pause