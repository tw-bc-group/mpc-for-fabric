#!/usr/bin/env bash

function run_as_org1() {
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

function run_as_org2()
{
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
}

heaas_container_name="heaas-server"
docker stop $heaas_container_name > /dev/null 2>&1
docker container rm $heaas_container_name > /dev/null 2>&1

if [[ ! -d "./fabric-samples" ]]; then
  curl -sSL https://bit.ly/2ysbOFE | bash -s
fi

docker image inspect $heaas_container_name > /dev/null 2>&1

if [[ "$?" -ne "0" ]]; then
  echo "local heaas_server image not found, starting building it."
  docker build -t $heaas_container_name heaas_server/.
fi

cd "fabric-samples/test-network" || exit

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true

./network.sh down
./network.sh up createChannel

echo "start heaas-server and join network..."
network="net_test"
docker run --name $heaas_container_name --network $network -p 10000:10000 -d $heaas_container_name

echo "packaging demo..."
peer lifecycle chaincode package demo.tar.gz --path ../../demo --lang golang --label demo_1


echo "install demo on org1..."
run_as_org1
peer lifecycle chaincode install demo.tar.gz


echo "install demo on org2..."
run_as_org2
peer lifecycle chaincode install demo.tar.gz

package=$(peer lifecycle chaincode queryinstalled | sed -n 2p | sed 's/^Package ID: \(.*\),.*$/\1/g')

export export CC_PACKAGE_ID=$package


echo "approve demo on org2..."
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name demo --version 1.0 --package-id "$CC_PACKAGE_ID" --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "approve demo on org1..."
run_as_org1
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name demo --version 1.0 --package-id "$CC_PACKAGE_ID" --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "commit demo to test network..."
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name demo --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
