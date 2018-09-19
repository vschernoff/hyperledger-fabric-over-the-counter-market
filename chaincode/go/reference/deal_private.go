package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	"errors"
	"strings"
)

const (
	dealPrivateIndex = "DealPrivate"
)

type DealPrivateValue struct {
	Borrower  string  	  `json:"borrower"`
	Lender    string  	  `json:"lender"`
}

type DealPrivate struct {
	Key       DealKey   `json:"key"`
	Value     DealPrivateValue `json:"value"`
}


func (deal *DealPrivate) UpdateOrInsertIn(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	value, err := deal.ToLedgerValue()
	if err != nil {
		return err
	}

	collectionName, err := deal.GetDestinationCollectionName(stub)
	if err != nil {
		return err
	}

	// === Save members to state ===
	logger.Debug("Collection name: " +collectionName)

	if err = stub.PutPrivateData(collectionName, compositeKey, value); err != nil {
		return err
	}

	return nil
}

func (deal *DealPrivate) ExistsIn(stub shim.ChaincodeStubInterface, collectionName string) bool {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return false
	}

	if data, err := stub.GetPrivateData(collectionName, compositeKey); err != nil || data == nil {
		return false
	}

	return true
}

func (deal *DealPrivate) LoadFrom(stub shim.ChaincodeStubInterface, collectionName string) error {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	data, err := stub.GetPrivateData(collectionName, compositeKey)
	if err != nil {
		return err
	}

	return deal.FillFromLedgerValue(data)
}

func (deal *DealPrivate) FillFromLedgerValue(ledgerValue []byte) error {
	if err := json.Unmarshal(ledgerValue, &deal.Value); err != nil {
		return err
	} else {
		return nil
	}
}

func (deal *DealPrivate) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		deal.Key.ID,
	}

	return stub.CreateCompositeKey(dealPrivateIndex, compositeKeyParts)
}

func (deal *DealPrivate) ToLedgerValue() ([]byte, error) {
	return json.Marshal(deal.Value)
}

func (deal *DealPrivate) GetDestinationCollectionName(stub shim.ChaincodeStubInterface) (string, error) {
	Col := Collections{}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return "", errors.New(message)
	}
	Col.OrganizationName = creator

	if !Col.ExistsIn(stub) {
		message := fmt.Sprintf("available collections for organization %s not found", Col.OrganizationName)
		logger.Error(message)
		return "", errors.New(message)
	}

	if err := Col.LoadFrom(stub); err != nil {
		message := fmt.Sprintf("cannot load existing collections: %s", err.Error())
		logger.Error(message)
		return "", errors.New(message)
	}

	for _, Collection := range Col.AvailableCollections{
		//0    1    2
		//Org1 Org2 Deal
		SplitArray := strings.Split(string(Collection),"-")

		if deal.Value.Borrower == SplitArray[0]{
			if deal.Value.Lender == SplitArray[1]{
				return string(Collection), nil
			}
		}

		if deal.Value.Lender == SplitArray[0]{
			if deal.Value.Lender == SplitArray[1]{
				return string(Collection), nil
			}
		}

	}

	message := fmt.Sprintf("cannot find name of collections for organizations %s and %s", deal.Value.Lender, deal.Value.Borrower)
	logger.Error(message)
	return "", errors.New(message)
}




