package org

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
)

// ExecuteChaincodeTranctionEvent : execute chaincode event services for invoking chaincode
func (setup *Setup) ExecuteChaincodeTranctionEvent(eventID, fcnName string, args [][]byte, chaincodeID string, channelClient *channel.Client, ccEvent *event.Client) (*channel.Response, error) {
	fmt.Println(" ############# ExecuteChaincodeTranctionEvent - " + eventID + " ############## ")

	fmt.Println("  Execute Org -- " + setup.OrgName)
	fmt.Println("  Execute CCId -- " + chaincodeID)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data to invoke in the ledger")

	if ccEvent == nil {
		fmt.Println(" ############### Event is Nil")
	}

	// registering the chaincode event. after registering the event, always deregister it
	fmt.Println("\n ######## Registering for Chaincode Event ###### ")
	registration, notifier, err := ccEvent.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("In chaincodeEvents.go and found error when registering chaincode event")
		return nil, fmt.Errorf("Blockchain ..... failed to register event: %v", err)
	}
	fmt.Println("\n ###### RegistrationChaincodeEvent executed successfully")

	// deregister the chaincode event
	defer ccEvent.Unregister(registration)

	// Create a request (proposal) and send it to the endorser peers and returns the proposal responses from peer(s)
	fmt.Println("\n In chaincodeEvents.go and creating a proposal(request) ")
	response, err := channelClient.Execute(channel.Request{
		ChaincodeID:  chaincodeID,
		Fcn:          fcnName,
		Args:         args,
		TransientMap: transientDataMap,
	}, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargets(setup.Peers[0], setup.Peers[1]))

	if err != nil {
		return nil, fmt.Errorf("failed to Invoke request and unable to move funds: %v", err)
	}

	if response.ChaincodeStatus == 0 {
		return nil, errors.WithMessage(nil, "Expected ChaincodeStatus")
	}

	if response.Responses[0].ChaincodeStatus != response.ChaincodeStatus {
		return nil, errors.WithMessage(nil, "Expected the chaincode status returned by successful Peer Endorsement to be same as Chaincode status for client response")
	}

	select {
	case ccEventNotify := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEventNotify)
	case <-time.After(time.Second * 100):
		return nil, fmt.Errorf("timeout while waiting for chaincode event hence did NOT receive CC event for eventId(%s)", eventID)
	}
	fmt.Println("\n\n completed ExecuteChaincodeTranctionEvent functionalities")
	return &response, nil
}
