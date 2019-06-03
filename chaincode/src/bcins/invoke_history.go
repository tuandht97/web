package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

func getHistoryForUFO(APIStub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dto := struct {
		Id        string `json:"asset_id"`
		AssetType string `json:"asset_type"`
	}{}

	err := json.Unmarshal([]byte(args[0]), &dto)
	if err != nil {
		return shim.Error(err.Error())
	}

	ufoId := dto.Id

	var key string

	switch dto.AssetType {
	case "shareholder":
		key, err = getShareholderKey(APIStub, ufoId)
	case "publisher":
		key, err = getPublisherKey(APIStub, ufoId)
	case "real_estate":
		key, err = getRealEstateKey(APIStub, ufoId)
	case "publish_contract":
		key, err = getPublishContractKey(APIStub, ufoId)
	case "transfer_contract":
		key, err = getTransferContractKey(APIStub, ufoId)
	case "change_price_contract":
		key, err = getChangePriceContractKey(APIStub, ufoId)
	case "sell_trit_advertising_contract":
		key, err = getSellTritAdvertisingContractKey(APIStub, ufoId)
	default:
		return shim.Error("Don't know type" + dto.AssetType)

	}

	if err != nil {
		return shim.Error(err.Error())
	}

	resultsIterator, err := APIStub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)

		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- History returning:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}
