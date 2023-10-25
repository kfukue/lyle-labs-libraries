package gethlyleaddresses

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofrs/uuid"
	"github.com/kfukue/lyle-labs-libraries/asset"
	"github.com/kfukue/lyle-labs-libraries/utils"
)

type GethAddress struct {
	ID            *int      `json:"id"`
	UUID          string    `json:"uuid"`
	Name          string    `json:"name"`
	AlternateName string    `json:"alternateName"`
	Description   string    `json:"description"`
	AddressStr    string    `json:"addressStr"`
	AddressTypeID *int      `json:"addressTypeId"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedBy     string    `json:"updatedBy"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func CreateOrGetContractAddressFromAsset(asset *asset.Asset) (*GethAddress, error) {
	contractAddress, err := GetGethAddressByAddressStr(asset.ContractAddress)
	if err != nil {
		log.Printf("Failed GetGethAddressByAddressStr: %v\n", err.Error())
		log.Fatal(err)
		return nil, err
	}
	// add as new address (contract) if doesn't exists
	if contractAddress == nil {
		contractName := fmt.Sprintf("Contract : %s", asset.Name)
		contractTypeID := utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_TYPE_ID
		newContractAddress := GethAddress{
			Name:          contractName,
			AlternateName: contractName,
			AddressStr:    asset.ContractAddress,
			AddressTypeID: &contractTypeID,
			CreatedBy:     utils.SYSTEM_NAME,
		}
		contractAddressId, err := InsertGethAddress(newContractAddress)
		if err != nil {
			log.Printf("Failed InsertGethAddress: %v\n", err.Error())
			log.Fatal(err)
			return nil, err
		}
		newContractAddress.ID = &contractAddressId
		contractAddress = &newContractAddress
	}
	return contractAddress, nil
}

func CreateOrGetEOAAddress(addressStr string) (*GethAddress, error) {
	eoaAddress, err := GetGethAddressByAddressStr(addressStr)
	if err != nil {
		log.Printf("Failed GetGethAddressByAddressStr: %v\n", err.Error())
		log.Fatal(err)
		return nil, err
	}
	// add as new address (contract) if doesn't exists
	if eoaAddress == nil {
		contractName := fmt.Sprintf("EOA: %s", addressStr)
		contractTypeID := utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID
		newContractAddress := GethAddress{
			Name:          contractName,
			AlternateName: contractName,
			AddressStr:    addressStr,
			AddressTypeID: &contractTypeID,
			CreatedBy:     utils.SYSTEM_NAME,
		}
		contractAddressId, err := InsertGethAddress(newContractAddress)
		if err != nil {
			log.Printf("Failed InsertGethAddress: %v\n", err.Error())
			log.Fatal(err)
			return nil, err
		}
		newContractAddress.ID = &contractAddressId
		eoaAddress = &newContractAddress
	}
	return eoaAddress, nil
}

func CreateOrGetContractAddress(addressStr string) (*GethAddress, error) {
	contractAddress, err := GetGethAddressByAddressStr(addressStr)
	if err != nil {
		log.Printf("Failed GetGethAddressByAddressStr: %v\n", err.Error())
		log.Fatal(err)
		return nil, err
	}
	// add as new address (contract) if doesn't exists
	if contractAddress == nil {
		contractName := fmt.Sprintf("Contract: %s", addressStr)
		contractTypeID := utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_TYPE_ID
		newContractAddress := GethAddress{
			Name:          contractName,
			AlternateName: contractName,
			AddressStr:    addressStr,
			AddressTypeID: &contractTypeID,
			CreatedBy:     utils.SYSTEM_NAME,
		}
		contractAddressId, err := InsertGethAddress(newContractAddress)
		if err != nil {
			log.Printf("Failed InsertGethAddress: %v\n", err.Error())
			log.Fatal(err)
			return nil, err
		}
		newContractAddress.ID = &contractAddressId
		contractAddress = &newContractAddress
	}
	return contractAddress, nil
}

func CreateContractAddress(addressStr string, isEOA bool) *GethAddress {
	var addressName string
	var contractTypeID int
	gethAddressUUID, err := uuid.NewV4()
	if err != nil {
		msg := fmt.Sprintf("error: uuid.NewV4(), during CreateContractAddress : %s", err.Error)
		log.Fatal(msg)
		return nil
	}
	if isEOA {
		addressName = fmt.Sprintf("EOA: %s", addressStr)
		contractTypeID = utils.EOA_ADDRESS_TYPE_STRUCTURED_VALUE_ID
	} else {
		addressName = fmt.Sprintf("Contract: %s", addressStr)
		contractTypeID = utils.CONTRACT_ADDRESS_TYPE_STRUCTURED_VALUE_TYPE_ID
	}

	gethAddress := GethAddress{
		Name:          addressName,
		UUID:          gethAddressUUID.String(),
		AlternateName: addressName,
		AddressStr:    addressStr,
		AddressTypeID: &contractTypeID,
		CreatedBy:     utils.SYSTEM_NAME,
	}
	return &gethAddress
}

func CreateEOAOrContractAddress(addressStr string, cl *ethclient.Client) (*GethAddress, error) {
	address := common.HexToAddress(addressStr)
	codeAtResult, err := cl.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Printf("Failed CodeAt: %v\n", err.Error())
		log.Fatal(err)
		return nil, err
	}
	//if result len is 0 EOA otherwise contract
	var gethAddress *GethAddress
	var isEOA bool
	// EOA
	if len(codeAtResult) == 0 {
		isEOA = true
	} else {
		isEOA = false
	}
	gethAddress = CreateContractAddress(addressStr, isEOA)
	return gethAddress, nil
}
