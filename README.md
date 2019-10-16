# MultiOrgs
This project is a Multi-Organization setup built using [HLFv1.4.2](https://github.com/hyperledger/fabric/tree/v1.4.2) and [GO-SDKv1.0.0-alpha5](https://github.com/hyperledger/fabric-sdk-go/tree/v1.0.0-alpha5). In this project, we have three organizations `OrdererOrganization`, `Org1` and `Org2`. So `Org1` and `Org2` are like a bank with some amount associated with each organization. We can add a new `User` with some amount say `X` to an organization and it's amount will be added to the total amount of the organization. when a user withdrawls the amount then its amount will be removed from the organization.

## Getting Started
#### Check for vendor directory in chaincode directory and project directory itself
1. Go to `chaincode` directory and check for `vendor` directory. If it is not present then run command `dep ensure` for creating a `vendor` directory **OR** 

2. If the vendor directory is present then go to `chaincode/vendor/github.com/docker/docker/integration-cli/fixtures/https` directory and run `tar -zxvf certs.tar` command.

Now you are good to go. 

#### Steps to start the network
