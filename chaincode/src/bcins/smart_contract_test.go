package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
	"time"
)

var testLog = shim.NewLogger("<chaincodeName>_test")

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		testLog.Info("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		testLog.Info("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	} else {
		testLog.Info("State value", name, "is", string(bytes), "as expected")
	}
}

func checkNoState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes != nil {
		testLog.Info("State", name, "should be absent; found value")
		t.FailNow()
	} else {
		testLog.Info("State", name, "is absent as it should be")
	}
}

func checkQueryOneArg(t *testing.T, stub *shim.MockStub, function string, argument string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(argument)})
	if res.Status != shim.OK {
		testLog.Info("Query", function, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		testLog.Info("Query", function, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		testLog.Info("Query value", function, "was", payload, "and not", value, "as expected")
		t.FailNow()
	} else {
		testLog.Info("Query value", function, "is", payload, "as expected")
	}
}

func checkQueryTwoArgs(t *testing.T, stub *shim.MockStub, function string, argument1 string, argument2 string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(argument1), []byte(argument2)})
	if res.Status != shim.OK {
		testLog.Info("Query", function, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		testLog.Info("Query", function, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		testLog.Info("Query value", function, "was", payload, "and not", value, "as expected")
		t.FailNow()
	} else {
		testLog.Info("Query value", function, "is", payload, "as expected")
	}
}

func createDefaultUsers(t *testing.T, mockStub *shim.MockStub) {
	var functionAndArgs []string
	functionName := "create_publisher"

	dto0 := struct {
		Username string `json:"username"`
	}{}

	dto0.Username = "C"

	arg, _ := json.Marshal(dto0)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)

	functionAndArgs = []string{}
	functionName = "create_shareholder"

	dto1 := struct {
		Username string `json:"username"`
	}{}

	dto1.Username = "A"

	arg, _ = json.Marshal(dto1)

	args = []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)

	functionAndArgs = []string{}
	functionName = "create_shareholder"
	dto2 := struct {
		Username string `json:"username"`
	}{}

	dto2.Username = "B"

	arg, _ = json.Marshal(dto2)

	args = []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)
}

func createDefaultRealEstate(t *testing.T, stub *shim.MockStub) {

	createDefaultUsers(t, stub)

	functionAndArgs := []string{}
	functionName := "create_and_publish_real_estate"

	dto := struct {
		UUID        string    `json:"uuid"`
		Id          string    `json:"id"`           // mã của bất động sản - ví dụ CHA, CHB
		Price       int       `json:"price"`        // giá trị
		SquareMeter float64   `json:"square_meter"` // diện tích
		Address     string    `json:"address"`      // Địa chỉ
		OwnerId     string    `json:"owner_id"`     //id của chủ sở hữu bất động sản
		Amount      int       `json:"amount"`       // số lượng chứng chỉ quỹ được bán ra
		PublisherId string    `json:"publisher_id"` // nguoi tao
		CreateTime  time.Time `json:"create_time"`  // thoi diem tao

	}{}

	dto.UUID = "1"
	dto.Id = "1"
	dto.Price = 10000
	dto.SquareMeter = 1000
	dto.Address = "VN"
	dto.OwnerId = "A"
	dto.Amount = 2000
	dto.PublisherId = "C"
	dto.CreateTime = time.Now()
	arg, _ := json.Marshal(dto)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, stub, functionAndArgs)
}

func checkBadQuery(t *testing.T, stub *shim.MockStub, function string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status == shim.OK {
		testLog.Info("Query", function, "unexpectedly succeeded")
		t.FailNow()
	} else {
		testLog.Info("Query", function, "failed as espected, with message: ", res.Message)

	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := [][]byte{}

	l := len(functionAndArgs)
	for i := 0; i < l; i++ {
		v := []byte(functionAndArgs[i])
		functionAndArgsAsBytes = append(functionAndArgsAsBytes, v)
	}

	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status != shim.OK {
		testLog.Info("Invoke", functionAndArgs, "failed", string(res.Message))
		t.FailNow()
	} else {
		testLog.Info("Invoke", functionAndArgs, "successful", string(res.Message))
	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {

	functionAndArgsAsBytes := [][]byte{}

	l := len(functionAndArgs)
	for i := 0; i < l; i++ {
		v := []byte(functionAndArgs[i])
		functionAndArgsAsBytes = append(functionAndArgsAsBytes, v)
	}
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status == shim.OK {
		testLog.Info("Invoke", functionAndArgs, "unexpectedly succeeded")
		t.FailNow()
	} else {
		testLog.Info("Invoke", functionAndArgs, "failed as espected, with message: "+res.Message)
	}
}

////==================================================================
//// TestSellTritAd - Test the 'create_sell_trit_advertising_contract' function
//// =================================================================
//func TestCreateSellTritAdvertisingContract(t *testing.T) {
//	simpleChaincode := new(SmartContract)
//	mockStub := shim.NewMockStub("Test Feature Creation", simpleChaincode)
//
//	var functionAndArgs []string
//	functionName := "create_sell_trit_advertising_contract"
//
//	dto := struct {
//		UUID			 string		 `json:"uuid"`
//		Seller           string      `json:"seller"`					    // người tao
//		TritId			 string      `json:"trit_id"`						// id chứng chỉ quỹ
//		Amount           int	     `json:"amount"`						// tổng số lượng chỉnh chỉ quỹ trao đổi
//		Price            float64 	 `json:"price"`							// giá bán
//		Time        	 time.Time   `json:"time"`							// thời điểm mua
//	}{}
//
//	dto.UUID = "1"
//	dto.Seller = "A"
//	dto.TritId = "1"
//	dto.Amount = 20
//	dto.Price = 10
//	dto.Time = time.Now();
//	arg, _ := json.Marshal(dto);
//
//	args := []string{string(arg)}
//	functionAndArgs = append(functionAndArgs, functionName)
//	functionAndArgs = append(functionAndArgs, args...)
//
//	checkInvoke(t, mockStub, functionAndArgs)
//
//
//}

//==================================================================
// TestCreateShareHolder - Test the 'create_sell_trit_advertising_contract' function
// =================================================================
func TestCreateShareholderA(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Create Shareholder", simpleChaincode)

	var functionAndArgs []string
	functionName := "create_shareholder"

	dto := struct {
		Username string `json:"username"`
	}{}

	dto.Username = "A"

	arg, _ := json.Marshal(dto)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)

}

func TestCreateShareholderB(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Create Shareholder", simpleChaincode)

	var functionAndArgs []string
	functionName := "create_shareholder"

	dto := struct {
		Username string `json:"username"`
	}{}

	dto.Username = "B"

	arg, _ := json.Marshal(dto)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)

}

func TestCreateShareholderC(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Create Shareholder", simpleChaincode)

	var functionAndArgs []string
	functionName := "create_shareholder"

	dto := struct {
		Username string `json:"username"`
	}{}

	dto.Username = "C"

	arg, _ := json.Marshal(dto)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)

}

func TestCreatePublisher(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Create Publisher", simpleChaincode)

	var functionAndArgs []string
	functionName := "create_publisher"

	dto := struct {
		Username string `json:"username"`
	}{}

	dto.Username = "C"

	arg, _ := json.Marshal(dto)

	args := []string{string(arg)}
	functionAndArgs = append(functionAndArgs, functionName)
	functionAndArgs = append(functionAndArgs, args...)

	checkInvoke(t, mockStub, functionAndArgs)
}

func TestCreateRealEstateNPublishTrit(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Create Real Estate And Publish Trit", simpleChaincode)
	createDefaultRealEstate(t, mockStub)
}

func TestTransferTrit(t *testing.T) {
	simpleChaincode := new(SmartContract)
	mockStub := shim.NewMockStub("Test Transfer Trit", simpleChaincode)
	createDefaultRealEstate(t, mockStub)

}
