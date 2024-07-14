build:
	go build -ldflags "-X xgrpc.XGrpcServerBuildVersion=$(date -u +%Y-%m-%d-%H:%M:%S)" -o build/server cmd/main.go
	./build/server -v

clean-build:
	rm build/*

build-version:
	./build/server -v