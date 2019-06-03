package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("main")

type SmartContract struct {
}

var bcFunctions = map[string]func(shim.ChaincodeStubInterface, []string) peer.Response{
	// Shareholder
	"create_shareholder": createShareholder,
	"get_shareholder":    getShareholder,
	"list_shareholder":   listShareholder,

	// Publisher
	"create_publisher": createPublisher,
	"get_publisher":    getPublisher,
	"list_publisher":   listPublisher,

	// Coin
	"pay_in": payIn,

	// Contract
	"create_sell_trit_advertising_contract":       createSellTritAdvertisingContract,
	"list_sell_trit_advertising_contract":         listSellTritAdvertisingContract,
	"change_price_sell_trit_advertising_contract": changePriceSellTritAdvertisingContract,
	"delete_sell_trit_advertising_contract":       deleteSellTritAdvertisingContract,
	"list_transfer_contract_by_buyer":             listTransferContractByBuyer,
	"list_transfer_contract_by_seller":            listTransferContractBySeller,
	"list_sell_trit_advertising_contract_by_user": listSellTritAdvertisingContractByUser,
	"create_change_price_contract":                createChangePriceContract,
	"list_real_estate":                            listRealEstate,
	"create_real_estate":                          createRealEstate,
	"create_and_publish_real_estate":              createAndPublishRealEstate,
	"list_transfer_contract":                      listTransferContract,
	"list_change_price_contract":                  listChangePriceContract,
	"list_publish_contract":                       listPublishContract,
	"create_publish_contract":                     createPublishContract,
	"create_transfer_contract_for_buyer":          createTransferContractForBuyer,
	"list_publish_contract_by_publisher":          listPublishContractByPublisher,
	"get_history_for_ufo":                         getHistoryForUFO,
	"pay_in_by_shareholder":                       payInByShareholder,
}

// Init callback representing the invocation of a chaincode
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke Function accept blockchain code invocations.
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "init" {
		return t.Init(stub)
	}
	bcFunc := bcFunctions[function]
	if bcFunc == nil {
		return shim.Error("Invalid invoke function. " + function)
	}
	return bcFunc(stub, args)
}

func main() {
	logger.SetLevel(shim.LogInfo)
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
