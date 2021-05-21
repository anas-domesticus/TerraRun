test:
	docker build -t terrarun-builder --target tester .
	docker run terrarun-builder go test ./...

lint:
	docker build -t terrarun-builder --target tester .
	docker run terrarun-builder staticcheck ./...

build:
	docker build -t terrarun-builder --target builder .
	docker run -v $(PWD):/src/bin terrarun-builder cp /src/terrarun /src/bin/terrarun

build-mac:
	docker build --build-arg GOOS=darwin --build-arg GOARCH=amd64 -t terrarun-builder --target builder .
	docker run -v $(PWD):/src/bin terrarun-builder cp /src/terrarun /src/bin/terrarun

build-mac-arm:
	docker build --build-arg GOOS=darwin --build-arg GOARCH=arm64 -t terrarun-builder --target builder .
	docker run -v $(PWD):/src/bin terrarun-builder cp /src/terrarun /src/bin/terrarun