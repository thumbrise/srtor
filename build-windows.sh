GOOS=windows GOARCH=386 go build -o build/windows-386/srtor.exe
GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/srtor.exe
GOOS=windows GOARCH=arm go build -o build/windows-arm/srtor.exe
GOOS=windows GOARCH=arm64 go build -o build/windows-arm64/srtor.exe