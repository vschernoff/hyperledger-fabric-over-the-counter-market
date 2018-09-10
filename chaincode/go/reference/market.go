
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
	"fmt"
	"time"
	"strings"
	"encoding/pem"
	"crypto/x509"
	"encoding/json"
)

var logger = shim.NewLogger("MarketChaincode")

// MarketChaincode example simple Chaincode implementation
type MarketChaincode struct {
}

func (t *MarketChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")

	_, args := stub.GetFunctionAndParameters()
	if len(args) > 0 {
		logger.Debug("Organization names: " + strings.Join(args, ", "))

		compositeKey, err := stub.CreateCompositeKey(orgsIndex, []string{})
		if err != nil {
			message := fmt.Sprintf("unable to create orgs composite key: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		value, err := json.Marshal(args)
		if err != nil {
			message := fmt.Sprintf("unable to marshal organization names: %s", err.Error())
			logger.Error(message)
			return shim.Error(message)
		}

		if err := stub.PutState(compositeKey, value); err != nil {
			message := fmt.Sprintf("persistence error: %s", err.Error())
			logger.Error(message)
			return pb.Response{Status: 500, Message: message}
		}
	}

	return shim.Success(nil)
}

func (t *MarketChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")

	function, args := stub.GetFunctionAndParameters()
	if function == "placeBid" {
		// place a bid on the market
		return t.placeBid(stub, args)
	} else if function == "cancelBid" {
		return t.cancelBid(stub, args)
	} else if function == "makeDeal" {
		return t.makeDeal(stub, args)
	} else if function == "queryBids" {
		return t.queryBids(stub, args)
	} else if function == "queryDeals" {
		return t.queryDeals(stub, args)
	} else if function == "queryFromCollection" {
		return t.queryFromCollection(stub, args)
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

	deal := Deal {
		Key: DealKey {
			ID: uuid.Must(uuid.NewV4()).String(),
		},
		Value: DealValue {
			Borrower: borrower,
			Lender: lender,
			Amount: bid.Value.Amount,
			Rate: bid.Value.Rate,
			Timestamp: time.Now().UTC().Unix(),
		},
	}

	if bytes, err := json.Marshal(deal); err == nil {
		logger.Debug("Deal: " + string(bytes))
	}

	if err := bid.UpdateOrInsertIn(stub); err != nil {
		message := fmt.Sprintf("persistence error: %s", err.Error())
		logger.Error(message)
		return pb.Response{Status: 500, Message: message}
	}

	if err := deal.UpdateOrInsertIn(stub); err != nil {
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

func (t *MarketChaincode) queryDeals(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.queryDeals is running")
	logger.Debug("MarketChaincode.queryDeals")

	it, err := stub.GetStateByPartialCompositeKey(dealIndex, []string{})
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

		if err := entry.FillFromLedgerValue(response.Value); err != nil {
			message := fmt.Sprintf("cannot fill deal value from response value: %s", err.Error())
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
			message := fmt.Sprintf("cannot fill deal key from composite key parts: %s", err.Error())
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

	logger.Info("MarketChaincode.queryDeals exited without errors")
	logger.Debug("Success: MarketChaincode.queryDeals")
	return shim.Success(result)
}

func (t *MarketChaincode) queryFromCollection(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("MarketChaincode.queryFromCollection is running")
	logger.Debug("MarketChaincode.queryFromCollection")

	if len(args) < 1 {
		message := fmt.Sprintf("insufficient number of arguments: expected %d, got %d",
			1, len(args))
		logger.Error(message)
		return shim.Error(message)
	}

	collection := args[0]

	it, err := stub.GetPrivateDataByPartialCompositeKey(collection, dealIndex, []string{})
	if err != nil {
		message := fmt.Sprintf("unable to get state from collection %s by partial composite key %s: %s",
			collection, dealIndex, err.Error())
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

		if err := entry.FillFromLedgerValue(response.Value); err != nil {
			message := fmt.Sprintf("cannot fill deal value from response value: %s", err.Error())
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
			message := fmt.Sprintf("cannot fill deal key from composite key parts: %s", err.Error())
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

	logger.Info("MarketChaincode.queryFromCollection exited without errors")
	logger.Debug("Success: MarketChaincode.queryFromCollection")
	return shim.Success(result)
}

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

func main() {
	err := shim.Start(new(MarketChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
