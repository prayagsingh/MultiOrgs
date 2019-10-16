package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MultiOrgs/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (t *SimpleChaincode) updateUserData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println(" ******** Invoke Update User ******** ")

	var user model.User
	var name, email, mobile, eventID string
	var queryCreatorOrg string
	var queryCreatorRole string
	var queryCreator string

	var needHistory bool

	/* User Data Parameter */
	name = args[1]
	email = args[2]
	mobile = args[3]
	eventID = args[4]
	queryCreatorOrg = args[5]
	queryCreatorRole = args[6]
	needHistory, _ = strconv.ParseBool(args[7])

	indexName := model.COLLECTION_KEY
	userNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{email})

	if err != nil {
		return shim.Error(err.Error())
	}

	err = getDataFromLedger(stub, userNameIndexKey, &user)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve userData in the ledger: %v", err))
	}

	userdata := &model.User{
		ID:     user.ID,
		Name:   name,
		Email:  user.Email,
		Mobile: mobile,
		Owner:  user.Owner,
		Role:   user.Role,
		Time:   user.Time,
	}

	userDataJSONasBytes, err := json.Marshal(userdata)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(userNameIndexKey, userDataJSONasBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(eventID, []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(" ###### Update Data Parameters ###### ")
	fmt.Println(" Email			= " + email)
	fmt.Println(" Name 			= " + name)
	fmt.Println(" Mobile 		= " + mobile)
	fmt.Println(" ################################## ")

	/*	Created History for Read by email Transaction */

	if needHistory {
		if strings.EqualFold(queryCreatorRole, model.ADMIN) {
			queryCreator = model.GetCustomOrgName(queryCreatorOrg) + " Admin"
		} else {
			queryCreator = email
		}

		fmt.Println(" ###### Query Access Details ###### ")
		fmt.Println(" queryCreatorRole = " + queryCreatorRole)
		fmt.Println(" queryCreator = " + queryCreator)
		fmt.Println(" ################################## ")

		var change []string

		if !strings.EqualFold(name, user.Name) {
			change = append(change, " Name to "+name+" , ")
		}

		if !strings.EqualFold(mobile, user.Mobile) {
			change = append(change, " Mobile number to "+mobile+" , ")
		}
		/*
			if !strings.EqualFold(age, user.Age) {
				change = append(change, " Age to "+age+" , ")
			}

			if !strings.EqualFold(salary, user.Salary) {
				change = append(change, " Salary to "+salary+" , ")
			}
		*/
		query := args[0]
		remarks := queryCreator + " has done following changes \n " + " [ " + strings.Join(change[:], "\n") + " ] "
		t.createHistory(stub, queryCreator, queryCreatorOrg, email, query, remarks)
	}

	return shim.Success(nil)
}
