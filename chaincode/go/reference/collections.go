package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	collectionsIndex = "Collections"
)

//type lenderT   string;
//type borrowerT string;

//var collectionDealPrivate = map [lenderT]map[borrowerT]string{
//	"a": {
//		"b" : "abDeals",
//		"c" : "acDeals",
//	},
//	"b": {
//		"a" : "abDeals",
//		"c" : "bcDeals",
//	},
//	"c": {
//		"a" : "acDeals",
//		"b" : "bcDeals",
//	},
//}
type DealKey struct {
	ID string `json:"id"`
}

type DealPrivateValue struct {
	Borrower  string  	  `json:"borrower"`
	Lender    string  	  `json:"lender"`
}

type Collection struct {
	Name        string   `json:"name"`
}

type Collections struct {
	Org         string       `json:"org"`
	Collections []Collection `json:"collections"`
}

type CollectionFromConfig struct{
	Name    string `json:"name"`
	Policy  string `json:"policy"`
}

var CollectionConfig = "[ { \"name\": \"a-b-Deals\", \"policy\": \"OR('aMSP.member', 'bMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"a-c-Deals\", \"policy\": \"OR('aMSP.member','cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"b-c-Deals\", \"policy\": \"OR('bMSP.member', 'cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 } ]"

func (collection *Collection) PutCollections(stub shim.ChaincodeStubInterface, CollectionConfig string) error {

	var Entries []string
	CollectionsFromConfig := []CollectionFromConfig{}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return err
	}

	logger.Debug("Creator: " + creator)


	if err := json.Unmarshal([]byte(CollectionConfig), &CollectionsFromConfig); err != nil {
		return err
	} else {
		for _, Col := range CollectionsFromConfig{
			ok := strings.Contains(string(Col.Policy), string(creator)+"MSP.member")
			if ok{
			Entries = append(Entries, Col.Name)
			}
		}
		if err = stub.PutState(compositeKey, json.Marshal(Capabilities)); err != nil {
			return err
		}
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