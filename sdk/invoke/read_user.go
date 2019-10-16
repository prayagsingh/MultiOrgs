package invoke

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MultiOrgs/chaincode/model"
)

// GetAllUserFromLedger : get all users data
func (orginvoke *OrgInvoke) GetAllUserFromLedger() ([]model.User, error) {
	fmt.Println(" ############## Invoke Read All User ################")

	eventID := "getAllUsersInvoke"

	response, err := orginvoke.User.OrgUserSetup.ExecuteChaincodeTranctionEvent(eventID, "invoke",
		[][]byte{
			[]byte("readAllUser"),
			[]byte(eventID),
		}, orginvoke.User.OrgUserSetup.ChaincodeId, orginvoke.User.ChannelClient, orginvoke.User.Event)

	if err != nil {
		return nil, fmt.Errorf("Error - addUserToLedger : %s", err.Error())
	}

	fmt.Println("Response Received")

	allUsers := make([]model.User, 0)

	if response != nil && response.Payload == nil {
		return nil, fmt.Errorf("unable to get response for the query: %v", err)
	}

	if response != nil {
		err = json.Unmarshal(response.Payload, &allUsers)
		if err != nil {
			return nil, fmt.Errorf("unable to convert response to the object given for the query: %v", err)
		}
	}

	if len(allUsers) < 1 {
		return nil, fmt.Errorf("No records found")
	}

	return allUsers, nil
}

// GetUserFromLedger : get a single user
func (orginvoke *OrgInvoke) GetUserFromLedger(email string, needHistory bool) (*model.User, error) {

	fmt.Println(" ############## Invoke Read User From Ledger ################")

	eventID := "getUserInvoke"
	queryCreatorOrg := orginvoke.User.OrgUserSetup.OrgName

	fmt.Println(" Email = " + email)

	fmt.Println(" Need History - " + strconv.FormatBool(needHistory))

	response, err := orginvoke.User.OrgUserSetup.ExecuteChaincodeTranctionEvent(eventID, "invoke",
		[][]byte{
			[]byte("readUser"),
			[]byte(email),
			[]byte(eventID),
			[]byte(queryCreatorOrg),
			[]byte(strconv.FormatBool(needHistory)),
		}, orginvoke.User.OrgUserSetup.ChaincodeId, orginvoke.User.ChannelClient, orginvoke.User.Event)

	if err != nil {
		return nil, fmt.Errorf("Error - Get User From Ledger : %s", err.Error())
	}

	if response == nil {
		return nil, fmt.Errorf("Error - No User found ")
	}

	var user *model.User

	err = json.Unmarshal(response.Payload, &user)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to the object given for the query: %v", err)
	}

	fmt.Println("#### User Found #### ")
	fmt.Println(" Email 	= " + user.Email)
	fmt.Println(" Mobile  	= " + user.Mobile)
	fmt.Println(" ################################## ")

	return user, nil
}
