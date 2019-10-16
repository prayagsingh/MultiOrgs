package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/MultiOrgs/chaincode/model"
	"github.com/MultiOrgs/sdk/invoke"
	"github.com/MultiOrgs/sdk/org"
)

// AllUsers : all users
type AllUsers struct {
	Users []model.User `json:"users"`
}

// SampleUsers : creating sample users
func SampleUsers() {

	var org string
	fmt.Println("Enter the Org - ( Org1, Org2) ")
	fmt.Scanln(&org)

	byteValue, _ := ioutil.ReadFile("sample/" + org + ".json")

	users := []model.User{}

	v := []byte(byteValue)

	err := json.Unmarshal(v, &users)

	if err != nil {
		fmt.Printf("unable to convert the result to object: %v", err)
	}

	threads := len(users)
	var wg sync.WaitGroup

	respond := make(chan string, threads)
	wg.Add(threads)

	for i := 0; i < len(users); i++ {

		email := users[i].Email
		name := users[i].Name
		mobile := users[i].Mobile
		owner := users[i].Owner
		role := users[i].Role

		go createUser(respond, &wg, email, name, mobile, owner, role)
	}

	wg.Wait()
	close(respond)

	for queryResp := range respond {
		fmt.Println("Query Response: " + queryResp)
	}
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	return sha1Hash
}

func createUser(respond chan<- string, wg *sync.WaitGroup, email, name, mobile, owner, role string) {

	fmt.Println(" Create User ---- " + email + " ---- Org = " + owner)

	password := hash("Test@#123")

	setup := &org.Setup{}
	_ = setup.Init(false)

	Org, err := setup.InitializeOrg(owner)
	if err != nil {
		defer wg.Done()
		respond <- fmt.Sprintf("Web ----- Unable to initialize org  - %s", err.Error())
	} else {

		orgUser, err := Org.RegisterUserWithCA(owner, email, name, password, role)

		orgInvoke := invoke.OrgInvoke{
			User: orgUser,
		}

		if err != nil {

			defer wg.Done()
			respond <- fmt.Sprintf("Web Error ----->>> Unable to Register Error Msg  - %s", err.Error())

		} else {

			err := orgInvoke.InvokeCreateUser(name, mobile)

			defer wg.Done()
			if err != nil {
				respond <- fmt.Sprintf("failed to invoke user  - %s", err.Error())
			} else {
				respond <- fmt.Sprintf("User Created Successfully  - %s", email)
			}

		}
	}
}
