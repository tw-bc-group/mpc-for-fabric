## Deploy local fabric test network

### Prerequisites

- Git

- cURL

- Docker and Docker Compose

- Golang

  

#### Note

If you are using Docker Toolbox or macOS, you will need to use a location under `/Users` (macOS) when installing and running the samples.

If you are using Docker for Mac, you will need to use a location under `/Users`, `/Volumes`, `/private`, or `/tmp`. To use a different location, please consult the Docker documentation for [file sharing](https://docs.docker.com/docker-for-mac/#file-sharing)



### Prepare Tools

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s
```

The command above will do these things:

1. If needed, clone the [hyperledger/fabric-samples](https://github.com/hyperledger/fabric-samples) repository
2. Checkout the appropriate version tag
3. Install the Hyperledger Fabric platform-specific binaries and config files for the version specified into the /bin and /config directories of fabric-samples
4. Download the Hyperledger Fabric docker images for the version specified



```
cd fabric-samples/test-network
```



#### Note

The $PWD of next steps always fabric-samples/test-network


### Prepare Demo

```
cd demo
go mod vendor
```


### Bring up the test network

```bash
./network.sh up createChannel
```



### Bring down the test network

```
./network.sh up down
```



## Deploy demo chaincode to test network

mv demo and deploy.sh to test-network's directory

```bash
./deploy.sh
```

every time run this scipt will renew a test-network and deploy the newest demo chaincode
