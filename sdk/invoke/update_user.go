package invoke

import (
	"fmt"
	"strconv"
)

// UpdateUserFromLedger : update user
func (orginvoke *OrgInvoke) UpdateUserFromLedger(email, name, mobile, role string) error {

	fmt.Println(" ############## Invoke Update Data ################")

	eventID := "updateInvoke"
	queryCreatorOrg := orginvoke.User.OrgUserSetup.OrgName
	queryCreatorRole := orginvoke.Role
	needHistory := strconv.FormatBool(true)

	_, err := orginvoke.User.OrgUserSetup.ExecuteChaincodeTranctionEvent(eventID, "invoke",
		[][]byte{
			[]byte("updateUserData"),
			[]byte(name),
			[]byte(email),
			[]byte(mobile),
			[]byte(eventID),
			[]byte(queryCreatorOrg),
			[]byte(queryCreatorRole),
			[]byte(needHistory),
		}, orginvoke.User.OrgUserSetup.ChaincodeId, orginvoke.User.OrgUserSetup.ChannelClient, orginvoke.User.OrgUserSetup.Event)

	if err != nil {
		fmt.Errorf(" Error - Update User Data From Ledger : %s ", err.Error())
	}

	fmt.Println(" ###################################################### ")

	return nil
}
