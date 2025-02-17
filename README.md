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
