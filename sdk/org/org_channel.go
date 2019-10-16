package org

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

// CreateChannel : use to create a channel for organizations
func (setup *Setup) CreateChannel() error {
	fmt.Println("\n Inside CreateChannel method in org_channel.go file ")
	fmt.Println("Create Channel for Org - "+setup.OrgName, len(setup.OrgName))

	if len(setup.OrgName) == 0 {
		return errors.WithMessage(nil, " empty Org Name")
	}

	var OrgJoined bool
	var err error
	fmt.Println("\n Value of setup.OrgName is: ", setup.OrgName)
	fmt.Println("\n Value of setup.OrdererName is: ", setup.OrdererName)

	if !strings.EqualFold(setup.OrgName, setup.OrdererName) {
		// returns boolean val if channel is already joined or not
		OrgJoined, err = setup.IsJoinedChannel(setup.Resmgmt, setup.Peers[0])
		if err != nil {
			fmt.Println("failed to check isJoin channel")
			return errors.WithMessage(err, "  failed to check isJoin channel")
		}
	}

	if !OrgJoined {

		fmt.Println("Create Channel == ChannelID = " + setup.ChannelID)
		fmt.Println("Create Channel == OrdererID = " + Orderer.OrdererID)
		//fmt.Println("Create Channel using setup.OrdererID == OrdererID = "+setup.OrdererID)
		fmt.Println("Create Channel == channelConfigPath = " + setup.ChannelConfig)
		fmt.Println("Create Channel == SigningIdentities are = ", len(setup.SigningIdentities))

		// SaveChannelRequest struct for creating channel
		req := resmgmt.SaveChannelRequest{
			ChannelID:         setup.ChannelID,
			ChannelConfigPath: setup.ChannelConfig,     // path of channel.tx, Org1MSPanchor, Org2MSPAnchor
			SigningIdentities: setup.SigningIdentities, // signing identities of all the orgs
		}

		// creating channel with ordererID. All other Orgs will join this channel
		txID, err := setup.Resmgmt.SaveChannel(req, resmgmt.WithOrdererEndpoint(Orderer.OrdererID))
		if err != nil || txID.TransactionID == "" {
			return errors.WithMessage(err, "failed to save anchor channel for - "+setup.OrgName)
		}

		var lastConfigBlock uint64
		lastConfigBlock, err = WaitForOrdererConfigUpdate(setup.Resmgmt, setup.ChannelID, true, lastConfigBlock)

		if err != nil {
			return errors.WithMessage(err, "failed to get Orderer config update")
		}

		fmt.Printf("Channel Orderer Config Update %lld ", lastConfigBlock)

	} else {
		fmt.Println(" Peers Already Joined channel")
	}

	return nil
}

// JoinChannelForOrg : organizations joining channel
func (setup *Setup) JoinChannelForOrg() error {

	fmt.Println("Join Channel for Org - "+setup.OrgName, len(setup.OrgName))

	if len(setup.OrgName) == 0 {
		return errors.WithMessage(nil, " empty Org Name")
	}

	var OrgJoined bool
	var err error

	if !strings.EqualFold(setup.OrgName, setup.OrdererName) {

		OrgJoined, err = setup.IsJoinedChannel(setup.Resmgmt, setup.Peers[0])
		if err != nil {
			fmt.Println("failed to check isJoin channel")
			return errors.WithMessage(err, "  failed to check isJoin channel")
		}
	}

	if !OrgJoined {

		fmt.Println("Initiating Join Channel - " + setup.OrgName)

		fmt.Println("  JoinChannel " + setup.ChannelID)

		if err := setup.Resmgmt.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(Orderer.OrdererID)); err != nil {
			return errors.WithMessage(err, "failed to make admin join channel")
		}

		fmt.Println("  Successfully Joined the Channel " + setup.ChannelID)

	} else {

		fmt.Println(" Peers Already Joined channel")
	}

	return nil
}

// IsJoinedChannel : checking whether the channel is joined or not
func (setup *Setup) IsJoinedChannel(orgResmgmt *resmgmt.Client, peer fabAPI.Peer) (bool, error) {

	resp, err := orgResmgmt.QueryChannels(resmgmt.WithTargets(peer))
	if err != nil {
		fmt.Println("IsJoinedChannel : failed to Query >>> " + err.Error())
		return false, err
	}
	for _, chInfo := range resp.Channels {
		fmt.Println("IsJoinedChannel : " + chInfo.ChannelId + " --- " + setup.ChannelID)
		if chInfo.ChannelId == setup.ChannelID {
			return true, nil
		}
	}
	return false, nil
}
