
# Integration rextporter and skycoin using Docker

## Dockerfile for build rextporter image
```docker
FROM golang:alpine as builder
RUN apk add git
RUN mkdir -p $GOPATH/src/github.com/simelo
RUN cd $GOPATH/src/github.com/simelo && git clone https://github.com/simelo/rextporter.git
RUN cd $GOPATH/src/github.com/simelo/rextporter/cmd/rextporter && go install ./...

FROM alpine:latest as final_layer
WORKDIR /bin
COPY --from=build_layer /go/bin/rextporter .
```
## Docker-compose
```yaml
version: '3.3'
services:
  
  rextporter:
    image: rextporter
    volumes:
      - rextporter-config:/bin/tomlconfig/
    command: "rextporter -config tomlconfig/main.toml"
    ports:
      - "8080:8080"
    
  skycoin:
    image: skycoin/skycoin
    container_name: skycoin
    ports:
      - "6420:6420"
      - "6000:6000"
    volumes:
      - skycoin-data:/data/.skycoin
      - skycoin-wallet:/wallet
      
  

volumes:
  skycoin-data:
    external: true
  skycoin-wallet:
    external: true
  rextporter-config:
    external: true
```
## Steps guide
### 1- Build `rextporter` image
```
docker build -t rextporter /path/to/dockerfile
```
### 2- Create volume `rextporter-config`
```
docker volume create rextporter-config
```
This volume will contain the rextporter config files, if you want to specified a directory containing the config files you can use yhe next command instead of the previous:
```
docker volume create --driver local \
                     --opt type=none \
                     --opt o=bind \
                     --opt device=/path/to/config-directory \
                     rextporter-config
```
### 3- Run docker-compose
```
cd /path/to/docker-compose.yml
docker-compose up
```
