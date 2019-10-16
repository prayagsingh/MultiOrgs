package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MultiOrgs/sdk/invoke"

	"github.com/MultiOrgs/web/model"
)

// UpdateUserHandler : handler for updating the user info
func (app *RestApp) UpdateUserHandler() http.HandlerFunc {

	return app.isAuthorized(func(w http.ResponseWriter, req *http.Request) {

		orgUser := app.OrgSetup.GetOrgUsers()
		if orgUser == nil {
			respondJSON(w, map[string]string{"error": "No Session Available"})
		} else {
			var userData model.UserData
			_ = json.NewDecoder(req.Body).Decode(&userData)

			name := userData.Name
			email := userData.Email
			mobile := userData.Mobile
			role := userData.Role

			fmt.Println(" ####### Rest Input for Update ####### ")

			fmt.Println(" Update Email	 	= " + email)
			fmt.Println(" Update Name 		= " + name)
			fmt.Println(" Update Mobile 	= " + mobile)
			fmt.Println(" Update Role 		= " + role)
			fmt.Println(" ###################################### ")

			orgInvoke := invoke.OrgInvoke{
				User: orgUser,
			}

			user, _ := orgInvoke.GetUserFromLedger(email, true)
			if user != nil {
				err := orgInvoke.UpdateUserFromLedger(email, name, mobile, role)
				if err != nil {
					respondJSON(w, map[string]string{"error": "Error Update User Data = " + err.Error()})
				} else {
					UserData, err := orgInvoke.GetUserFromLedger(email, false)

					if err != nil {
						respondJSON(w, map[string]string{"error": "No User Data Found"})
					} else {
						respondJSON(w, UserData)
					}
				}
			} else {
				respondJSON(w, map[string]string{"error": "No User Data Found"})
			}
		}
	})
}
