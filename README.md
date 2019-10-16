# MultiOrgs
This project is a Multi-Organization setup built using [HLFv1.4.2](https://github.com/hyperledger/fabric/tree/v1.4.2) and [GO-SDKv1.0.0-alpha5](https://github.com/hyperledger/fabric-sdk-go/tree/v1.0.0-alpha5). In this project, we have three organizations `OrdererOrganization`, `Org1` and `Org2`. So `Org1` and `Org2` are like a bank with some amount associated with each organization. We can add a new `User` with some amount say `X` to an organization and it's amount will be added to the total amount of the organization. when a user withdrawls the amount then its amount will be removed from the organization.

## Getting Started
1. Clone this project in `$GOPATH/src/github.com` directory

2. Check for `vendor` directory in `$GOPATH/src/github.com/MultiOrgs/chaincode` directory and `$GOPATH/src/github.com/MultiOrgs` directory itself.
  
      2.a: Go to `chaincode` directory and check for `vendor` directory. If it is not present then run command `dep ensure` for creating a `vendor` directory **OR** 

      2.b: If the vendor directory is present then go to `chaincode/vendor/github.com/docker/docker/integration-cli/fixtures/https` directory and run `tar -zxvf certs.tar` command.

### Prerequisites

1. GoLang: `go1.13.1 linux/amd64`

2. Docker: `18.09.3, build 774a1f4`

3. Docker-compose: `1.24.0, build 0aa59064`

4. Hyperledger-Fabric: `v1.4.2`

5. Fabric-Go-Sdk: `v1.0.0-alpha5 revision: 0e710cebffe5674d411a625bc103ff0e266bb5e1`

6. Using `Dep` for dependency.

## Installing 
##### Go to `$GOPATH/src/github.com/MultiOrgs` directory.
1. Run `make` command. Operations carried out by `make` command are:-

   a. Remove previous containers, network and volumes from docker. 
   
   b. Remove previous certs and keys from `wallet` directory in case you are restarting the network.
   
   c. Create the fresh docker containers, network and volume. 
   
   d. Deploy the network.
   
    d.1: Create channel.
    
    d.2: Join channel --> channel joined by every organization.
    
    d.3: Install chaincode on `peer0` of both the organization.
    
    d.4: Instantiate the chaincode.
    
  e. Start the Rest Server. Functionality offered by rest server are:
    
    e.1: Create new Identity using api
    
    curl --header "Content-Type: application/json"  --request POST  --data '{"Email": "beta@gmail.com","Name": "beta","Mobile": "+91 1234567891","Owner": "Org_Name", "Role": "admin","No": "1","Org": "Org_Name","Password": "Beta@123"}' http://localhost:5050/api/register_user
    
    
    e.2 Login User 
    
    curl --header "Content-Type: application/json" -H "Authorization: Bearer <ACCESS_TOKEN_GENERATED_BY_REGISTERING_USER>" --request POST --data '{"Org": "Org1", "email": "alpha@gmail.com", "password": "<sha256(password)>"}' http://localhost:5050/api/login_user
    
    
    e.3: Read a single user in an Organization
    
    curl --header "Content-Type: application/json" -H "Authorization: Bearer <Access_Token>" --request GET --data '{"email": "beta@gmail.com"}' http://localhost:5050/api/read_user
    
    
    e.4: Read multiple users
    
    curl --header "Content-Type: application/json" -H "Authorization: Bearer <Access_Token>" --request GET --data '{"email": "beta@gmail.com"}' http://localhost:5050/api/read_users
    
    
    e.5: Update the user info
    
    curl --header "Content-Type: application/json" -H "Authorization: Bearer <Access_Token>" --request PUT --data '{"Email": "beta@gmail.com","Name": "Beta","Mobile": "+91 1234567892", "Role": "client"}' http://localhost:5050/api/update_user
    
    
    e.6: Delete user
    
    curl --header "Content-Type: application/json" -H "Authorization: Bearer <Access_Token>" --request DELETE --data '{"Email": "beta@gmail.com", "Role": "user", "Org":"Org_Name"}' http://localhost:5050/api/delete_user
