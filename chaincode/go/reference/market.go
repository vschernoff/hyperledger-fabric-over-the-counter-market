
package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

var logger = shim.NewLogger("MarketChaincode")

// MarketChaincode example simple Chaincode implementation
type MarketChaincode struct {
}

func (t *MarketChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")
	var CollectionConfig = "[ { \"name\": \"a-b-Deals\", \"policy\": \"OR('aMSP.member', 'bMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"a-c-Deals\", \"policy\": \"OR('aMSP.member','cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 }, { \"name\": \"b-c-Deals\", \"policy\": \"OR('bMSP.member', 'cMSP.member')\", \"requiredPeerCount\": 0, \"maxPeerCount\": 3, \"blockToLive\":3 } ]"

	if err := InitCollections(stub, CollectionConfig); err != nil {
		message := fmt.Sprintf("cannot init collections: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	return shim.Success(nil)
}

func (t *MarketChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")

	function, args := stub.GetFunctionAndParameters()
	if function == "placeBid" {
		// place a bid on the market
		return t.placeBid(stub, args)
	} else if function == "editBid" {
		return t.editBid(stub, args)
	} else if function == "cancelBid" {
		return t.cancelBid(stub, args)
	} else if function == "makeDeal" {
		return t.makeDeal(stub, args)
	} else if function == "queryBids" {
		return t.queryBids(stub, args)
	} else if function == "queryBidsCreator" {
		return t.queryBidsCreator(stub, args)
	} else if function == "queryDeals" {
		return t.queryDeals(stub, args)
	//} else if function == "queryDealsCreatorByTime" {
	//	return t.queryDealsCreatorByTime(stub, args)
	}

	return pb.Response{Status:403, Message:"Invalid invoke function name."}
}

func (t *MarketChaincode) placeBid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.placeBid is running")
	logger.Debug("MarketChaincode.placeBid")

	if len(args) < bidBasicArgumentsNumber {
		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
			bidBasicArgumentsNumber, len(args))
		logger.Error(message)
		return shim.Error(message)
	}

	bid := Bid{}
	if err := bid.FillFromArguments(args); err != nil {
		message := fmt.Sprintf("cannot fill a bid from arguments: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Debug("Creator: " + creator)

	bid.Key.ID = uuid.Must(uuid.NewV4()).String()

	bid.Value.Creator = creator
	bid.Value.Status = statusActive
	bid.Value.Timestamp = time.Now().UTC().Unix()

	if bytes, err := json.Marshal(bid); err == nil {
		logger.Debug("Bid: " + string(bytes))
	}

	if err := bid.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	logger.Info("MarketChaincode.placeBid exited without errors")
	logger.Debug("Success: MarketChaincode.placeBid")
	return shim.Success(nil)
}

func (t *MarketChaincode) editBid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.editBid is running")
	logger.Debug("MarketChaincode.editBid")

	if len(args) < bidBasicArgumentsNumber + bidKeyFieldsNumber {
		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
			bidBasicArgumentsNumber + bidKeyFieldsNumber, len(args))
		logger.Error(message)
		return shim.Error(message)
	}

	var bid, bidToUpdate Bid
	if err := bid.FillFromCompositeKeyParts(args[:1]); err != nil {
        message := fmt.Sprintf("cannot fill a bid from arguments: %s", err.Error())
        logger.Error(message)
        return shim.Error(message)
    }
    bidToUpdate.Key.ID = bid.Key.ID
	if err := bid.FillFromArguments(args[1:]); err != nil {
		message := fmt.Sprintf("cannot fill a bid from arguments: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Debug("Bid to edit: " + bid.Key.ID)

	if !bidToUpdate.ExistsIn(stub) {
		message := fmt.Sprintf("bid with ID %s not found", bid.Key.ID)
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	if err := bidToUpdate.LoadFrom(stub); err != nil {
		message := fmt.Sprintf("cannot load existing bid: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	if bidToUpdate.Value.Creator != creator {
		message := "unable to edit a bid owned by not the creator"
		logger.Error(message)
		return pb.Response{Status: 403, Message: message}
	}

	if bidToUpdate.Value.Status != statusActive {
		message := fmt.Sprintf("unable to edit the bid: bid status is not Active")
		logger.Error(message)
		return shim.Error(message)
	}

	bidToUpdate.Value.Amount = bid.Value.Amount
	bidToUpdate.Value.Rate = bid.Value.Rate

    if bytes, err := json.Marshal(bidToUpdate); err == nil {
        logger.Debug("Bid: " + string(bytes))
    }

	if err := bidToUpdate.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	logger.Info("MarketChaincode.editBid exited without errors")
	logger.Debug("Success: MarketChaincode.editBid")
	return shim.Success(nil)
}

func (t *MarketChaincode) cancelBid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.cancelBid is running")
	logger.Debug("MarketChaincode.cancelBid")

	if len(args) < bidKeyFieldsNumber {
		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
			bidKeyFieldsNumber, len(args))
		logger.Error(message)
		return shim.Error(message)
	}

	bid := Bid{}
	if err := bid.FillFromCompositeKeyParts(args[:1]); err != nil {
		message := fmt.Sprintf("cannot fill a bid from arguments: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Debug("Bid to cancel: " + bid.Key.ID)

	if !bid.ExistsIn(stub) {
		message := fmt.Sprintf("bid with ID %s not found", bid.Key.ID)
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	if err := bid.LoadFrom(stub); err != nil {
		message := fmt.Sprintf("cannot load existing bid: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	if bid.Value.Creator != creator {
		message := fmt.Sprintf(
			"no privileges to cancel specified bid from the side of organization %s", creator)
		logger.Error(message)
		return pb.Response{Status: 403, Message: message}
	}

	if bid.Value.Status != statusActive {
		message := fmt.Sprintf("unable to cancel the bid: bid status is not Active")
		logger.Error(message)
		return shim.Error(message)
	}

	bid.Value.Status = statusCancelled

	if bytes, err := json.Marshal(bid); err == nil {
		logger.Debug("Bid: " + string(bytes))
	}

	if err := bid.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	logger.Info("MarketChaincode.cancelBid exited without errors")
	logger.Debug("Success: MarketChaincode.cancelBid")
	return shim.Success(nil)
}

func (t *MarketChaincode) makeDeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.makeDeal is running")
	logger.Debug("MarketChaincode.makeDeal")

	if len(args) < bidKeyFieldsNumber {
		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
			bidKeyFieldsNumber, len(args))
		logger.Error(message)
		return shim.Error(message)
	}

	bid := Bid{}
	if err := bid.FillFromCompositeKeyParts(args[:1]); err != nil {
		message := fmt.Sprintf("cannot fill a bid from arguments: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Debug("Bid to accept: " + bid.Key.ID)

	if !bid.ExistsIn(stub) {
		message := fmt.Sprintf("bid with ID %s not found", bid.Key.ID)
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	if err := bid.LoadFrom(stub); err != nil {
		message := fmt.Sprintf("cannot load existing bid: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 404, Message: message}
	}

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	if bid.Value.Creator == creator {
		message := "unable to make a deal on a bid owned by the creator"
		logger.Error(message)
		return pb.Response{Status: 403, Message: message}
	}

	if bid.Value.Status != statusActive {
		message := fmt.Sprintf("unable to accept the bid: bid status is not Active")
		logger.Error(message)
		return shim.Error(message)
	}

	bid.Value.Status = statusInactive

	var lender, borrower string
	if bid.Value.Type == typeLend {
		lender, borrower = bid.Value.Creator, creator
	} else { // typeBorrow
		lender, borrower = creator, bid.Value.Creator
	}
    ID := uuid.Must(uuid.NewV4()).String()
	dealPublic := DealPublic {
		Key: DealKey {
			ID: ID,
		},
		Value: DealPublicValue {
			Amount: bid.Value.Amount,
			Rate: bid.Value.Rate,
			Timestamp: time.Now().UTC().Unix(),
		},
	}

	dealPrivate := DealPrivate{
		Key: DealKey {
			ID: ID,
		},
		Value: DealPrivateValue {
			Borrower: borrower,
			Lender  : lender,
		},
	}

	if bytes, err := json.Marshal(dealPublic); err == nil {
		logger.Debug("Deal Public: " + string(bytes))
	}

	if bytes, err := json.Marshal(dealPrivate); err == nil {
		logger.Debug("Deal Private: " + string(bytes))
	}

	if err := bid.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	if err := dealPublic.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	if err := dealPrivate.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	logger.Info("MarketChaincode.makeDeal exited without errors")
	logger.Debug("Success: MarketChaincode.makeDeal")
	return shim.Success(nil)
}

func (t *MarketChaincode) queryBids(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.queryBids is running")
	logger.Debug("MarketChaincode.queryBids")

	it, err := stub.GetStateByPartialCompositeKey(bidIndex, []string{})
	if err != nil {
		message := fmt.Sprintf("unable to get state by partial composite key %s: %s", bidIndex, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}
	defer it.Close()

	entries := []Bid{}
	for it.HasNext() {
		response, err := it.Next()
		if err != nil {
			message := fmt.Sprintf("unable to get an element next to a queryBids iterator: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		logger.Debug(fmt.Sprintf("Response: {%s, %s}", response.Key, string(response.Value)))

		entry := Bid{}

		if err := entry.FillFromLedgerValue(response.Value); err != nil {
			message := fmt.Sprintf("cannot fill bid value from response value: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(response.Key)
		if err != nil {
			message := fmt.Sprintf("cannot split response key into composite key parts slice: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if err := entry.FillFromCompositeKeyParts(compositeKeyParts); err != nil {
			message := fmt.Sprintf("cannot fill bid key from composite key parts: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if bytes, err := json.Marshal(entry); err == nil {
			logger.Debug("Entry: " + string(bytes))
		}

		entries = append(entries, entry)
	}

	result, err := json.Marshal(entries)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Debug("Result: " + string(result))

	logger.Info("MarketChaincode.queryBids exited without errors")
	logger.Debug("Success: MarketChaincode.queryBids")
	return shim.Success(result)
}

func (t *MarketChaincode) queryBidsCreator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.queryBidsCreator is running")
	logger.Debug("MarketChaincode.queryBidsCreator")

	creator, err := GetCreatorOrganization(stub)
	if err != nil {
		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Debug("Creator: " + creator)

	it, err := stub.GetStateByPartialCompositeKey(bidIndex, []string{})
	if err != nil {
		message := fmt.Sprintf("unable to get state by partial composite key %s: %s", bidIndex, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}
	defer it.Close()

	entries := []Bid{}
	for it.HasNext() {
		response, err := it.Next()
		if err != nil {
			message := fmt.Sprintf("unable to get an element next to a queryBidsCreator iterator: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		logger.Debug(fmt.Sprintf("Response: {%s, %s}", response.Key, string(response.Value)))

		entry := Bid{}

		if err := entry.FillFromLedgerValue(response.Value); err != nil {
			message := fmt.Sprintf("cannot fill bid value from response value: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(response.Key)
		if err != nil {
			message := fmt.Sprintf("cannot split response key into composite key parts slice: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if err := entry.FillFromCompositeKeyParts(compositeKeyParts); err != nil {
			message := fmt.Sprintf("cannot fill bid key from composite key parts: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if bytes, err := json.Marshal(entry); err == nil {
			logger.Debug("Entry: " + string(bytes))
		}

		if creator == entry.Value.Creator {
			entries = append(entries, entry)
		}
	}

	result, err := json.Marshal(entries)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Debug("Result: " + string(result))

	logger.Info("MarketChaincode.queryBidsCreator exited without errors")
	logger.Debug("Success: MarketChaincode.queryBidsCreator")
	return shim.Success(result)
}

func (t *MarketChaincode) queryDeals(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.queryDeals is running")
	logger.Debug("MarketChaincode.queryDeals")

	it, err := stub.GetPrivateDataByPartialCompositeKey("collectionDeals", dealIndex, []string{})
	if err != nil {
		message := fmt.Sprintf("unable to get state by partial composite key %s: %s", dealIndex, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}
	defer it.Close()

	entries := []Deal{}
	for it.HasNext() {
		response, err := it.Next()
		if err != nil {
			message := fmt.Sprintf("unable to get an element next to a query iterator: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		logger.Debug(fmt.Sprintf("Response: {%s, %s}", response.Key, string(response.Value)))

		entry := Deal{}
		entryPublic := DealPublic{}
		//entryPrivate := DealPublic{}

		if err := entryPublic.FillFromLedgerValue(response.Value); err != nil {
			message := fmt.Sprintf("cannot fill deal public value from response value: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(response.Key)
		if err != nil {
			message := fmt.Sprintf("cannot split response key into composite key parts slice: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if err := entryPublic.FillFromCompositeKeyParts(compositeKeyParts); err != nil {
			message := fmt.Sprintf("cannot fill public deal key from composite key parts: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if bytes, err := json.Marshal(entryPublic); err == nil {
			logger.Debug("Entry: " + string(bytes))
		}

		entry.Key             = entryPublic.Key
		entry.Value.Amount    = entryPublic.Value.Amount
		entry.Value.Rate 	  = entryPublic.Value.Rate
		entry.Value.Timestamp = entryPublic.Value.Timestamp
		entry.Value.Borrower  = string("not available")
		entry.Value.Lender    = string("not available")

		entries = append(entries, entry)
	}

	result, err := json.Marshal(entries)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Debug("Result: " + string(result))

	logger.Info("MarketChaincode.queryDeals exited without errors")
	logger.Debug("Success: MarketChaincode.queryDeals")
	return shim.Success(result)
}


//func (t *MarketChaincode) queryDealsCreatorByTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	logger.Info("MarketChaincode.queryDealsCreatorByTime is running")
//	logger.Debug("MarketChaincode.queryDealsCreatorByTime")
//
//	if len(args) < timeBasicArgumentsNumber {
//		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
//			timeBasicArgumentsNumber, len(args))
//		logger.Error(message)
//		return shim.Error(message)
//	}
//
//	time := TimePeriod{}
//	if err := time.FillFromArguments(args[0:]); err != nil {
//		message := fmt.Sprintf("cannot fill a deal from arguments: %s", err.Error())
//		logger.Error(message)
//		return shim.Error(message)
//	}
//
//	logger.Debug("TimePeriodFrom: " + strconv.FormatInt(time.From,10))
//	logger.Debug("TimePeriodTo: " + strconv.FormatInt(time.To,10))
//
//	creator, err := GetCreatorOrganization(stub)
//	if err != nil {
//		message := fmt.Sprintf("cannot obtain creator's name from the certificate: %s", err.Error())
//		logger.Error(message)
//		return shim.Error(message)
//	}
//
//	logger.Debug("Creator: " + creator)
//
//	//it, err := stub.GetPrivateDataByPartialCompositeKey(string(collectionMembers[lenderT(members.Value.Lender)][borrowerT(members.Value.Borrower)]), dealIndex, []string{})
//	it, err := stub.GetPrivateDataByPartialCompositeKey("collectionDeals", dealIndex, []string{})
//	if err != nil {
//		message := fmt.Sprintf("unable to get state by partial composite key %s: %s", dealIndex, err.Error())
//		logger.Error(message)
//		return shim.Error(message)
//	}
//	defer it.Close()
//
//	entries := []Deal{}
//
//	for it.HasNext() {
//		response, err := it.Next()
//		if err != nil {
//			message := fmt.Sprintf("unable to get an element next to a query iterator: %s", err.Error())
//			logger.Error(message)
//			return shim.Error(message)
//		}
//
//		logger.Debug(fmt.Sprintf("Response: {%s, %s}", response.Key, string(response.Value)))
//
//		entry := DealPublic{}
//		members := Members{}
//
//		if err := entry.FillFromLedgerValue(response.Value); err != nil {
//			message := fmt.Sprintf("cannot fill deal value from response value: %s", err.Error())
//			logger.Error(message)
//			return shim.Error(message)
//		}
//
//		_, compositeKeyParts, err := stub.SplitCompositeKey(response.Key)
//		if err != nil {
//			message := fmt.Sprintf("cannot split response key into composite key parts slice: %s", err.Error())
//			logger.Error(message)
//			return shim.Error(message)
//		}
//
//		if err := entry.FillFromCompositeKeyParts(compositeKeyParts); err != nil {
//			message := fmt.Sprintf("cannot fill deal key from composite key parts: %s", err.Error())
//			logger.Error(message)
//			return shim.Error(message)
//		}
//
//		members.Key.ID =
//		creatorCollections := collectionMembers[lenderT(creator)]
//
//		for k := range creatorCollections {
//
//			logger.Debug(fmt.Sprintf("Get private collection: %s", string(creatorCollections[k])))
//			members, err := stub.GetPrivateData(string(creatorCollections[k]), compositeKey)
//			compositeKey, err := members.ToCompositeKey(stub)
//
//			//if (creator == entry.Value.Borrower || creator == entry.Value.Lender) && time.From <= entry.Value.Timestamp &&
//			//	time.To >= entry.Value.Timestamp{
//			//	entries = append(entries, entry)
//			//}
//		}
//
//		if bytes, err := json.Marshal(entry); err == nil {
//
//			logger.Debug("Entry: " + string(bytes))
//		}
//
//	}
//
//	result, err := json.Marshal(entries)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	logger.Debug("Result: " + string(result))
//
//	logger.Info("MarketChaincode.queryDealsCreatorByTime exited without errors")
//	logger.Debug("Success: MarketChaincode.queryDealsCreatorByTime")
//	return shim.Success(result)
//}
//
func getOrganization(certificate []byte) (string, error) {
	data := certificate[strings.Index(string(certificate), "-----") : strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	organization := cert.Issuer.Organization[0]
	return strings.Split(organization, ".")[0], nil
}

func GetCreatorOrganization(stub shim.ChaincodeStubInterface) (string, error) {
	certificate, err := stub.GetCreator()
	if err != nil {
		return "", err
	}
	return getOrganization(certificate)
}

func InitCollections(stub shim.ChaincodeStubInterface, CollectionConfig string) (error) {

	CollectionsFromConfig := []CollectionFromConfig{}
	var Col Collections

	creator := "a"

	if err := json.Unmarshal([]byte(CollectionConfig), &CollectionsFromConfig); err != nil {
		return err
	} else {
		for _, i := range CollectionsFromConfig{
			ok := strings.Contains(string(i.Policy), string(creator)+"MSP.member")
			if ok{
				Col.ListCollections = append(Col.ListCollections, Collection{i.Name})
			}
		}

		if bytes, err := json.Marshal(Col.ListCollections); err == nil {
			logger.Debug("Bid: " + string(bytes))
		}

		if err := Col.UpdateOrInsertIn(stub); err != nil {
			message := fmt.Sprintf("persistence error: %s", err.Error())
			logger.Error(message)
			return err
		}

		return nil
	}
}

func main() {
	err := shim.Start(new(MarketChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
