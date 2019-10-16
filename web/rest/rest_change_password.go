package rest

import (
	"encoding/json"
	"net/http"

	"github.com/MultiOrgs/web/model"
)

// ChangePwdHandler : handler for managing the password changing request
func (app *RestApp) ChangePwdHandler() func(http.ResponseWriter, *http.Request) {

	return app.isAuthorized(func(w http.ResponseWriter, r *http.Request) {

		orgUser := app.OrgSetup.GetOrgUsers()

		if orgUser == nil {
			respondJSON(w, map[string]string{"error": "Error Session User "})
		} else {

			var userdata model.UserData

			_ = json.NewDecoder(r.Body).Decode(&userdata)

			email := userdata.Email
			role := userdata.Role
			name := userdata.Name
			oldPwd := hash(userdata.OldPassword)
			newPwd := hash(userdata.Password)

			verifyErr := verifyPassword(userdata.Password)

			if verifyErr != nil && len(verifyErr.Error()) > 0 {
				respondJSON(w, map[string]string{"error": verifyErr.Error(), "message": "Password must contain at least one number and one uppercase and lowercase letter, and at least 8 or more characters"})
			} else {

				err := orgUser.OrgUserSetup.ChangePassword(email, role, name, oldPwd, newPwd)

				if err != nil {
					respondJSON(w, map[string]string{"error": "Unable to Change user pwd - " + err.Error()})
				} else {
					respondJSON(w, map[string]string{"success": "Password successfully changed for - " + email})
				}
			}
		}
	})
}
