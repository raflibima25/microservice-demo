## Mac

```
brew install protobuf
```

## Ubuntu/Debian

```
apt-get install protobuf-compiler
```

## Windows (dengan chocolatey)

```
choco install protoc
```

## Install Go plugins untuk protoc

```
make install-tools
```

## Generate gRPC code dari proto files:

```
make proto
```

## Terminal 1 - Auth Service

```
make run-auth
```

## Terminal 2 - Product Service

```
make run-product
```

## Terminal 3 - Frontend

```
make run-frontend
```

## Terminal 4 - API Gateway

```
make run-api-gateway
```

## Import GRPC

```
go mod edit -require=grpc@v0.0.0
go mod edit -replace=grpc=../grpc
go mod tidy
```

## Install and run redis on windows

```
https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-windows/
```
