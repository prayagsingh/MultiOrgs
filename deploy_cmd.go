package main

import (
	"fmt"
	"strings"

	"github.com/MultiOrgs/sdk/org"
	"github.com/pkg/errors"
)

// DeployCMD : deploy cmd
func DeployCMD(setup *org.Setup, cmd string) error {

	fmt.Println("Deploy CMD " + cmd)

	err := setup.Init(true)
	if err != nil {
		return errors.WithMessage(err, "  failed to initialize  : "+err.Error())
	}
	fmt.Println("\nInside deplay_cmd.go and SDK initialized successfully")
	fmt.Println("\n Inside deploy_cmd and value of OrgList is: ", org.OrgList)

	switch id := cmd; id {

	case "1":

		/*  Channel Create for Org */
		for _, s := range org.OrgList {
			fmt.Print("\n Ranging for Org: ", s)
			err = s.CreateChannel()
			if err != nil {
				return errors.WithMessage(err, "  failed to Create Channel "+err.Error())
			}
		}

		break

	case "2":

		for _, s := range org.OrgList {

			if !strings.EqualFold(s.OrgName, s.OrdererName) {

				err = s.JoinChannelForOrg()
				if err != nil {
					return errors.WithMessage(err, "  failed to Join Channel "+err.Error())
				}

			}
		}

		break

	case "3":

		/* Install Chaincode */
		fmt.Println("\n Value of OrgList[0] is: ", org.OrgList[0],
			"\n Value of OrgList[1] is: ", org.OrgList[1],
			"\n Value of OrgList[2] is: ", org.OrgList[2])
		ccPkg, err := org.OrgList[1].CreateChaincodePkg()
		if err != nil {
			fmt.Println("  failed to create chaincode Pkg : " + err.Error())
			return errors.WithMessage(err, "  failed to create chaincode Pkg : "+err.Error())
		}

		Org1 := org.OrgList[1]
		_, err = Org1.InstallCC(ccPkg)
		if err != nil {
			return errors.WithMessage(err, "  failed to Install Chaincode "+" : "+err.Error())
		}
		fmt.Println("Install CC successfull for - " + Org1.OrgName)

		ccPkg, err = org.OrgList[2].CreateChaincodePkg()
		Org2 := org.OrgList[2]
		_, err = Org2.InstallCC(ccPkg)
		if err != nil {
			return errors.WithMessage(err, "  failed to Install Chaincode "+" : "+err.Error())
		}
		fmt.Println("Install CC successfull for - " + Org2.OrgName)
		break

	case "4":

		/* Instantiate Chaincode */

		Org1 := org.OrgList[1]
		org1Peers := Org1.Peers

		err = Org1.InstantiateCC(org1Peers)
		if err != nil {
			return errors.WithMessage(err, "  failed to Instantiate Chaincode "+" : "+err.Error())
		}
		fmt.Println("Instantiate CC successfull ")

		for _, s := range org.OrgList {

			if !strings.EqualFold(s.OrgName, s.OrdererName) {

				err := s.TestInvoke(s.OrgName)

				if err != nil {
					fmt.Println(" Invoke failed for - " + s.OrgName + " : " + err.Error())
				}

				amount, txID, err := s.TestQuery(s.OrgName)
				fmt.Println("\n Value of Query for Org: " + s.OrgName + " is: " + amount)
				fmt.Println("\n Value of Query-TxID for Org: " + s.OrgName + " is: " + txID)
				if err != nil {
					fmt.Println(" Invoke failed for - " + s.OrgName + " : " + err.Error())
				}
			}
		}

		break

	case "5":

		var testOrg string
		fmt.Println(" Enter the Org name- ( org1, org2) ")
		fmt.Scanln(&testOrg)

		s := org.OrgList[2]

		err := s.TestInvoke(testOrg)

		if err != nil {
			fmt.Println(" Invoke failed for - " + testOrg + " : " + err.Error())
		}

		break

	case "6":

		Org1 := org.OrgList[1]
		Org2 := org.OrgList[2]

		ccPkg, err := Org1.CreateChaincodePkg()
		if err != nil {
			fmt.Println("  failed to create chaincode Pkg : " + err.Error())
			return errors.WithMessage(err, "  failed to create chaincode Pkg : "+err.Error())
		}

		_, err = Org1.InstallCC(ccPkg)
		if err != nil {
			return errors.WithMessage(err, "  failed to Install Chaincode "+" : "+err.Error())
		}
		fmt.Println("Install CC successfull for - " + Org1.OrgName)

		_, err = Org2.InstallCC(ccPkg)
		if err != nil {
			return errors.WithMessage(err, "  failed to Install Chaincode "+" : "+err.Error())
		}
		fmt.Println("Install CC successfull for - " + Org2.OrgName)

		org1Peers := Org1.Peers

		err = Org1.UpgradeCC(org1Peers)
		if err != nil {
			return errors.WithMessage(err, "  failed to Upgrade Chaincode "+" : "+err.Error())
		}
		fmt.Println("Upgrade CC successfull ")

		for _, s := range org.OrgList {

			if !strings.EqualFold(s.OrgName, s.OrdererName) {

				err := s.TestInvoke(s.OrgName)

				if err != nil {
					fmt.Println(" Invoke failed for - " + s.OrgName + " : " + err.Error())
				}

			}
		}

		break

	case "7":

		for _, s := range org.OrgList {

			_, err = s.QueryInstalledCC(s.OrgName, s.ChaincodeId, s.ChainCodeVersion, s.ChaincodePath, s.Peers)
			if err != nil {
				return errors.WithMessage(err, "  failed to Query Installed Chaincode "+" : "+err.Error())
			}
			fmt.Println("Query CC successfull ")
		}
		break

	case "8":

		for _, s := range org.OrgList {

			_, err = s.QueryInstantiatedCC(s.ChannelID, s.OrgName, s.ChaincodeId, s.ChainCodeVersion, s.ChaincodePath, s.Peers)
			if err != nil {
				return errors.WithMessage(err, "  failed to Query Instantiated Chaincode "+" : "+err.Error())
			}
			fmt.Println("Query CC successfull ")
		}
		break
		/*
			case "9":

				Org3 := org.OrgList[3]

				err := Org3.AddAffiliationOrg()

				if err != nil {
					return errors.WithMessage(err, "  failed to affiliate org "+" : "+err.Error())
				}

				Org4 := org.OrgList[4]

				err = Org4.AddAffiliationOrg()

				if err != nil {
					return errors.WithMessage(err, "  failed to affiliate org "+" : "+err.Error())
				}

				break

			}
		*/
	}
	return nil
}
