package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getShareholderKey(stub shim.ChaincodeStubInterface, shareholderID string) (string, error) {
	shareholderKey, err := stub.CreateCompositeKey(prefixShareholder, []string{shareholderID})
	if err != nil {
		return "", err
	} else {
		return shareholderKey, nil
	}
}

func getPublisherKey(stub shim.ChaincodeStubInterface, publisherId string) (string, error) {
	publisherKey, err := stub.CreateCompositeKey(prefixPublisher, []string{publisherId})
	if err != nil {
		return "", err
	} else {
		return publisherKey, nil
	}
}

func getTransferContractKey(stub shim.ChaincodeStubInterface, contractID string) (string, error) {
	contractKey, err := stub.CreateCompositeKey(prefixTransferContract, []string{contractID})
	if err != nil {
		return "", err
	} else {
		return contractKey, nil
	}
}

func getPublishContractKey(stub shim.ChaincodeStubInterface, contractId string) (string, error) {
	contractKey, err := stub.CreateCompositeKey(prefixPublishContract, []string{contractId})
	if err != nil {
		return "", err
	} else {
		return contractKey, nil
	}
}

func getChangePriceContractKey(stub shim.ChaincodeStubInterface, contractId string) (string, error) {
	contractKey, err := stub.CreateCompositeKey(prefixChangePriceContract, []string{contractId})
	if err != nil {
		return "", err
	} else {
		return contractKey, nil
	}
}

func getRealEstateKey(stub shim.ChaincodeStubInterface, realEstateID string) (string, error) {
	contractKey, err := stub.CreateCompositeKey(prefixRealEstate, []string{realEstateID})
	if err != nil {
		return "", err
	} else {
		return contractKey, nil
	}
}

func getSellTritAdvertisingContractKey(stub shim.ChaincodeStubInterface, contractId string) (string, error) {
	contractKey, err := stub.CreateCompositeKey(prefixSellTritAdvertisingContract, []string{contractId})
	if err != nil {
		return "", err
	} else {
		return contractKey, nil
	}
}
