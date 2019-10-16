package org

import (
	"fmt"
	"os"
	"strings"

	caMsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"     // use to register the user with CA
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp" // helps in creation and updation of user

	// Package resmgmt enables creation and update of resources on a Fabric network
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	// Identity represents a Fabric client identity
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	// Client supplies the configuration and signing identity to client objects
	// ClientProvider returns client context
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	// Peer interface which provides MSPID() <-- gets the Peer mspid and URL() <--- gets the peer address
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

const (
	channelID   = "mychannel"
	ccPath      = "github.com/MultiOrgs/chaincode/"
	ccID        = "mycc"
	ccVersion   = "1"
	ccPolicy    = "OR ('Org1MSP.member','Org2MSP.member')"
	ordererName = "OrdererOrg"
	ordererID   = "orderer.example.com"
)

// Setup : basic config
type Setup struct {

	// Network parameters
	OrdererID             string
	OrdererAdmin          string
	OrdererName           string
	OrdererClientContext  contextAPI.ClientProvider
	OrdererChannelContext contextAPI.ChannelProvider
	OrdererResmgmt        *resmgmt.Client

	// Channel parameters
	ChannelID     string
	ChannelConfig string

	// Chaincode parameters
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeId      string
	ChainCodeVersion string
	ChainCodePolicy  string

	CCPkg *resource.CCPackage

	ConfigFile string
	OrgCaID    string
	OrgName    string
	OrgAdmin   string
	UserName   string

	Sdk               *fabsdk.FabricSDK
	CaClient          *caMsp.Client
	Resmgmt           *resmgmt.Client
	Ctx               contextAPI.ClientProvider
	MspClient         *mspclient.Client
	Peers             []fabAPI.Peer
	ChannelContext    contextAPI.ChannelProvider
	ChannelClient     *channel.Client
	Event             *event.Client
	SigningIdentity   msp.SigningIdentity
	SigningIdentities []msp.SigningIdentity
}

// Initialize reads the configuration file and sets up the client
func initialize(setup Setup) (*Setup, error) {
	fmt.Println("\n#### Initialize " + setup.OrgName + " SDK ####")
	// Reading file
	//fmt.Println("\n Value of setup.ConfigFile is: ", setup.ConfigFile)
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read config.yaml file hence failed to create SDK")
	}
	fmt.Println("\n#### SDK created for " + setup.OrgName + " ####")

	// Register the User with CA
	// New: creates a new client Instance
	caClient, err := caMsp.New(sdk.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to create new CA client: %v", err)
	}
	fmt.Println("  CA Client created for " + setup.OrgName)

	// Get resource management context
	orgCtx := sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	fmt.Println("\n#### Context created for " + setup.OrgName)

	// New returns a resource management client instance
	resMgmtClient, err := resmgmt.New(orgCtx)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create resmgmt")
	}
	fmt.Println("  Resource management client created for " + setup.OrgName)

	// creating a new client instance. mspClient allows to retrieve user info from the identity
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create MSP client")
	}
	fmt.Println("\n MSP Client created for " + setup.OrgName)
	fmt.Println("\n Value of OrgAdmin is: " + setup.OrgAdmin)

	signingIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get admin signing identity")
	}
	//fmt.Println("\n Value of signing identity is: ", signingIdentity)
	fmt.Println("\n Signing Identity created for " + setup.OrgName)

	// if the OrgName is not same as ordererName then we are creating a channel
	if !strings.EqualFold(setup.OrgName, ordererName) {
		fmt.Println("\n #### In org_setup file and Inside EqualFold condition")
		signIdentities = append(signIdentities, signingIdentity)
		// returns Peer interface which provides MSPID and URL functions
		orgPeers, err = DiscoverLocalPeers(orgCtx, 2)
		if err != nil {
			_ = fmt.Errorf(" failed to Discover Local Peers: %v for "+setup.OrgName, err)
			return nil, nil
		}
		fmt.Println("\nPeers Discovered for " + setup.OrgName)
		fmt.Println("\n Inside org_setup and Value of channelID is: ", channelID)

		// Channel client is used to query and execute transactions
		channelCtx = sdk.ChannelContext(channelID,
			fabsdk.WithUser(setup.OrgAdmin),
			fabsdk.WithOrg(setup.OrgName))
		fmt.Println(" In org_setup.go file and Channel Client created for " + setup.OrgName)
	}

	// list of signing identities
	signingIdentities := []msp.SigningIdentity{signingIdentity}
	fmt.Println("\n value of ChannelConfig is: ", setup.ChannelConfig)

	return &Setup{
		ConfigFile:        setup.ConfigFile,
		ChannelID:         channelID,
		ChaincodeGoPath:   os.Getenv("GOPATH"),
		ChaincodePath:     ccPath,
		ChaincodeId:       ccID,
		ChainCodeVersion:  ccVersion,
		ChainCodePolicy:   ccPolicy,
		OrdererName:       ordererName,
		OrdererID:         ordererID,
		ChannelClient:     nil,
		Event:             nil,
		OrgCaID:           setup.OrgCaID,
		OrgName:           setup.OrgName,
		OrgAdmin:          setup.OrgAdmin,
		ChannelConfig:     getArtifactPath() + setup.ChannelConfig,
		Sdk:               sdk,
		CaClient:          caClient,
		Ctx:               orgCtx,
		Resmgmt:           resMgmtClient,
		MspClient:         mspClient,
		SigningIdentities: signingIdentities,
		Peers:             orgPeers,
		ChannelContext:    channelCtx,
	}, nil
}

func getArtifactPath() string {
	return os.Getenv("GOPATH") + "/src/github.com/MultiOrgs/networks/channel-artifacts/"
}

func getSigningIdentities() []msp.SigningIdentity {
	fmt.Println("\n\ngetSigningIdentities == ", len(signIdentities))
	return signIdentities
}
