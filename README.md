## Deploy local fabric test network

### Prerequisites

- Git

- cURL

- Docker and Docker Compose

- Golang

  

#### Note

If you are using Docker Toolbox or macOS, you will need to use a location under `/Users` (macOS) when installing and running the samples.

If you are using Docker for Mac, you will need to use a location under `/Users`, `/Volumes`, `/private`, or `/tmp`. To use a different location, please consult the Docker documentation for [file sharing](https://docs.docker.com/docker-for-mac/#file-sharing)


## Run Demo

###### !!! IMPORTANT !!! 

Modify demo/demo.go line 24 assign the server address, e.g. the docker0 bridge address.

```bash
./deploy.sh
```

```bash
./setAsk.sh
```

```
./setBids.sh
```

```bash
./getBids.sh
```