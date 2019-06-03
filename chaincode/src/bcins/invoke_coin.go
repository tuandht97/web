package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func payIn(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("func payIn Incorrect number of arguments. Expecting 1")
	}

	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
	if err != nil {

		return shim.Error(fmt.Sprintf("Error extracting creator identity info: %s\n", err.Error()))
	}

	if !authenticateRegulatorOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Regulator Org. Access denied.")
	}

	dto := struct {
		Amount   int    `json:"amount"`   // Số lượng coin
		Receiver string `json:"receiver"` // Người nạp
	}{}

	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	//validate amount
	if dto.Amount <= 0 {
		return shim.Error("Amount must be greater than zero")
	}

	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiver, err := getShareholderState(stub, dto.Receiver)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiver.Balance += dto.Amount
	err = receiver.saveShareHolderState(stub)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func payInByShareholder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("func payInByShareholder Incorrect number of arguments. Expecting 1")
	}

	creatorOrg, creatorCertIssuer, username, err := getTxCreatorInfo(stub)

	if err != nil {

		return shim.Error(fmt.Sprintf("Error extracting creator identity info: %s\n", err.Error()))
	}

	if !authenticateShareholderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of RealEstate Org or Trader Org. Access denied.")
	}

	dto := struct {
		Amount int `json:"amount"` // Số lượng coin
	}{}

	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	//validate amount
	if dto.Amount <= 0 {
		return shim.Error("Amount must be greater than zero")
	}

	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiver, err := getShareholderState(stub, username)
	if err != nil {
		return shim.Error(err.Error())
	}

	receiver.Balance += dto.Amount
	err = receiver.saveShareHolderState(stub)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
