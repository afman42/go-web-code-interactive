run/api:
	go run main.go

run/web:
	cd web/; npm run dev;

build/linux:
	@echo "Build Binary Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/app main.go 
	@echo "Build Done"

build/windows:
	@echo "Build Binary Windows"
	GOOS=windows GOARCH=amd64 go build -ldflags="-s" -o=./bin/windows_amd64/app main.go 
	@echo "Build Done"

build/web:
	@echo "Build Dist Web"
	cd web/; rm -rf dist/;
	cd web/; npm run build;
	@echo "Build Dist Web Done"

preview/api_linux:
	@echo "Preview"
	./bin/linux_amd64/app

preview/dist_web:
	@echo "Preview Dist Web"
	cd web; npm run preview;

build: build/linux build/windows build/web

