package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
)

const (
	dealIndex = "Deal"
)

const (
	dealKeyFieldsNumber      = 1
)

type DealKey struct {
	ID string `json:"id"`
}

type DealPublicValue struct {
	Amount    float32 `json:"amount"`
	Rate      float32 `json:"rate"`
	Timestamp int64   `json:"timestamp"`
}

type DealPublic struct {
	Key   DealKey   `json:"key"`
	Value DealPublicValue `json:"value"`
}

type DealValue struct {
	Amount    float32      `json:"amount"`
	Rate      float32      `json:"rate"`
	Timestamp int64   	   `json:"timestamp"`
	Borrower  string  	   `json:"borrower"`
	Lender    string  	   `json:"lender"`
}

type Deal struct {
	Key   DealKey       `json:"key"`
	Value DealValue     `json:"value"`
}

func (deal *DealPublic) FillFromCompositeKeyParts(compositeKeyParts []string) error {
	if len(compositeKeyParts) < dealKeyFieldsNumber {
		return errors.New(fmt.Sprintf("composite key parts array must contain at least %d items", dealKeyFieldsNumber))
	}

	if id, err := uuid.FromString(compositeKeyParts[0]); err != nil {
		return errors.New(fmt.Sprintf("unable to parse an ID from \"%s\"", compositeKeyParts[0]))
	} else if id.Version() != uuid.V4 {
		return errors.New("wrong ID format; expected UUID version 4")
	}

	deal.Key.ID = compositeKeyParts[0]

	return nil
}

func (deal *DealPublic) FillFromLedgerValue(ledgerValue []byte) error {
	if err := json.Unmarshal(ledgerValue, &deal.Value); err != nil {
		return err
	} else {
		return nil
	}
}

func (deal *DealPublic) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		deal.Key.ID,
	}

	return stub.CreateCompositeKey(dealIndex, compositeKeyParts)
}

func (deal *DealPublic) ToLedgerValue() ([]byte, error) {
	return json.Marshal(deal.Value)
}

func (deal *DealPublic) ExistsIn(stub shim.ChaincodeStubInterface) bool {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return false
	}

	if data, err := stub.GetState(compositeKey); err != nil || data == nil {
		return false
	}

	return true
}

func (deal *DealPublic) LoadFrom(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	data, err := stub.GetState(compositeKey)
	if err != nil {
		return err
	}

	return deal.FillFromLedgerValue(data)
}

func (deal *DealPublic) UpdateOrInsertIn(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	value, err := deal.ToLedgerValue()
	if err != nil {
		return err
	}

	// === Save deal to state ===
	if err = stub.PutState(compositeKey, value); err != nil {
		return err
	}

	return nil
}


