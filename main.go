package main

import (
	"fmt"
	"strings"

	"github.com/MultiOrgs/sdk/org"
	"github.com/MultiOrgs/web"
	"github.com/MultiOrgs/web/rest"
)

func main() {

	fmt.Println(" Choose the following ")
	fmt.Println(" 1. Deploy the network")
	fmt.Println(" 2. Start the Rest Server (Listening (http://localhost:5050) ...)")
	//fmt.Println(" 3. Start the Html Web App (Listening (http://localhost:6000) ...)")
	fmt.Println(" 3. Create Dummy Users")

	var choose string

	fmt.Scanln(&choose)

	setup := &org.Setup{}
	_ = setup.Init(false)

	if strings.EqualFold(choose, "1") {

		fmt.Println(" Deployement of a network")
		fmt.Println("   1.  Create Channel")
		fmt.Println("   2.  Join Channel")
		fmt.Println("   3.  Install Chaincode")
		fmt.Println("   4.  Instantiate Chaincode")
		fmt.Println("   5.  Test Invoke")
		fmt.Println("   6.  Upgrade Chaincode")
		fmt.Println("   7.  Query Installed Chaincode")
		fmt.Println("   8.  Query Instantiate Chaincode")
		fmt.Println("   9.  Affiliate an Org")

		var cmd string
		fmt.Scanln(&cmd)

		err := DeployCMD(&org.Setup{}, cmd)
		if err != nil {
			fmt.Println(" setup Failed " + err.Error())
			return
		}
	}

	if strings.EqualFold(choose, "2") {

		app := &rest.RestApp{
			OrgSetup: setup,
		}
		web.ServeRestAPI(app)
	}

	if strings.EqualFold(choose, "3") {

		SampleUsers()

	}
}
