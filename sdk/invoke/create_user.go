package invoke

import (
	"fmt"
	"strconv"

	"github.com/MultiOrgs/sdk/org"
)

// OrgInvoke : storing info for creating User
type OrgInvoke struct {
	User *org.OrgUser
	Role string
}

// InvokeCreateUser : creating User on Ledger
func (orginvoke *OrgInvoke) InvokeCreateUser(name, mobile string) error {

	fmt.Println(" ############## Invoke Create User ################")

	var queryCreatorOrg string

	queryCreatorOrg = orginvoke.User.OrgUserSetup.OrgName
	email := orginvoke.User.Username
	eventID := "userInvoke"
	needHistory := strconv.FormatBool(true)

	// Trying a separate channel client and event.
	// RegistrationCahincodeEvent isn't working if using email id to create a channelContext, for admin its is working fine
	//channelClient, event, _ := orginvoke.User.OrgUserSetup.CreateChannelClient(orginvoke.User.OrgUserSetup.OrgName, email)

	fmt.Println(" ###### Create Data Parameters ###### ")
	fmt.Println(" Email 			= " + email)
	fmt.Println(" Name 			= " + name)
	fmt.Println(" Mobile 			= " + mobile)
	fmt.Println(" Owner 			= " + queryCreatorOrg)
	fmt.Println(" Value of orginvoke.User.OrgUserSetup.ChaincodeId: ", orginvoke.User.OrgUserSetup.ChaincodeId)
	fmt.Printf("\n Value of orginvoke.User.ChannelClient: %v", *orginvoke.User.ChannelClient)
	fmt.Printf("\n Value of orginvoke.User.Event: %v", *orginvoke.User.Event)
	//fmt.Printf("\n Value of new channel client is %v ", *channelClient)
	//fmt.Printf("\n Value of new event is %v ", *event)
	fmt.Println(" ")

	_, err := orginvoke.User.OrgUserSetup.ExecuteChaincodeTranctionEvent(eventID, "invoke",
		[][]byte{
			[]byte("createUser"),
			[]byte(name),
			[]byte(email),
			[]byte(mobile),
			[]byte(eventID),
			[]byte(queryCreatorOrg),
			[]byte(needHistory),
		}, orginvoke.User.OrgUserSetup.ChaincodeId, orginvoke.User.ChannelClient, orginvoke.User.Event)

	if err != nil {
		return fmt.Errorf("Error - addUserToLedger : %s", err.Error())
	}

	fmt.Println("#### User added Successfully ####")
	return nil
}
