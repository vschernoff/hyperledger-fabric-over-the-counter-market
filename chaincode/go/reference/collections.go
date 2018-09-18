package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	collectionsIndex = "Collections"
)

type Collections struct {
	OrganizationName     string      `json:"organizationName"`
	AvailableCollections []string    `json:"availableCollections"`
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
		col.OrganizationName,
	}

	return stub.CreateCompositeKey(collectionsIndex, compositeKeyParts)
}

func (col *Collections) ToLedgerValue() ([]byte, error) {
	return json.Marshal(col.AvailableCollections)
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
	if err := json.Unmarshal(ledgerValue, &col.AvailableCollections); err != nil {
		return err
	} else {
		return nil
	}
}