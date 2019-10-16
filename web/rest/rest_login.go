package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MultiOrgs/sdk/invoke"

	"github.com/MultiOrgs/web/model"
)

// LoginHandler : handles login related queries
func (app *RestApp) LoginHandler(w http.ResponseWriter, req *http.Request) {

	var userData model.UserData

	_ = json.NewDecoder(req.Body).Decode(&userData)

	orgName := userData.Org
	email := userData.Email
	password := userData.Password
	//role := userData.Role
	fmt.Println("\n value of userData is: ", userData)

	Org, err := app.OrgSetup.InitializeOrg(orgName)
	if err != nil {
		respondJSON(w, map[string]string{"error": "failed to invoke user " + err.Error()})
	}

	fmt.Println("Sign In --->  emailValue = " + email)

	orgUser, err := Org.LoginUserWithCA(email, password)

	orgInvoke := invoke.OrgInvoke{
		User: orgUser,
	}

	if err != nil {
		respondJSON(w, map[string]string{"error": "Unable to Login : " + err.Error()})
	} else {
		fmt.Println("Logged In User : " + orgUser.Username)
		token := app.processAuthentication(w, email)

		if len(token) > 0 {
			UserData, err := orgInvoke.GetUserFromLedger(email, true)

			if err != nil {
				respondJSON(w, map[string]string{"error": "No User Data Found"})
			} else {
				respondJSON(w, map[string]string{

					"token":  token,
					"name":   UserData.Name,
					"email":  UserData.Email,
					"mobile": UserData.Mobile,
				})
			}
		} else {
			respondJSON(w, map[string]string{"error": "Failed to generate token"})
		}
	}
}
