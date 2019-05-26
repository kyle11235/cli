# version
VERSION=v0.0.1

# linux
GOOS=linux GOARCH=amd64 go build -o release/$VERSION/linux-amd64/cli

# windows
GOOS=windows GOARCH=amd64 go build -o release/$VERSION/windows-amd64/cli.exe
GOOS=windows GOARCH=386 go build -o release/$VERSION/windows-386/cli.exe

# adrwin
GOOS=darwin GOARCH=amd64 go build -o release/$VERSION/darwin-amd64/cli
