package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MultiOrgs/web/rest"
	"github.com/gorilla/mux"
)

// ServeRestAPI : manages all the handler for rest calls
func ServeRestAPI(app *rest.RestApp) {

	r := mux.NewRouter()

	// read User
	r.HandleFunc("/api/read_users", app.GetAllUsersDataHandler()).Methods("GET")
	r.HandleFunc("/api/read_user", app.GetUserDataByEmailHandler()).Methods("GET")

	// register and login user
	r.HandleFunc("/api/register_user", app.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login_user", app.LoginHandler).Methods("POST")

	// update and delete user
	r.HandleFunc("/api/update_user", app.UpdateUserHandler()).Methods("PUT")
	r.HandleFunc("/api/delete_user", app.DeleteUserHandler()).Methods("DELETE")

	// change password
	r.HandleFunc("/api/change_password", app.ChangePwdHandler()).Methods("POST")

	fmt.Println("Listening (http://localhost:5050) ...")
	log.Fatal(http.ListenAndServe(":5050", r))
}
