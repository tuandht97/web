package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (s shareholder) saveShareHolderState(stub shim.ChaincodeStubInterface) error {

	userKey, err := getShareholderKey(stub, s.Username)
	if err != nil {
		return err
	}
	userAsBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = stub.PutState(userKey, userAsBytes)
	return err
}

func (s *shareholder) setTritIdToShareholderTritList(tritId string, amount int) {

	if tritId == "" {
		return
	}

	if amount <= 0 {
		return
	}

	if len(s.TritList) == 0 {
		s.TritList = make(map[string]int)
	}

	s.TritList[tritId] = amount

}

func (s *shareholder) setTritIdToShareholderPublishedTrits(tritId string, amount int) {

	if tritId == "" {
		return
	}

	if amount <= 0 {
		return
	}

	if len(s.PublishedTrits) == 0 {
		s.PublishedTrits = []string{}
	}

	s.PublishedTrits = append(s.PublishedTrits, tritId)

	if len(s.TritList) == 0 {
		s.TritList = make(map[string]int)
	}

	s.TritList[tritId] = amount

}

func (s *shareholder) deleteTritIdFromShareholderTritList(tritId string) {
	if tritId == "" {
		return
	}

	delete(s.TritList, tritId)

}

func (s *shareholder) setIdToShareholderSellTritAdvertisingContractList(tritId string, contractId string) {

	if tritId == "" || contractId == "" {
		return
	}

	if len(s.SellTritAdvertisingContractIdList) == 0 {
		s.SellTritAdvertisingContractIdList = make(map[string][]string)
	}

	if _, ok := s.SellTritAdvertisingContractIdList[tritId]; ok {
		if len(s.SellTritAdvertisingContractIdList[tritId]) == 0 {
			s.SellTritAdvertisingContractIdList[tritId] = []string{}
		}
	} else {
		s.SellTritAdvertisingContractIdList[tritId] = []string{}
	}

	if !Contains(s.SellTritAdvertisingContractIdList[tritId], contractId) {
		s.SellTritAdvertisingContractIdList[tritId] = append(s.SellTritAdvertisingContractIdList[tritId], contractId)
	}

}

func (s *shareholder) deleteIdToShareholderSellTritAdvertisingContractList(tritId string, contractId string) {

	if tritId == "" || contractId == "" {
		return
	}

	if len(s.SellTritAdvertisingContractIdList) == 0 {
		s.SellTritAdvertisingContractIdList = make(map[string][]string)
	}

	if v, ok := s.SellTritAdvertisingContractIdList[tritId]; ok {
		v = Remove(v, contractId)
		s.SellTritAdvertisingContractIdList[tritId] = v
	}

}

func (s shareholder) checkEnoughCoin(require int) bool {
	if s.Balance < require {
		return false
	} else {
		return true
	}
}

func (s shareholder) getTotalAdvertisingTritOf(stub shim.ChaincodeStubInterface, tritId string) (int, error) {
	if list, ok := s.SellTritAdvertisingContractIdList[tritId]; ok {
		total := 0
		for _, id := range list {
			advertisingContractAsBytes, err := validateSellTritAdvertisingContractState(stub, id)
			if err != nil {
				return 0, err
			}

			adContract := sellTritAdvertisingContract{}
			err = json.Unmarshal(advertisingContractAsBytes, &adContract)
			if err != nil {
				return 0, err
			}

			total += adContract.Amount

		}
		return total, nil
	} else {
		return 0, nil
	}
}

func getShareholderState(stub shim.ChaincodeStubInterface, id string) (shareholder, error) {
	if len(id) == 0 {
		return shareholder{}, errors.New("func getShareholderState id is empty ")
	}

	userKey, err := getShareholderKey(stub, id)
	if err != nil {
		return shareholder{}, err
	}
	userBytes, _ := stub.GetState(userKey)
	if len(userBytes) == 0 {
		return shareholder{}, errors.New("User name does not exist " + id)
	}

	response := shareholder{}
	err = json.Unmarshal(userBytes, &response)
	return response, err
}
