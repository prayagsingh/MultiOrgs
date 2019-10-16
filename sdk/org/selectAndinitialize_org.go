package org

import (
	"fmt"

	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
)

// Orderer :
var Orderer Setup

// OrgList :
var OrgList []Setup

// OrgNames :
var OrgNames []string
var totalOrg = 3
var orgPeers []fabAPI.Peer
var signIdentities = make([]msp.SigningIdentity, 0, totalOrg-1)
var channelCtx contextAPI.ChannelProvider

// Init : for initializing all orgs
func (setup *Setup) Init(processAll bool) error {
	// OrgList type is slice --> make([]A, length, capacity)
	// capacity of the slice is the number of elements in the underlying array starting from the index from which the slice is created
	OrgList = make([]Setup, 0, totalOrg-1)
	OrgNames = []string{"Org1", "Org2"}

	if processAll {
		setup.InitializeAllOrgs()
	}

	return nil
}

// InitializeOrg : initializing Orgs with their values
func (setup *Setup) InitializeOrg(org string) (Setup, error) {
	var obj Setup

	switch name := org; name {

	case "Org1":
		obj = Setup{
			OrgAdmin:      "Admin",
			OrgName:       "Org1",
			ConfigFile:    "config-org1.yaml",
			OrgCaID:       "ca.org1.example.com",
			ChannelConfig: "Org1MSPanchors.tx",
		}
		break

	case "Org2":
		obj = Setup{
			OrgAdmin:      "Admin",
			OrgName:       "Org2",
			ConfigFile:    "config-org2.yaml",
			OrgCaID:       "ca.org2.example.com",
			ChannelConfig: "Org2MSPanchors.tx",
		}
		break
	}

	orgSetup, err := InitiateOrg(obj)
	if err != nil {
		return Setup{}, errors.WithMessage(err, " failed to initiate Org")
	}

	return orgSetup, nil
}

// InitiateOrg :
func InitiateOrg(obj Setup) (Setup, error) {
	org, err := initialize(obj)
	if err != nil {
		return Setup{}, fmt.Errorf("  failed to setup Org - " + obj.OrgName + " - " + err.Error())
	}

	if org == nil {
		return Setup{}, fmt.Errorf("  failed to setup Org - " + obj.OrgName)
	}

	fmt.Println("\n ****** Setup Created for " + org.OrgName + " ****** ")

	return *org, nil
}

// InitiateOrderer : setting up required info for Orderer org
func InitiateOrderer() (Setup, error) {

	obj := Setup{
		OrgAdmin:      "Admin",
		OrgName:       "OrdererOrg",
		ConfigFile:    "config-org1.yaml",
		ChannelConfig: "channel.tx",
	}

	orderer, err := initialize(obj)

	if err != nil {
		return Setup{}, fmt.Errorf("  failed to setup Org - " + obj.OrgName + " - " + err.Error())
	}

	if orderer == nil {
		return Setup{}, fmt.Errorf("  failed to setup Org - " + obj.OrgName)
	}

	Orderer = *orderer

	Orderer.OrdererID = "orderer.example.com"

	fmt.Println("\n **** Setup Created for " + Orderer.OrgName + " **** ")

	return Orderer, nil
}

// InitializeAllOrgs : initialize all orgs
func (setup *Setup) InitializeAllOrgs() error {
	fmt.Println("\n In selectAndInitialize.go file and Inside InitializeAllOrgs func")
	ordererSetup, err := InitiateOrderer()
	if err != nil {
		return errors.WithMessage(err, " failed to initiate Orderer")
	}
	OrgList = append(OrgList, ordererSetup)

	for _, org := range OrgNames {

		orgSetup, err := setup.InitializeOrg(org)
		if err != nil {
			return errors.WithMessage(err, " failed to initiate Org - "+org)
		}
		OrgList = append(OrgList, orgSetup)
	}

	Orderer.SigningIdentities = getSigningIdentities()
	OrgList[0] = Orderer

	return nil
}

// SelectOrg : allow to choose the org
func (setup *Setup) SelectOrg(org string) (*Setup, error) {

	fmt.Println("\n Inside SelectOrg method in selectAndinitialize.go file \n Input  Org " + org)

	orgSetup, err := setup.InitializeOrg(org)
	if err != nil {
		return nil, errors.WithMessage(err, "\n Falied to select Organization")
	}

	fmt.Println(" ########### Org Details ################# ")
	fmt.Println(" OrgName - " + orgSetup.OrgName)

	return &orgSetup, nil
}
