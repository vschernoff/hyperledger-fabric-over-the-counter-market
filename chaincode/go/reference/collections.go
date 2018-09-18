package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	collectionsIndex = "Collections"
)

type CollectionFromConfig struct{
	Name    string `json:"name"`
	Policy  string `json:"policy"`
}

type Collection struct {
	Name        string   `json:"name"`
}

type Collections struct {
	OrgName          string            `json:"org"`
	ListCollections  []Collection      `json:"collections"`
}

func (col *Collections) UpdateOrInsertIn(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := col.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	value, err := col.ToLedgerValue()
	if err != nil {
		return err
	}

	// === Save collections to state ===
	if err = stub.PutState(compositeKey, value); err != nil {
		return err
	}

	return nil
}

func (col *Collections) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		col.OrgName,
	}

	return stub.CreateCompositeKey(collectionsIndex, compositeKeyParts)
}

func (collections *Collections) ToLedgerValue() ([]byte, error) {
	return json.Marshal(collections.ListCollections)
}

func (col *Collections) ExistsIn(stub shim.ChaincodeStubInterface) bool {
	compositeKey, err := col.ToCompositeKey(stub)
	if err != nil {
		return false
	}

	if data, err := stub.GetState(compositeKey); err != nil || data == nil {
		return false
	}

	return true
}

func (col *Collections) LoadFrom(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := col.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	data, err := stub.GetState(compositeKey)
	if err != nil {
		return err
	}

	return col.FillFromLedgerValue(data)
}

func (col *Collections) FillFromLedgerValue(ledgerValue []byte) error {
	if err := json.Unmarshal(ledgerValue, &col.ListCollections); err != nil {
		return err
	} else {
		return nil
	}
}