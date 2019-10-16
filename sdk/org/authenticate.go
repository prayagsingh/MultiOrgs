package org

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	caMsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// OrgUser struct
type OrgUser struct {
	Username      string
	ChannelClient *channel.Client // Client enables access to a channel on a Fabric network.
	Event         *event.Client   // Client enables access to a channel events on a Fabric network.
	OrgUserSetup  Setup
}

var sessionOrgUser *OrgUser
var sessionOrgName = make(map[string]string)
var sessionUser = make(map[string]string)
var secretKey = make(map[string]string)

//GetOrgUsers : get the user
func (setup *Setup) GetOrgUsers() *OrgUser {
	return sessionOrgUser
}

// RegisterUserWithCA : register new User with their role and Org Name
func (setup *Setup) RegisterUserWithCA(org, email, name, password, role string) (*OrgUser, error) {

	fmt.Println("\n Inside RegisterUserWithCA method in authenticate.go file")
	orgSetup, _ := setup.SelectOrg(org) // returns Channel client created

	caID := orgSetup.OrgCaID
	caClient := orgSetup.CaClient
	affiliation := strings.ToLower(org) + ".department1"

	fmt.Println("CA Register Org      === " + org)
	fmt.Println("CA Register CaID     === " + caID)
	fmt.Println("CA Register Email 	  === " + email)
	fmt.Println("CA Register Name 	  === " + name)
	fmt.Println("CA Register Password === " + password)
	fmt.Println("CA Register Role 	  === " + role)
	fmt.Println("Affiliation === " + affiliation)
	fmt.Println("CA Register CaClient 	  === ", caClient)

	// Get CaInfo
	var mp interface{}
	mp, _ = caClient.GetCAInfo()
	fmt.Printf("\n Inside authenticate.go and value of GetCAInfo is: %s", mp)

	// Register registers a User with the Fabric CA
	//  Parameters:
	//  request is registration request
	//
	//  Returns:
	//  enrolment secret
	// RegistrationRequest defines the attributes required to register a user with the CA
	registerSecret, err := caClient.Register(&caMsp.RegistrationRequest{
		Name: email, // Name is the unique name of the identity
		// Secret is an optional password.  If not specified, a random secret is generated.  In both cases, the secret
		// is returned from registration.
		Secret: password,
		Type:   "peer", // Type of identity being registered (e.g. "peer, app, user")
		// MaxEnrollments is the number of times the secret can be reused to enroll.
		// if omitted, this defaults to max_enrollments configured on the server
		MaxEnrollments: -1,          // -1 means infinite enrollment, 0 is no enrollment allowed
		Affiliation:    affiliation, // The identity's affiliation e.g. org1.department1
		// Optional attributes associated with this identity
		// refer to this link: https://hyperledger-fabric-ca.readthedocs.io/en/latest/users-guide.html#attribute-based-access-control
		Attributes: []caMsp.Attribute{
			{
				Name:  "role", // <-- Need to change the t.Role() in chaincode dir if using any variable instead of static value
				Value: role,
				ECert: true,
			},
		},
		CAName: caID, // CAName is the name of the CA to connect to
	})
	fmt.Println("\n Value of registerSecret is: ", registerSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to register user with CA '%s': %v", email, err)
	}

	sessionOrgName["orgName"] = org
	sessionUser["name"] = email
	secretKey["secret"] = password

	fmt.Println("\n Successfully register user ", email)
	fmt.Println("\n In auth.go and value of secretKey is: ", secretKey["secret"])

	orgUser, err := setup.LoginUserWithCA(email, password)
	if err != nil {
		return nil, fmt.Errorf("unable to login '%s': %v", email, err)
	}
	fmt.Println("\n ##### Org Register Name: " + orgUser.OrgUserSetup.OrgName)

	return orgUser, nil
}

// LoginUserWithCA :
func (setup *Setup) LoginUserWithCA(email, password string) (*OrgUser, error) {
	fmt.Println("\n ####### Login User ####### ")

	caClient := setup.CaClient

	// Enroll enrolls a registered user in order to receive a signed X509 certificate.
	// A new key pair is generated for the user. The private key and the
	// enrollment certificate issued by the CA are stored in SDK stores.
	// They can be retrieved by calling IdentityManager.GetSigningIdentity().
	//attrs := []*AttributeRequest{{Name: "role", Optional: true}
	//Attributes := []*caMsp.AttributeRequest{{Name: "role", Optional: true}}  use -->  caMsp.WithAttributeRequests(Attributes)
	err := caClient.Enroll(email, caMsp.WithSecret(password), caMsp.WithType("peer"))
	if err != nil {
		return nil, fmt.Errorf("failed to enroll identity '%s': %v", email, err)
	}

	// For testing purpose only
	orgEnrolledUser, err := caClient.GetSigningIdentity(email)

	if orgEnrolledUser.Identifier().ID != email {
		return nil, fmt.Errorf("In authenticate.go and Enrolled user name doesn't match")
	}

	if setup.OrgName == "Org1" {
		if orgEnrolledUser.Identifier().MSPID != "Org1MSP" {
			return nil, fmt.Errorf("In authenticate.go and Enrolled user mspID doesn't match")
		}
	} else if setup.OrgName == "Org2" {
		if orgEnrolledUser.Identifier().MSPID != "Org2MSP" {
			return nil, fmt.Errorf("In authenticate.go and Enrolled user mspID doesn't match")
		}
	}
	fmt.Println("\n Value of orgEnrolledUser.Identifier().ID is: ", orgEnrolledUser.Identifier().ID)

	sessionOrgName["orgName"] = setup.OrgName
	sessionUser["name"] = email
	secretKey["secret"] = password

	fmt.Println("\n ###### Org: " + setup.OrgName)

	channelClient, event, err := setup.CreateChannelClient(setup.OrgName, email)

	if err != nil {
		return nil, fmt.Errorf("unable to create channel client '%s': %v", email, err)
	}

	fmt.Println("\n ####### Org Enroll Name: " + setup.OrgName)

	sessionOrgUser = &OrgUser{
		Username:      email,
		ChannelClient: channelClient,
		Event:         event,
		OrgUserSetup:  *setup,
	}

	return sessionOrgUser, nil

}

// ChangePassword : method for changing the password
func (setup *Setup) ChangePassword(email, role, name, oldPwd, newPwd string) error {

	//fmt.Println("Change PWD : Email = " + email + " , OLD PWD = " + oldPwd + " , Saved PWD = " + secretKey["secret"] + " ,  New PWD = " + newPwd)

	if !strings.EqualFold(oldPwd, secretKey["secret"]) {
		return fmt.Errorf("Old password don't matched, can't change pwd for the email: '%s'", email)
	}

	if strings.EqualFold(oldPwd, newPwd) {
		return fmt.Errorf("failed old password, new password should not be same: '%s'", email)
	}

	orgUser := setup.GetOrgUsers()

	err := orgUser.RemoveUserFromCA(email, setup.OrgCaID, setup.CaClient)

	if err != nil {
		return fmt.Errorf("failed to remove identity '%s': %v", email, err)
	}

	fmt.Println("User Removed")

	org := setup.OrgName

	_, err = setup.RegisterUserWithCA(org, email, name, newPwd, role)

	if err != nil {
		return fmt.Errorf("failed to register with CA '%s': %v", email, err)
	}

	fmt.Println("User Re-Registered")

	_, err = setup.LoginUserWithCA(email, newPwd)

	if err != nil {
		return fmt.Errorf("failed to enroll user '%s': %v", email, err)
	}

	fmt.Println("User Enrolled")

	err = setup.ReEnrollUser(email)

	if err != nil {
		return fmt.Errorf("failed to re-enroll user '%s': %v", email, err)
	}

	fmt.Println("User Re-Enrolled")

	return nil
}

// ReEnrollUser : again enrolling/adding the user
func (setup *Setup) ReEnrollUser(email string) error {

	err := setup.CaClient.Reenroll(email)

	if err != nil {
		return fmt.Errorf("failed to re-enroll user '%s': %v", email, err)
	}
	return nil
}

//################################################################
//###########         REMOVE USER SECTION            #############
//################################################################

// RemoveUserFromCA : removes the user from the organization
func (orguser *OrgUser) RemoveUserFromCA(email string, caID string, caClient *caMsp.Client) error {
	_, err := caClient.RemoveIdentity(&caMsp.RemoveIdentityRequest{

		ID:     email,
		Force:  true,
		CAName: caID,
	})

	if err != nil {
		return fmt.Errorf("failed to remove signing identity for '%s': %v", email, err)
	}
	return nil
}

// RemoveUser : remove user
func (orguser *OrgUser) RemoveUser(email string, caID string, caClient *caMsp.Client) error {
	err := orguser.RemoveUserFromCA(email, caID, caClient)
	if err != nil {
		return fmt.Errorf("failed to remove signing identity for '%s': %v", email, err)
	}

	return nil
}

// RevokeUser : black listing the user
func (orguser *OrgUser) RevokeUser(email string) error {

	_, err := orguser.OrgUserSetup.CaClient.Revoke(&caMsp.RevocationRequest{
		Name: email,
	})

	if err != nil {
		return fmt.Errorf("failed to revoke signing identity for '%s': %v", email, err)
	}

	return nil
}

// CreateChannelClient : enabling ChannelContext and Events services
func (setup *Setup) CreateChannelClient(org string, email string) (*channel.Client, *event.Client, error) {

	SigningIndentity, err := setup.CaClient.GetSigningIdentity(email)
	if err != nil {
		fmt.Println(" failed to get signing identity: " + email + "\n" + err.Error())
		return nil, nil, fmt.Errorf("failed to get signing identity for '%s': %v", email, err)
	}
	fmt.Println("In file authenticate.go and Signing Identity Created for " + email)
	fmt.Println("\n Previous channelContext created in org_setup.go file is: ", setup.ChannelContext)
	fmt.Println("\n Inside authenticate.go and ChannelID is: ", setup.ChannelID)
	fmt.Println("\n Inside authentication.go and org is: ", org)

	// ChannelContext is used to create the channel client and event client objects
	// using channel client and event client, we can query, execute txn and call channel events respectively.
	clientContext := setup.Sdk.ChannelContext(setup.ChannelID,
		fabsdk.WithUser(email),
		fabsdk.WithOrg(org),
		fabsdk.WithIdentity(SigningIndentity))

	fmt.Println("\n In authentication.go and Channel Context Created and value is:", &clientContext)

	// New returns a Client instance.
	// Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	ChannelClient, err := channel.New(clientContext)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new channel client for '%s': %v", email, err)
	}

	setup.ChannelClient = ChannelClient
	fmt.Println("In authenticate.go file and Channel client created")

	// Creation of the client which will enables access to our channel events
	event, err := event.New(clientContext)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new event client %v", err)
	}
	setup.Event = event
	fmt.Println("Event client created")

	return setup.ChannelClient, setup.Event, nil
}
