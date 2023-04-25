
go get -d github.com/therecipe/qt@494d2f8c1ebe65dba7aee361a044511d67147357

```sh
export GO111MODULE=on; go get -v github.com/therecipe/qt && go install -v -tags=no_env github.com/therecipe/qt/cmd/... && go mod vendor && rm -rf vendor/github.com/therecipe/env_linux_amd64_513 && git clone https://github.com/therecipe/env_linux_amd64_513.git vendor/github.com/therecipe/env_linux_amd64_513 && $(go env GOPATH)/bin/qtsetup
```

```sh
export GO111MODULE=on; go install -v -tags=no_env github.com/therecipe/qt/cmd/... && go mod vendor && git clone https://github.com/therecipe/env_linux_amd64_513.git vendor/github.com/therecipe/env_linux_amd64_513 && $(go env GOPATH)/bin/qtdeploy test desktop ./cmd/wallet
```