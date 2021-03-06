'use strict';

var Fabric_Client = require('fabric-client');
var path = require('path');
var util = require('util');
var os = require('os');
var fs = require('fs');

var orgName = "org1.example.com";
if (process.env.ORG_NAME) {
    orgName = process.env.ORG_NAME;
}

var peerAddress = "localhost:7051";
if (process.env.PEER_ADDRESS) {
    peerAddress = process.env.PEER_ADDRESS;
}
var peerName = "peer0.org1.example.com";
if (process.env.PEER_NAME) {
    peerName = process.env.PEER_NAME;
}

var peerTLSCertPath = "./tls/peer0.org1/ca.crt";
if (process.env.PEER_TLS_CERT) {
    peerTLSCertPath = process.env.PEER_TLS_CERT;
}

//
var fabric_client = new Fabric_Client();

// setup the fabric network
var channel = fabric_client.newChannel('mychannel');
var peer
if (process.env.TLS_ENABLED && process.env.TLS_ENABLED === 'true') {
	console.log("process.env.TLS_ENABLED is 'true'");
	let peerTLSCert = fs.readFileSync(path.join(__dirname, peerTLSCertPath));
	peer = fabric_client.newPeer('grpcs://' + peerAddress, {
		pem: Buffer.from(peerTLSCert).toString(),
		"ssl-target-name-override": peerName
	});
} else {
	console.log("process.env.TLS_ENABLED is 'false' or not existing.");
	peer = fabric_client.newPeer('grpc://' + peerAddress);
}
channel.addPeer(peer);

//
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store', orgName);
console.log('Store path:'+store_path);
var tx_id = null;

// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
Fabric_Client.newDefaultKeyValueStore({ path: store_path}).then((state_store) => {
	// assign the store to the fabric client
	fabric_client.setStateStore(state_store);
	var crypto_suite = Fabric_Client.newCryptoSuite();
	// use the same location for the state store (where the users' certificate are kept)
	// and the crypto store (where the users' keys are kept)
	var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
	crypto_suite.setCryptoKeyStore(crypto_store);
	fabric_client.setCryptoSuite(crypto_suite);

	// get the enrolled user from persistence, this user will sign all requests
	return fabric_client.getUserContext('user1', true);
}).then((user_from_store) => {
	if (user_from_store && user_from_store.isEnrolled()) {
		console.log('Successfully loaded `user1` from persistence');
		member_user = user_from_store;
	} else {
		throw new Error('Failed to get `user1`.... run registerUser.js');
	}

	// queryCar chaincode function - requires 1 argument, ex: args: ['CAR4'],
	// queryAllCars chaincode function - requires no arguments , ex: args: [''],
	const request = {
		//targets : --- letting this default to the peers assigned to the channel
		chaincodeId: 'mycc',
		fcn: 'queryTransferMesg',
		args: ['000002','000002']
	};

	// send the query proposal to the peer
	return channel.queryByChaincode(request);
}).then((query_responses) => {
	console.log("Query has completed, checking results");
	// query_responses could have more than one  results if there multiple peers were used as targets
	if (query_responses && query_responses.length == 1) {
		if (query_responses[0] instanceof Error) {
			console.error("error from query = ", query_responses[0]);
		} else {
			console.log("Response is ", query_responses[0].toString());
		}
	} else {
		console.log("No payloads were returned from query");
	}
}).catch((err) => {
	console.error('Failed to query :: ' + err);
});