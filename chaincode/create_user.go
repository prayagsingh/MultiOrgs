package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MultiOrgs/chaincode/model"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(" ******** Invoke Create User ******** ")

	var queryCreatorOrg string
	var email, name, mobile string
	var eventID string
	var owner string

	var needHistory bool

	var timestamp *timestamp.Timestamp

	/* User Data Parameter */
	name = args[1]
	email = args[2]
	mobile = args[3]

	eventID = args[4]
	queryCreatorOrg = args[5]
	needHistory, _ = strconv.ParseBool(args[6])

	role, err := t.getRole(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to get roles from the account: %v", err))
	}

	userID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request owner: %v", err))
	}

	orgID, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request org: %v", err))
	}
	fmt.Println("\n value of orgID is: ", orgID)

	timestamp, err = stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Timestamp Error " + err.Error())
	}

	tm := model.GetTime(timestamp)

	user := &model.User{
		ID:     userID,
		Name:   name,
		Email:  email,
		Mobile: mobile,
		Owner:  queryCreatorOrg,
		Role:   role,
		Time:   tm,
	}

	userJSONBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := model.COLLECTION_KEY
	userNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{user.Email})
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(userNameIndexKey, userJSONBytes)
	if err != nil {
		return shim.Error("###### Error Put Private Create User Data Failed " + err.Error())
	}

	err = stub.SetEvent(eventID, []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(" ###### Create Data Parameters ###### ")
	fmt.Println(" ID 			= " + userID)
	fmt.Println(" Email			= " + email)
	fmt.Println(" Name 			= " + name)
	fmt.Println(" Mobile 			= " + mobile)
	fmt.Println(" Owner 			= " + owner)
	fmt.Println(" Role			= " + role)
	fmt.Println(" Time			= " + tm)
	fmt.Println(" ################################## ")

	/*	Created History for Create user Transaction */

	if needHistory {
		query := args[0]
		queryCreator := email
		remarks := email + " user created"
		t.createHistory(stub, queryCreator, queryCreatorOrg, email, query, remarks)
	}

	fmt.Println("User Invoked into the Ledger Successfully")

	return shim.Success(nil)
}
