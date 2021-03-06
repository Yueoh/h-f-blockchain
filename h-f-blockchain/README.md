### Prepare
```
sudo apt-get install apt-transport-https ca-certificates curl software-properties-common vim libltdl-dev python make node-gyp -y
```

### Install Golang (v1.11.1)
```
# download go SDK
wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz -P ~/Downloads/
sudo tar zxf ~/Downloads/go1.11.1.linux-amd64.tar.gz -C /usr/local/

# set Golang related environment variables
export PATH=$PATH:/usr/local/go/bin
export GOPATH=~/gopath

# create $GOPATH folders
mkdir ~/gopath
cd $GOPATH
mkdir src bin pkg
```

### Install Docker
```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install docker-ce -y
sudo gpasswd -a ${USER} docker
# relogin linux user required
```

### Install docker-compose (v1.22.0)
```
wget https://github.com/docker/compose/releases/download/1.22.0/docker-compose-`uname -s`-`uname -m` -P ~/Downloads
sudo cp ~/Downloads/docker-compose-Linux-x86_64 /usr/local/bin/docker-compose
```

### Install Node.JS SDK (v8.x)
```
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
source ~/.bashrc
nvm install 8
node -v
npm -v
```

### Download docker images
```
cd ~/h-f-blockchain
./pull-images.sh
```

#download fabric binaries
```
wget https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/linux-amd64-1.4.0/hyperledger-fabric-linux-amd64-1.4.0.tar.gz -P ~/Downloads/

#extract fabric binaries to `fabric-bin`
tar zxf ~/Downloads/hyperledger-fabric-linux-amd64-1.4.0.tar.gz -C ~/h-f-blockchain/fabric-bin
```

---

### Scenario 1: build network
```
#setup fabric-cli into $PATH
cd ~/h-f-blockchain
export PATH=$PATH:~/h-f-blockchain/fabric-bin/bin

#generate crypto certs and channel-artifacts
cd ~/h-f-blockchain/fabric-network/example.com
./generate.sh

#startup fabric-network 'example.com'
./startup.sh

#run fabric-cli to setup channel, chaincode, and invoke chaincode
./runscript.sh
```

### Scenario 2: access FABRIC-CA 'ca.org1.example.com' and chaincode on FABRIC-PEER 'peer0.org1.example.com'

```
echo '######## - Prepare(!!!required only at first time!!!) - ########'
#install required SDK modules
cd ~/h-f-blockchain/apps/example02
npm install --registry=https://registry.npm.taobao.org


echo ###################### - BEGIN - ######################'
#enable tls feature
export TLS_ENABLED=true

#clear local cert wallet
rm -r hfc-key-store

echo '######## - ORG1 - ########'
export ORG_NAME=org1.example.com
export MSP_ID=Org1MSP
export CA_ADDRESS=localhost:7054
export CA_NAME=ca.org1.example.com
export PEER_ADDRESS=localhost:7051
export PEER_NAME=peer0.org1.example.com
export PEER_TLS_CERT="./tls/peer0.org1/ca.crt"

#enroll 'admin' with password 
echo '[ORG1]enroll `admin` and store certs in `hfc-key-store/org1.example.com`'
node enroll_admin_hospital.js

#register user 'user1'
echo '[ORG1]register `user1` and store certs in `hfc-key-store/org1.example.com`'
node register-user1_hospital.js

#add patient message
node addPatientMesg_hospital.js

#query medicine message
node queryMedicineMesg_hospital.js

#query transfer message
node queryTransferMesg_hospital.js



echo '######## - ORG2 - ########'
export ORG_NAME=org2.example.com
export MSP_ID=Org2MSP
export CA_ADDRESS=localhost:6054
export CA_NAME=ca.org2.example.com
export PEER_ADDRESS=localhost:6051
export PEER_NAME=peer0.org2.example.com
export PEER_TLS_CERT="./tls/peer0.org2/ca.crt"

#enroll 'admin' with password 
echo '[ORG2]enroll `admin` and store certs in `hfc-key-store/org2.example.com`'
node enroll_admin_factory.js

#register user 'user2'
echo '[ORG2]register `user2` and store certs in `hfc-key-store/org2.example.com`'
node register-user2_factory.js

#add medicine message
node addMedicineMesg_factory.js

#query patient message
node queryPatientMesg_factory.js

#add transfer message
node addTransferMesg_factory.js



echo '###################### - END - ######################'
```

### Finally: stop network 'example.com'

```
#shutdown fabric-network 'example.com' and teardown
cd ~/h-f-blockchain/fabric-network/example.com
./teardown.sh
```
