package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func listShareholder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(prefixShareholder, []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var results []interface{}
	for resultsIterator.HasNext() {
		kvResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		c := struct {
			Username     string         `json:"username"`
			TritList     map[string]int `json:"trit_list"`
			FirstName    string         `json:"first_name"`
			LastName     string         `json:"last_name"`
			IdentityCard string         `json:"identity_card"`
		}{}
		err = json.Unmarshal(kvResult.Value, &c)
		if err != nil {
			return shim.Error(err.Error())
		}

		results = append(results, c)
	}

	resultsAsBytes, err := json.Marshal(results)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsAsBytes)
}
func createShareholder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
	if err != nil {

		return shim.Error(fmt.Sprintf("Error extracting creator identity info: %s\n", err.Error()))
	}

	if !authenticateShareholderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of RealEstate Org or Trader Org. Access denied.")
	}

	if len(args) != 1 {
		return shim.Error("func createshareholder Incorrect number of arguments. Expecting 1")
	}

	dto := struct {
		Username     string `json:"username"`      //tên đăng nhập
		FirstName    string `json:"first_name"`    //tên
		LastName     string `json:"last_name"`     //họ
		IdentityCard string `json:"identity_card"` //Số CMT

	}{}
	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	//validate username
	if len(dto.Username) == 0 {
		return shim.Error("Username is empty")
	}

	//validate identity card
	if len(dto.IdentityCard) == 0 {
		return shim.Error("IdentityCard is empty")
	}

	//validate first name
	if len(dto.FirstName) == 0 {
		return shim.Error("FirstName is empty")
	}

	//validate last name
	if len(dto.LastName) == 0 {
		return shim.Error("LastName is empty")
	}

	err = validateShareholderId(stub, dto.Username)
	if err != nil {
		return shim.Error(err.Error())
	}

	shareholder := shareholder{
		Username:     dto.Username,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		IdentityCard: dto.IdentityCard,
	}

	shareholderAsBytes, err := json.Marshal(shareholder)
	if err != nil {
		return shim.Error(err.Error())
	}

	shareholderKey, err := getShareholderKey(stub, dto.Username)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(shareholderKey, shareholderAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func getShareholder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("func getShareholder Invalid argument count. expected 1")
	}

	input := struct {
		Username string `json:"username"`
	}{}

	err := json.Unmarshal([]byte(args[0]), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	userBytes, err := validateShareholderState(stub, input.Username)
	if err != nil {
		return shim.Error(err.Error())
	}

	response := shareholder{}
	err = json.Unmarshal(userBytes, &response)
	if err != nil {
		return shim.Error(err.Error())
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(responseBytes)
}
