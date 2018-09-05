package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

const (
	membersIndex = "Members"
)

const (
	abDeals = iota
	acDeals
	caDeals
)

type lenderT   string;
type borrowerT string;

var collectionMembers = map [lenderT]map[borrowerT]string{
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

type MembersKey struct {
	ID string `json:"id"`
}

type MembersValue struct {
	Borrower  string  	  `json:"borrower"`
	Lender    string  	  `json:"lender"`
}

type Members struct {
	Key       MembersKey   `json:"key"`
	Value     MembersValue `json:"value"`
}

func (members *Members) UpdateOrInsertIn(stub shim.ChaincodeStubInterface) error {
	compositeKey, err := members.ToCompositeKey(stub)
	if err != nil {
		return err
	}

	value, err := members.ToLedgerValue()
	if err != nil {
		return err
	}

	// === Save members to state ===
	if err = stub.PutPrivateData(collectionMembers[lenderT(members.Value.Lender)][borrowerT(members.Value.Borrower)], compositeKey, value); err != nil {
		return err
	}

	return nil
}

func (members *Members) ToCompositeKey(stub shim.ChaincodeStubInterface) (string, error) {
	compositeKeyParts := []string {
		members.Key.ID,
	}

	return stub.CreateCompositeKey(membersIndex, compositeKeyParts)
}

func (members *Members) ToLedgerValue() ([]byte, error) {
	return json.Marshal(members.Value)
}