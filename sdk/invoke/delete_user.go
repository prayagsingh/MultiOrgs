package invoke

import (
	"fmt"
	"strconv"
)

// DeleteUserFromLedger : delete user info
func (orginvoke *OrgInvoke) DeleteUserFromLedger(email, role string) error {

	fmt.Println(" ############## Invoke Delete User ################")

	eventID := "deleteInvoke"
	queryCreatorOrg := orginvoke.User.OrgUserSetup.OrgName
	queryCreatorRole := orginvoke.Role
	needHistory := strconv.FormatBool(true)

	_, err := orginvoke.User.OrgUserSetup.ExecuteChaincodeTranctionEvent(eventID, "invoke",
		[][]byte{
			[]byte("deleteUser"),
			[]byte(email),
			[]byte(eventID),
			[]byte(role),
			[]byte(queryCreatorOrg),
			[]byte(queryCreatorRole),
			[]byte(needHistory),
		}, orginvoke.User.OrgUserSetup.ChaincodeId, orginvoke.User.OrgUserSetup.ChannelClient, orginvoke.User.OrgUserSetup.Event)

	if err != nil {
		fmt.Errorf("Error - DeleteUserFromLedger : %s", err.Error())
	}

	return nil
}
