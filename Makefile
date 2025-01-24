.ONESHELL:

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
	@if [ $$l = "alpine" ]; then
		@echo "Build Binary $l linux";
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/tmp/app main.go
		@echo "Build Binary $l Done";
		return
	fi
	@echo "Build Binary linux";
	GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/tmp/app main.go;
	@echo "Build Binary Done";

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

build: 
	@make build/rmfldr; 
	@make build/linux l="etc";
	@make build/windows; 
	@make build/compress_linux; 
	@make build/compress_windows; 
	@make build/web;

build_alpine: 
	@make build/rmfldr; 
	@make build/linux l="alpine"; 
	@make build/compress_linux; 
	@make build/web;

deploy:
	caprover deploy -h $$CAPROVER_HOST -p $$CAPROVER_PASSWORD -t deploy.tar -a $$CAPROVER_APP_NAME -n $$CAPROVER_MACHINE_NAME;

deploy/tar:
	rm -f deploy.tar;
	tar -zcvf deploy.tar ./bin/linux_amd64/app ./web/dist/ Dockerfile captain-definition .env.prod;

deploy/tar_alpine:
	rm -f deploy.tar;
	sed -i 's/Dockerfile/Dockerfile-alpine/g' ./captain-definition
	tar -zcvf deploy.tar ./bin/linux_amd64/app ./web/dist/ Dockerfile-alpine captain-definition .env.prod;
	sed -i 's/Dockerfile-alpine/Dockerfile/g' ./captain-definition

deploy/prod_alpine: build_alpine deploy/tar_alpine deploy

deploy/prod: build deploy/tar deploy

npm-install:
	@echo "Install Package Web";
	read -p "Install lib on package.json: " lib;
	cd web/; npm i $$lib;
	@echo "install Package Web done";

npm-uni:
	@echo "Uninstall Package Web";
	read -p "uninstall lib on package.json: " lib
	cd web/; npm uninstall $$lib;
	@echo "uninstall Package Web done";

preview/api_linux:
	@echo "Preview";
	./bin/linux_amd64/tmp/app -mode preview;
