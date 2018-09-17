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

type Collection struct{
	Name    string `json:"name"`
	Policy  string `json:"policy"`
}

var CollectionConfig = "[ { \"name\": \"a-b-Deals\", \"policy\": \"OR('aMSP.member', 'bMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"a-c-Deals\", \"policy\": \"OR('aMSP.member','cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"b-c-Deals\", \"policy\": \"OR('bMSP.member', 'cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 } ]"

type DealPrivateValue struct {
	Borrower  string  	  `json:"borrower"`
	Lender    string  	  `json:"lender"`
}

type DealPrivate struct {
	Key       DealKey   `json:"key"`
	Value     DealPrivateValue `json:"value"`
}

func (deal *DealPrivate) GetCollectionsFromCreator(ledgerValue []byte) error {

	Collections := []Collection{}

	if err := json.Unmarshal(ledgerValue, &Collections); err != nil {
		return err
	} else {
		//Parsing Collections

		return nil
	}
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