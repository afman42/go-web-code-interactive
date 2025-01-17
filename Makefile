run/api:
	go run main.go

run/web:
	cd web/; npm run dev;

run/preview_linux: build/linux build/web-staging preview/api_linux

build/linux:
	@echo "Build Binary Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/tmp/app main.go 
	@echo "Build Done"

build/windows:
	@echo "Build Binary Windows"
	GOOS=windows GOARCH=amd64 go build -ldflags="-s" -o=./bin/windows_amd64/tmp/app main.go 
	@echo "Build Done"

build/web:
	@echo "Build Dist Web";
	cd web/; rm -rf dist/;
	cd web/; npm run build;
	@echo "Build Dist Web Done";

build/web-staging:
	@echo "Build Dist Web Staging";
	cd ./web/; rm -rf dist/;
	cd ./web/; npm run build:staging;
	@echo "Build Dist Web Staging Done";

preview/api_linux:
	@echo "Preview";
	./bin/linux_amd64/tmp/app -mode preview;

build: build/linux build/windows build/compress build/web

deploy:
	caprover deploy -t deploy.tar;

deploy/tar:
	rm -f deploy.tar;
	tar -zcvf deploy.tar ./bin/linux_amd64/app ./web/dist/ Dockerfile captain-definition .env.prod;

.ONESHELL:
npm-install:
	@echo "Install Package Web";
	read -p "Install lib on package.json: " lib;
	cd web/; npm i $$lib;
	@echo "install Package Web done";
	
build/compress:
	@echo "Start Compress file";
	./upx ./bin/linux_amd64/tmp/app -o  ./bin/linux_amd64/app;
	./upx ./bin/windows_amd64/tmp/app -o ./bin/windows_amd64/app;
	@echo "Finish Compress file";
