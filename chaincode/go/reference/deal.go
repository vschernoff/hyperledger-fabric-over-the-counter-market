package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

const (
	dealIndex = "Deal"
	orgsIndex = "Orgs"
)

const (
	dealKeyFieldsNumber      = 1
	dealBasicArgumentsNumber = 3
)

type DealKey struct {
	ID string `json:"id"`
}

type DealValue struct {
	Borrower  string  `json:"borrower"`
	Lender    string  `json:"lender"`
	Amount    float32 `json:"amount"`
	Rate      float32 `json:"rate"`
	Timestamp int64   `json:"timestamp"`
}

type Deal struct {
	Key   DealKey   `json:"key"`
	Value DealValue `json:"value"`
}

func (deal *Deal) FillFromArguments(args []string) error {
	//if len(args) < dealBasicArgumentsNumber {
	//	return errors.New(fmt.Sprintf("arguments array must contain at least %d items", dealBasicArgumentsNumber))
	//}
	//
	//dealType, err := strconv.Atoi(args[0])
	//if err != nil {
	//	return errors.New(fmt.Sprintf("unable to parse deal type: %s", err.Error()))
	//}
	//if dealType != typeLend && dealType != typeBorrow {
	//	return errors.New(fmt.Sprintf("unsupported type of a deal: %d", dealType))
	//}
	//deal.Value.Type = dealType
	//
	//amount, err := strconv.ParseFloat(args[1], 32)
	//if err != nil {
	//	return errors.New(fmt.Sprintf("unable to parse the amount: %s", err.Error()))
	//}
	//if amount < 0 {
	//	return errors.New("amount must be larger than zero")
	//}
	//deal.Value.Amount = float32(amount)
	//
	//rate, err := strconv.ParseFloat(args[2], 32)
	//if err != nil {
	//	return errors.New(fmt.Sprintf("unable to parse the rate: %s", err.Error()))
	//}
	//if rate < 0 {
	//	return errors.New("rate must be larger than zero")
	//}
	//deal.Value.Rate = float32(rate)

	return nil
}

func (deal *Deal) FillFromCompositeKeyParts(compositeKeyParts []string) error {
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

func (deal *Deal) FillFromLedgerValue(ledgerValue []byte) error {
	if err := json.Unmarshal(ledgerValue, &deal.Value); err != nil {
		return err
	} else {
		return nil
	}
}

func (deal *Deal) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		deal.Key.ID,
	}

	return stub.CreateCompositeKey(dealIndex, compositeKeyParts)
}

func (deal *Deal) ToLedgerValue() ([]byte, error) {
	return json.Marshal(deal.Value)
}

func (deal *Deal) ExistsIn(stub shim.ChaincodeStubInterface) bool {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return false
	}

	if data, err := stub.GetState(compositeKey); err != nil || data == nil {
		return false
	}

	return true
}

func (deal *Deal) LoadFrom(stub shim.ChaincodeStubInterface) error {
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

func (deal *Deal) UpdateOrInsertIn(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := deal.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	value, err := deal.ToLedgerValue()
	if err != nil {
		return err
	}

	var collection string
	if strings.Compare(deal.Value.Borrower, deal.Value.Lender) < 0 {
		collection = fmt.Sprintf("%s-%s", deal.Value.Borrower, deal.Value.Lender)
	} else {
		collection = fmt.Sprintf("%s-%s", deal.Value.Lender, deal.Value.Borrower)
	}
	fmt.Sprintf("Collection: %s, deal key: %s", collection, compositeKey)

	//if err = stub.PutPrivateData(collection, compositeKey, value); err != nil {
	if err = stub.PutState(compositeKey, value); err != nil {
		return err
	}

	if err = stub.PutPrivateData(collection, compositeKey, value); err != nil {
		return err
	}

	return nil
}