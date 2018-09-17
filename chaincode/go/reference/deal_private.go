package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
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

	// === Save members to state ===
	logger.Debug("Collection name: " + string(collectionDealPrivate[lenderT(deal.Value.Lender)][borrowerT(deal.Value.Borrower)]))

	if err = stub.PutPrivateData(string(collectionDealPrivate[lenderT(deal.Value.Lender)][borrowerT(deal.Value.Borrower)]), compositeKey, value); err != nil {
		return err
	}

	return nil
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