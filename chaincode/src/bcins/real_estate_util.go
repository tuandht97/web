package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (estate realEstate) publishTrit(stub shim.ChaincodeStubInterface, amount int, owner *shareholder) error {

	if amount <= 0 {
		return errors.New(fmt.Sprintf("Amount must be greater than zero: %d", amount))
	}

	if len(estate.Shareholders) == 0 {
		estate.Shareholders = make(map[string]int)
	}

	estate.Shareholders[estate.OwnerId] = amount

	//set trit id for owner

	owner.setTritIdToShareholderPublishedTrits(estate.Id, amount)

	realEstateAsBytes, err := json.Marshal(estate)
	if err != nil {
		return err
	}
	realEstateKey, err := getRealEstateKey(stub, estate.Id)
	err = stub.PutState(realEstateKey, realEstateAsBytes)
	if err != nil {
		return err
	}

	return nil

}

func (estate realEstate) changePrice(stub shim.ChaincodeStubInterface, price int) error {

	estate.Price = price
	realEstateAsBytes, err := json.Marshal(estate)
	if err != nil {
		return err
	}
	realEstateKey, err := getRealEstateKey(stub, estate.Id)
	err = stub.PutState(realEstateKey, realEstateAsBytes)
	if err != nil {
		return err
	}
	return nil

}

func (estate realEstate) transferTrit(stub shim.ChaincodeStubInterface, buyer *shareholder, seller *shareholder, amount int) error {
	shareholders := estate.Shareholders

	if shareholders == nil {
		return errors.New("real estate has not yet published trit " + estate.Id)
	}

	if tritOfseller, ok := shareholders[seller.Username]; ok {

		if tritOfseller < amount {
			return errors.New("Number trit of seller must greater than number trit is transferred " + estate.Id)
		}

		tritOfseller -= amount

		tritOfbuyer, ok := shareholders[buyer.Username]

		if ok {
			tritOfbuyer += amount
		} else {
			tritOfbuyer = amount
		}

		if tritOfseller == 0 {
			delete(shareholders, seller.Username)
			seller.deleteTritIdFromShareholderTritList(estate.Id)

		} else {
			shareholders[seller.Username] = tritOfseller
			seller.setTritIdToShareholderTritList(estate.Id, tritOfseller)

		}

		shareholders[buyer.Username] = tritOfbuyer
		buyer.setTritIdToShareholderTritList(estate.Id, tritOfbuyer)

		totalTrit := estate.totalTrit()
		//if buyer own all trit => change owner of real estate
		if tritOfbuyer == totalTrit {
			estate.OwnerId = buyer.Username
		}

		realEstateAsBytes, _ := json.Marshal(estate)
		realEstateKey, err := getRealEstateKey(stub, estate.Id)
		err = stub.PutState(realEstateKey, realEstateAsBytes)

		if err != nil {
			return errors.New("Failed to modified asset " + estate.Id)
		}

	} else {
		return errors.New("seller " + seller.Username + " doesn't have trit with id = " + estate.Id)
	}

	return nil
}

func getRealEstateState(stub shim.ChaincodeStubInterface, id string) (realEstate, error) {
	if len(id) == 0 {
		return realEstate{}, errors.New("func getRealEstateState id is empty ")
	}

	key, err := getRealEstateKey(stub, id)
	if err != nil {
		return realEstate{}, err
	}
	estateBytes, _ := stub.GetState(key)
	if len(estateBytes) == 0 {
		return realEstate{}, errors.New("Real Estate id does not exist " + id)
	}

	response := realEstate{}
	err = json.Unmarshal(estateBytes, &response)
	return response, err
}
