package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MultiOrgs/sdk/invoke"

	"github.com/MultiOrgs/web/model"
)

// DeleteUserHandler : handler for deleting the user
func (app *RestApp) DeleteUserHandler() http.HandlerFunc {
	fmt.Println("\n In rest_delete_user.go file and Inside DeletUserHandler")
	return app.isAuthorized(func(w http.ResponseWriter, req *http.Request) {

		orgUser := app.OrgSetup.GetOrgUsers()
		if orgUser == nil {
			respondJSON(w, map[string]string{"error": "No Session Available"})
		} else {
			var userData model.UserData

			_ = json.NewDecoder(req.Body).Decode(&userData)

			email := userData.Email
			role := userData.Role
			orgName := userData.Org // owner of that user i.e Org1 or Org2

			fmt.Println("DeleteUserHandler : Email = " + email)

			orgInvoke := invoke.OrgInvoke{
				User: orgUser,
			}

			orgSetup, err := orgUser.OrgUserSetup.SelectOrg(orgName)
			if err != nil {
				respondJSON(w, map[string]string{"error": "Unable to select the org"})
			}

			err = orgUser.RemoveUser(email, orgSetup.OrgCaID, orgSetup.CaClient)
			if err != nil {
				fmt.Println("DeleteUserHandler : RemoveUser = Error : " + err.Error())
				respondJSON(w, map[string]string{"error": "Error Session User  " + err.Error()})
			} else {
				fmt.Println("Success RemoveUser ")
				// ReInitialize to Session Org

				_, err = orgUser.OrgUserSetup.SelectOrg(strings.ToLower(orgUser.OrgUserSetup.OrgName))

				user, _ := orgInvoke.GetUserFromLedger(email, false)

				if user != nil {
					err = orgInvoke.DeleteUserFromLedger(email, role)

					if err != nil {
						fmt.Println("DeleteUserHandler : Error Deleting User from ledger : " + err.Error())
						respondJSON(w, map[string]string{"error": "Error Deleting User from ledger " + err.Error()})
					}
					respondJSON(w, map[string]string{"success": "Succesfully delete the user with email - " + email})
				} else {
					respondJSON(w, map[string]string{"error": "No User Data Found"})
				}
			}
		}
	})
}
