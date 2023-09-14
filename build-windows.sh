GOOS=windows GOARCH=386 go build -o build/windows-386/subtrans.exe
GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/subtrans.exe
GOOS=windows GOARCH=arm go build -o build/windows-arm/subtrans.exe
GOOS=windows GOARCH=arm64 go build -o build/windows-arm64/subtrans.exe