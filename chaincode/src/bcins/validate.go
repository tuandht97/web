package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func validatePublisherState(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	key, err := getPublisherKey(stub, id)
	if err != nil {
		return nil, err
	}

	publisherAsBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if publisherAsBytes == nil {
		return nil, errors.New(fmt.Sprintf("publisher with this username %s does not exist.", id))
	}

	return publisherAsBytes, nil
}

func validateShareholderState(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	key, err := getShareholderKey(stub, id)
	if err != nil {
		return nil, err
	}

	shareholderAsBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if shareholderAsBytes == nil {
		return nil, errors.New(fmt.Sprintf("shareholder with this username %s does not exist.", id))
	}

	return shareholderAsBytes, nil
}

func validateRealEstateState(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	key, err := getRealEstateKey(stub, id)
	if err != nil {
		return nil, err
	}

	realEstateAsBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if realEstateAsBytes == nil {
		return nil, errors.New(fmt.Sprintf("real estate with this id %s does not exist.", id))
	}

	return realEstateAsBytes, nil
}

func validateSellTritAdvertisingContractState(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	key, err := getSellTritAdvertisingContractKey(stub, id)
	if err != nil {
		return nil, err
	}

	contractAsBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if contractAsBytes == nil {
		return nil, errors.New(fmt.Sprintf("Sell Trit Advertising Contract with this id %s does not exist.", id))
	}

	return contractAsBytes, nil
}

func validateRealEstateId(stub shim.ChaincodeStubInterface, id string) error {
	realEstateKey, err := getRealEstateKey(stub, id)
	if err != nil {
		return err
	}

	realEstateAsBytes, err := stub.GetState(realEstateKey)
	if realEstateAsBytes != nil {
		return errors.New("The id of this Real Estate is taken, try another " + id)
	}
	return nil
}

func validateSellTritAdvertisingContractId(stub shim.ChaincodeStubInterface, id string) error {
	contractKey, err := getSellTritAdvertisingContractKey(stub, id)
	if err != nil {
		return err
	}

	contractAsBytes, err := stub.GetState(contractKey)
	if err != nil {
		return err
	}
	if contractAsBytes != nil {
		return errors.New(fmt.Sprintf("Advertising Contract with this UUID: " + id + " is taken"))
	}

	return nil
}

func validateTransferContractId(stub shim.ChaincodeStubInterface, id string) error {
	contractKey, err := getTransferContractKey(stub, id)
	if err != nil {
		return err
	}

	contractAsBytes, err := stub.GetState(contractKey)
	if err != nil {
		return err
	}
	if contractAsBytes != nil {
		return errors.New(fmt.Sprintf("Transfer Contract with this UUID: " + id + " is taken"))
	}

	return nil
}

func validateChangePriceContractId(stub shim.ChaincodeStubInterface, id string) error {
	contractKey, err := getChangePriceContractKey(stub, id)
	if err != nil {
		return err
	}

	contractAsBytes, err := stub.GetState(contractKey)
	if err != nil {
		return err
	}
	if contractAsBytes != nil {
		return errors.New(fmt.Sprintf("Change Price Contract with this UUID: " + id + " is taken"))
	}

	return nil
}

func validatePublishContractId(stub shim.ChaincodeStubInterface, id string) error {
	contractKey, err := getPublishContractKey(stub, id)
	if err != nil {
		return err
	}

	contractAsBytes, err := stub.GetState(contractKey)
	if err != nil {
		return err
	}
	if contractAsBytes != nil {
		return errors.New(fmt.Sprintf("Publish Contract with this UUID: " + id + " is taken"))
	}

	return nil
}

func validateShareholderId(stub shim.ChaincodeStubInterface, id string) error {
	key, err := getShareholderKey(stub, id)
	if err != nil {
		return err
	}

	userAsBytes, err := stub.GetState(key)
	if userAsBytes != nil {
		return errors.New("The id of this shareholder is taken, try another " + id)
	}
	return nil
}

func validatePublisherId(stub shim.ChaincodeStubInterface, id string) error {
	key, err := getPublisherKey(stub, id)
	if err != nil {
		return err
	}

	userAsBytes, err := stub.GetState(key)
	if userAsBytes != nil {
		return errors.New("The id of this publisher is taken, try another " + id)
	}
	return nil
}
