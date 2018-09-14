package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

const (
	dealPrivateIndex = "DealPrivate"
)

type lenderT   string;
type borrowerT string;

var collectionDealPrivate = map [lenderT]map[borrowerT]string{
	"a": {
		"b" : "abDeals",
		"c" : "acDeals",
	},
	"b": {
		"a" : "abDeals",
		"c" : "bcDeals",
	},
	"c": {
		"a" : "acDeals",
		"b" : "bcDeals",
	},
}

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
	logger.Debug("Collection name: " + string(collectionDealPrivate[lenderT(deal.Value.Lender)][borrowerT(members.Value.Borrower)]))

	if err = stub.PutPrivateData(string(collectionDealPrivate[lenderT(members.Value.Lender)][borrowerT(members.Value.Borrower)]), compositeKey, value); err != nil {
		return err
	}

	return nil
}

func (members *DealPrivate) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		members.Key.ID,
	}

	return stub.CreateCompositeKey(dealPrivateIndex, compositeKeyParts)
}

func (members *DealPrivate) ToLedgerValue() ([]byte, error) {
	return json.Marshal(members.Value)
}