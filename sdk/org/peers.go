package org

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/pkg/errors"
)

// DiscoverLocalPeers : returns the list of peers
func DiscoverLocalPeers(ctxProvider context.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {

	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		fmt.Println("    error creating local context :  %s " + err.Error())
		return nil, errors.Wrap(err, "error creating local context")
	}
	// NewInvoker creates a new Retryable Invoker
	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(

		func() (interface{}, error) {

			// LocalDiscoveryService returns core discovery service
			// Returns: localDiscovery which points to "fab.DiscoveryService" in fab/provider.go
			// DiscoveryService(Interface) is used to discover eligible peers on specific channel
			// using GetPeers()
			peers, err := ctx.LocalDiscoveryService().GetPeers()

			if err != nil {
				fmt.Println("    error getting peers for MSP :  %s " + err.Error())
				return nil, errors.Wrapf(err, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}

			fmt.Println("  MSP ID -- "+ctx.Identifier().MSPID, expectedPeers, len(peers))

			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}

// LoadOrgPeers : create context to load the peers
func LoadOrgPeers(ctxProvider context.ClientProvider) error {

	fmt.Println("    LoadOrgPeers")

	_, err := contextImpl.NewLocal(ctxProvider)

	if err != nil {
		fmt.Printf("\ncontext creation failed : %s ", err.Error())
		return errors.WithMessage(err, "context creation failed")
	}

	return nil
}
