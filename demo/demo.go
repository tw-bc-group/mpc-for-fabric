package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(_ contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) HttpCall(_ contractapi.TransactionContextInterface) (string, error) {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Printf("Failed to http get from baidu, %s", err.Error())
		return "nil", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body from http response: %s", err.Error())
		return "nil", err
	}

	return string(body), nil

}

func main() {
	chainCode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create demo chainCode: %s", err.Error())
		return
	}

	if err := chainCode.Start(); err != nil {
		fmt.Printf("Error starting demo chainCode: %s", err.Error())
	}
}
