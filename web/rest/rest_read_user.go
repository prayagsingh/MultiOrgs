package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MultiOrgs/sdk/invoke"

	"github.com/MultiOrgs/web/model"
)

// GetUserDataByEmailHandler : handler for getting user via email
func (app *RestApp) GetUserDataByEmailHandler() http.HandlerFunc {

	return app.isAuthorized(func(w http.ResponseWriter, req *http.Request) {

		orgUser := app.OrgSetup.GetOrgUsers()
		if orgUser == nil {
			respondJSON(w, map[string]string{"error": "No Session Available"})
		} else {
			var userdata model.UserData
			_ = json.NewDecoder(req.Body).Decode(&userdata)
			email := userdata.Email

			fmt.Println(" Session User - " + orgUser.Username)
			fmt.Println(" Session OrgName - " + orgUser.OrgUserSetup.OrgName)

			orgInvoke := invoke.OrgInvoke{
				User: orgUser,
			}

			UserData, err := orgInvoke.GetUserFromLedger(email, true)
			if err != nil {
				respondJSON(w, map[string]string{"error": "No User Data Found"})
			} else {
				respondJSON(w, UserData)
			}
		}
	})
}

// GetAllUsersDataHandler : get all user data handler
func (app *RestApp) GetAllUsersDataHandler() http.HandlerFunc {
	return app.isAuthorized(func(w http.ResponseWriter, req *http.Request) {

		orgUser := app.OrgSetup.GetOrgUsers()
		if orgUser == nil {
			respondJSON(w, map[string]string{"error": "No Session Available"})
		} else {
			orgInvoke := invoke.OrgInvoke{
				User: orgUser,
			}

			allUsersData, err := orgInvoke.GetAllUserFromLedger()
			if err != nil {
				respondJSON(w, map[string]string{"error": "No User Data Found"})
			} else {
				respondJSON(w, allUsersData)
			}
		}
	})
}
