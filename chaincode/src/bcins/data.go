/**
 *  Nguyễn Tiến Thành - 20156455
 *	14/3/2019
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
	"time"
)

type UserRole int8

const (
	// Người sử dụng - người mua bán chứng chỉ quỹ === cổ đông
	ShareHolder UserRole = iota

	// Đơn vị chứng chỉ quỹ - gười tạo chứng chỉ quỹ
	Publisher

	// Chủ đầu tư - người tạo bất động sản
	Investor
)

//Cấu trúc đại diện người trong mạng
// Key consists of prefix + username
type shareholder struct {
	Username                          string              `json:"username"`                               //Tên đăng nhập
	TritList                          map[string]int      `json:"trit_list"`                              //Danh sach cac ma chung chi quy user nam giữ
	IdentityCard                      string              `json:"identity_card"`                          //Số CMT
	PublishedTrits                    []string            `json:"published_trits"`                        //Những mã TRIT đã được publish - chỉ dành cho
	FirstName                         string              `json:"first_name"`                             //Tên
	LastName                          string              `json:"last_name"`                              //Họ
	SellTritAdvertisingContractIdList map[string][]string `json:"sell_trit_advertising_contract_id_list"` //danh sách quảng cáo ccq đã đăng
	Balance                           int                 `json:"balance"`                                //số coin của user
}

type publisher struct {
	Username       string   `json:"username"`        //Tên đăng nhập
	IdentityCard   string   `json:"identity_card"`   //Số CMT
	PublishedTrits []string `json:"published_trits"` //Những mã TRIT đã publish
	FirstName      string   `json:"first_name"`      //Tên
	LastName       string   `json:"last_name"`       //Họ

}

//Cấu trúc đại diện cho hợp đồng user muốn bán chứng chỉ quỹ với số lượng và giá cụ thể
type sellTritAdvertisingContract struct {
	UUID   string    `json:"uuid"`
	Seller string    `json:"seller"`  // người tao
	TritId string    `json:"trit_id"` // id chứng chỉ quỹ
	Amount int       `json:"amount"`  // tổng số lượng chỉnh chỉ quỹ trao đổi
	Price  int       `json:"price"`   // giá bán
	Time   time.Time `json:"time"`    // thời điểm mua

}

//Cấu trúc đại diện cho hợp động tạo chứng chỉ quỹ
type publishContract struct {
	UUID      string    `json:"uuid"`
	Publisher string    `json:"publisher"` // người tao
	TritId    string    `json:"trit_id"`   // id chứng chỉ quỹ
	Amount    int       `json:"amount"`    // tổng số lượng chỉnh chỉ quỹ trao đổi
	Time      time.Time `json:"time"`      // thời điểm mua

}

//Cấu trúc đại diện cho hợp đồng thay đổi giá bất động sản
type changePriceContract struct {
	UUID      string    `json:"uuid"`
	Publisher string    `json:"publisher"` // người tao
	TritId    string    `json:"trit_id"`   // id chứng chỉ quỹ
	Price     int       `json:"price"`     // gia cua toan bo bds
	Time      time.Time `json:"time"`      // thời điểm mua
}

//Cấu trúc đại diện hợp đồng mua bán chưngs chỉ quỹ trong mạng
type transferContract struct {
	UUID      string    `json:"uuid"`
	Buyer     string    `json:"buyer"`      // người mua
	Seller    string    `json:"seller"`     // người bán
	TritId    string    `json:"trit_id"`    // id chứng chỉ quỹ
	Amount    int       `json:"amount"`     // tổng số lượng chỉnh chỉ quỹ trao đổi
	Time      time.Time `json:"time"`       // thời điểm mua
	TritPrice int       `json:"trit_price"` // giá của 1 chứng chỉ quỹ tại thời điểm trao đổi
}

//Cấu trúc đại diện bất động sản trong mạng
type realEstate struct {
	Id           string         `json:"id"`               // mã của bất động sản - ví dụ CHA, CHB
	Shareholders map[string]int `json:"shareholder_list"` /* danh sách các cố đông, id của cổ đông và số chứng chỉ quỹ người đó nắm giữ */
	Price        int            `json:"price"`            // giá trị
	SquareMeter  float64        `json:"square_meter"`     // diện tích
	Address      string         `json:"address"`          // Địa chỉ
	Description  string         `json:"description"`      // mô tả
	OwnerId      string         `json:"owner_id"`         //id của chủ sở hữu bất động sản
}

func (estate realEstate) totalTrit() int {
	shareHolders := estate.Shareholders
	total := 0
	for _, value := range shareHolders {
		total = total + value
	}
	return total
}

func (contract sellTritAdvertisingContract) changeAmount(stub shim.ChaincodeStubInterface, amount int, seller *shareholder) error {
	if amount <= 0 {
		return errors.New(fmt.Sprintf("func changeAmount amount must be greater than zero %d", amount))
	}
	key, err := getSellTritAdvertisingContractKey(stub, contract.UUID)
	if err != nil {
		return err
	}
	//change amount of advertising contract
	contract.Amount -= amount
	if contract.Amount == 0 {
		err = stub.DelState(key)
		if err != nil {
			return err
		}
		seller.deleteIdToShareholderSellTritAdvertisingContractList(contract.TritId, contract.UUID)

	} else {
		contractAsBytes, err := json.Marshal(contract)
		if err != nil {
			return err
		}

		err = stub.PutState(key, contractAsBytes)
		if err != nil {
			return err
		}
	}

	return nil

}

func (contract sellTritAdvertisingContract) changePrice(stub shim.ChaincodeStubInterface, price int) error {
	contract.Price = price
	advertisingContractAsBytes, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	advertisingContractKey, err := getSellTritAdvertisingContractKey(stub, contract.UUID)

	err = stub.PutState(advertisingContractKey, advertisingContractAsBytes)
	if err != nil {
		return err
	}

	return nil

}

// ghi đè lại hàm đọc json
func (s *UserRole) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	switch strings.ToUpper(value) {
	default:
		fallthrough
	case "S":
		*s = ShareHolder
	case "P":
		*s = Publisher
	case "I":
		*s = Investor
	}

	return nil
}

//ghi đè lại hàm ghi json
func (s UserRole) MarshalJSON() ([]byte, error) {
	var value string

	switch s {
	default:
		fallthrough
	case ShareHolder:
		value = "S"
	case Publisher:
		value = "P"
	case Investor:
		value = "I"
	}

	return json.Marshal(value)
}
