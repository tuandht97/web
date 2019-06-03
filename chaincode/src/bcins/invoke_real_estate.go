package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func listRealEstate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	resultsIterator, err := stub.GetStateByPartialCompositeKey(prefixRealEstate, []string{})
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

		c := realEstate{}
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

func createRealEstate(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	creatorOrg, creatorCertIssuer, _, err := getTxCreatorInfo(stub)
	if err != nil {

		return shim.Error(fmt.Sprintf("Error extracting creator identity info: %s\n", err.Error()))
	}

	if !authenticateRegulatorOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Regulator Org. Access denied. " + creatorCertIssuer + " " + creatorOrg)
	}

	if len(args) != 1 {
		return shim.Error("func createRealEstate Incorrect number of arguments. Expecting 1")
	}

	dto := struct {
		Id          string  `json:"id"`           // mã của bất động sản - ví dụ CHA, CHB
		Price       int     `json:"price"`        // giá trị
		SquareMeter float64 `json:"square_meter"` // diện tích
		Address     string  `json:"address"`      // Địa chỉ
		Description string  `json:"description"`  // Miêu tả
		OwnerId     string  `json:"owner_id"`     // id của chủ sở hữu bất động sản
	}{}

	err = json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	//validate id
	if len(dto.Id) == 0 {
		return shim.Error("Id is empty")
	}

	//validate price
	if dto.Price <= 0 {
		return shim.Error("Price of real estate must be greater than zero")
	}

	//validate owner state
	_, err = validateShareholderState(stub, dto.OwnerId)
	if err != nil {
		return shim.Error(err.Error())
	}

	//validate new real estate id
	err = validateRealEstateId(stub, dto.Id)
	if err != nil {
		return shim.Error(err.Error())
	}

	estate := realEstate{
		Id:           dto.Id,
		Shareholders: nil,
		Price:        dto.Price,
		SquareMeter:  dto.SquareMeter,
		Address:      dto.Address,
		Description:  dto.Description,
		OwnerId:      dto.OwnerId,
	}

	realEstateAsBytes, err := json.Marshal(estate)
	if err != nil {
		return shim.Error(err.Error())
	}

	realEstateKey, err := getRealEstateKey(stub, dto.Id)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(realEstateKey, realEstateAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
