package org

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// CreateChaincodePkg : created a chaincode package
func (setup *Setup) CreateChaincodePkg() (*resource.CCPackage, error) {
	fmt.Println(" Creating Chaincode Package with ChaincodePath " + setup.ChaincodePath + " ")
	fmt.Println(" 	- ChaincodeGoPath " + setup.ChaincodeGoPath)
	// NewCCPackage creates new go lang chaincode package
	CCPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create chaincode package")
	}
	setup.CCPkg = CCPkg
	fmt.Println("ccPkg created")
	return CCPkg, nil
}

// InstallCC : for installing chaincode on peers
func (setup *Setup) InstallCC(ccPkg *resource.CCPackage) ([]resmgmt.InstallCCResponse, error) {
	fmt.Println("\n ##### Installing Chaincode on Peers #####")

	// checking that if the OrgName is equal to OrderName then return
	if strings.EqualFold(setup.OrgName, setup.OrdererName) || len(setup.OrgName) == 0 {
		return nil, errors.WithMessage(nil, "setup.OrgName: "+setup.OrgName+" is not same as OrdererName: "+setup.OrdererName)
	}
	// Ensure that Gossip has propagated it's view of local peers before invoking
	// install since some peers may be missed if we call InstallCC too early
	orgPeers, err := DiscoverLocalPeers(setup.Ctx, 2)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to Discover Local Peers for "+setup.OrgName)
	}
	setup.Peers = orgPeers
	fmt.Println("  Peers Discovered for " + setup.OrgName)
	fmt.Println("\n  Installing Chaincode for " + setup.OrgName)
	fmt.Println("\n Value of setup.ChaincodeId is: ", setup.ChaincodeId,
		"\n Value of setup.Chaincode path is: ", setup.ChaincodePath,
		"\n Value of setup.ChaincodeVersion is: ", setup.ChainCodeVersion)

	if setup.CCPkg == nil {
		fmt.Println("\n Chaincode Pkg for Org: ", setup.OrgName, " is nil")
		return nil, errors.WithMessage(nil, "Chaincode Pkg is nil")
	}
	installCCReq := resmgmt.InstallCCRequest{
		Name:    setup.ChaincodeId,
		Path:    setup.ChaincodePath,
		Version: setup.ChainCodeVersion, //chaincode version. this is the very first version of our chaincode
		Package: setup.CCPkg,
	}

	resp, err := setup.Resmgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return resp, errors.WithMessage(err, "failed to install chaincode")
	}
	//fmt.Println("Value of resp.Status is: ", resp[1])
	fmt.Println("Chaincode installed for Org: ", setup.OrgName)
	return resp, nil
}

// InstantiateCC : initialize the chaincode
// Returns: upgrade chaincode response with transaction ID
func (setup *Setup) InstantiateCC(orgPeers []fab.Peer) error {
	if strings.EqualFold(setup.OrgName, setup.OrdererName) || len(setup.OrgName) == 0 {
		return nil
	}
	fmt.Println("\n  Instantiating Chaincode for ", setup.OrgName)
	fmt.Println("\nInstantiateCC CC Policy: ", setup.ChainCodePolicy)
	fmt.Println("\nInstantiateCC CC Name: ", setup.ChaincodeId)
	fmt.Println("\nInstantiateCC CC ChaincodePath: ", setup.ChaincodePath)
	fmt.Println("\nOrgPeers are:", orgPeers[0], " and ", orgPeers[1])

	ccPolicy, err := cauthdsl.FromString(setup.ChainCodePolicy)
	if err != nil {
		fmt.Println("failed policy : " + err.Error())
		return errors.WithMessage(err, "failed policy")
	}

	resp, err := setup.Resmgmt.InstantiateCC(
		setup.ChannelID,
		resmgmt.InstantiateCCRequest{

			Name:    setup.ChaincodeId,
			Path:    setup.ChaincodePath,
			Version: setup.ChainCodeVersion,
			Args:    [][]byte{[]byte("init"), []byte("Org1"), []byte("100"), []byte("Org2"), []byte("200")},
			Policy:  ccPolicy,
		}, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(orgPeers[0], orgPeers[1]))
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode for "+setup.OrgName)
	}
	fmt.Println("Chaincode instantiated and resp.TxId is: ", resp.TransactionID)
	fmt.Println("Chaincode Installation & Instantiation Successful")

	return nil
}

// UpgradeCC : upgrades the chaincode
// Returns: upgrade chaincode response with transaction ID
func (setup *Setup) UpgradeCC(orgPeers []fab.Peer) error {
	fmt.Println("\n##### Upgrading Chaincode")
	fmt.Println("Upgrade CC Policy with policy: " + setup.ChainCodePolicy)
	fmt.Println("Upgrade CC Name: " + setup.ChaincodeId)
	fmt.Println("Upgrade CC Version:  " + setup.ChainCodeVersion)
	fmt.Println("Upgrade CC ChaincodePath: " + setup.ChaincodePath)

	ccPolicy, err := cauthdsl.FromString(setup.ChainCodePolicy)
	if err != nil {
		return errors.WithMessage(err, "failed Chaincode policy")
	}

	// UpgradeCCRequest contains upgrade chaincode request parameters
	req := resmgmt.UpgradeCCRequest{
		Name:    setup.ChaincodeId,
		Version: setup.ChainCodeVersion,
		Path:    setup.ChaincodePath,
		Args:    [][]byte{[]byte("init"), []byte("Org1"), []byte("350"), []byte("Org2"), []byte("500")},
		Policy:  ccPolicy,
	}

	resp, err := setup.Resmgmt.UpgradeCC(setup.ChannelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(orgPeers[0], orgPeers[1]))

	if resp.TransactionID == "" || err != nil {
		return errors.WithMessage(err, "\n ####### failed to upgrade chaincode, no transaction ID")
	}

	return nil
}

// QueryInstalledCC : use to check if chaincode is installed or not on peers
func (setup *Setup) QueryInstalledCC(orgID, ccName, ccVersion, ccPath string, orgPeers []fab.Peer) (bool, error) {
	installedOnAllPeers := true
	for _, peer := range orgPeers {
		fmt.Println("\n Querying chaincode installed status for peer: ", peer.URL())
		chaincodeQueryResponse, err := setup.Resmgmt.QueryInstalledChaincodes(resmgmt.WithTargets(peer))
		if err != nil {
			return false, errors.WithMessage(err, "\nQueryChaincodeInstalled for peer '"+peer.URL()+"' failed")
		}
		found := false
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			if chaincode.Name == ccName && chaincode.Version == ccVersion {
				fmt.Println("   " + orgID + " found chaincode " + chaincode.Name + ": " + ccName + " with version " +
					chaincode.Version + ": " + ccVersion)
				found = true
				break
			}
		}

		if !found {
			fmt.Println("   " + orgID + " chaincode is not instatiated on peer " + peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers, nil
}

// QueryInstantiatedCC : use to check if chaincode is instantiated or not on peers
func (setup *Setup) QueryInstantiatedCC(channelID, orgID, ccName, ccVersion, ccPath string, orgPeers []fab.Peer) (bool, error) {
	installedOnAllPeers := true
	for _, peer := range orgPeers {
		fmt.Println("\n Querying chaincode installed status for peer: ", peer.URL())
		chaincodeQueryResponse, err := setup.Resmgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(peer))
		if err != nil {
			return false, errors.WithMessagef(err, " QueryChaincodeInstalled for peer '"+peer.URL()+"' failed")
		}
		found := false
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			if chaincode.Name == ccName && chaincode.Version == ccVersion && chaincode.Path == ccPath {
				fmt.Println("   " + orgID + " found chaincode " + chaincode.Name + " --- " + ccName + " with version " +
					chaincode.Version + " -- " + ccVersion)
				found = true
				break
			}
		}

		if !found {
			fmt.Println("   " + orgID + " chaincode is not installed on peer " + peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers, nil
}

// TestInvoke :
func (setup *Setup) TestInvoke(org string) error {

	fmt.Println(" ********** Test Invoke - " + org + " **********")

	eventID := "testInvoke-" + org

	orgSetup, _ := setup.SelectOrg(org)
	orgName := orgSetup.OrgName
	//orgSdk := orgSetup.Sdk
	orgAdmin := orgSetup.OrgAdmin
	//caClient := orgSetup.CaClient
	channelClient, event, _ := setup.CreateChannelClient(orgName, orgAdmin)

	response, err := setup.ExecuteChaincodeTranctionEvent(eventID, "invoke", [][]byte{
		[]byte("testInvoke"),
		[]byte(eventID),
	}, setup.ChaincodeId, channelClient, event)

	if err != nil {
		return fmt.Errorf("Error - Test Invoke failed for "+org+" : %s", err.Error())
	}
	fmt.Println(" ********** For " + org + ", Test Invoke Successful with txn ID: " + string(response.TransactionID) + "********** ")

	return nil
}

// TestQuery :
func (setup *Setup) TestQuery(org string) (string, string, error) {

	fmt.Println(" ********** Test Query with " + org + " **********")

	orgSetup, _ := setup.SelectOrg(org)
	orgName := orgSetup.OrgName
	//orgSdk := orgSetup.Sdk
	orgAdmin := orgSetup.OrgAdmin
	//caClient := orgSetup.CaClient
	channelClient, _, _ := setup.CreateChannelClient(orgName, orgAdmin)

	fmt.Println("\n\n value of setup.ChaincodeId is: ", setup.ChaincodeId)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data to invoke in the ledger")

	response, err := channelClient.Query(channel.Request{
		ChaincodeID:  setup.ChaincodeId,
		Fcn:          "invoke",
		Args:         [][]byte{[]byte("query"), []byte(orgName)},
		TransientMap: transientDataMap,
	}, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargets(setup.Peers[0], setup.Peers[1]))

	if err != nil {
		return "", "", fmt.Errorf("Error - Test Query failed for "+org+" : %s", err.Error())
	}
	fmt.Println("For Org: " + org + ", total Amount is: " + string(response.Payload))
	fmt.Println(" ********** " + org + " Test Query Successful ********** ")

	return string(response.Payload), string(response.TransactionID), nil
}
