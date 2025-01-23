run/api:
	go run main.go

run/web:
	cd web/; npm run dev;

run/preview_linux: build/linux build/web-staging preview/api_linux

build/rmfldr:
	@echo "Remove Folder";
	rm -rf ./bin/;
	@echo "Finish Remove Folder";

build/linux:
	@echo "Build Binary Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/tmp/app main.go 
	@echo "Build Done"

build/alpine_linux:
	@echo "Build Binary Linux"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/tmp/app main.go 
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


build/compress_linux:
	@echo "Start Compress file linux";
	./upx ./bin/linux_amd64/tmp/app -o  ./bin/linux_amd64/app;
	@echo "Finish Compress file linux";

build/compress_windows:
	@echo "Start Compress file Windows";
	./upx ./bin/windows_amd64/tmp/app -o ./bin/windows_amd64/app;
	@echo "Finish Compress file Windows";

build: build/rmfldr build/linux build/windows build/compress_linux build/compress_windows build/web

build_alpine: build/rmfldr build/alpine_linux build/compress_linux build/web

deploy:
	caprover deploy -h $$CAPROVER_HOST -p $$CAPROVER_PASSWORD -t deploy.tar -a $$CAPROVER_APP_NAME -n $$CAPROVER_MACHINE_NAME;

deploy/tar:
	rm -f deploy.tar;
	tar -zcvf deploy.tar ./bin/linux_amd64/app ./web/dist/ Dockerfile captain-definition .env.prod;

deploy/tar_alpine:
	# Set Up file captain-definition
	rm -f deploy.tar;
	tar -zcvf deploy.tar ./bin/linux_amd64/app ./web/dist/ Dockerfile-alpine captain-definition .env.prod;

deploy/prod_alpine: build_alpine deploy/tar_alpine deploy

deploy/prod: build deploy/tar deploy

.ONESHELL:
npm-install:
	@echo "Install Package Web";
	read -p "Install lib on package.json: " lib;
	cd web/; npm i $$lib;
	@echo "install Package Web done";

npm-uni:
	@echo "Uninstall Package Web";
	read -p "uninstall lib on package.json: " lib;
	cd web/; npm uninstall $$lib;
	@echo "uninstall Package Web done";

preview/api_linux:
	@echo "Preview";
	./bin/linux_amd64/tmp/app -mode preview;
