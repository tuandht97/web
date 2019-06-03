'use strict';
const config = require('./config');
const { OrganizationClient } = require('./utils');
const http = require('http');
const url = require('url');

let status = 'down';
let statusChangedCallbacks = [];
const { writeFileSync, unlinkSync } = require('fs');
const writeCryptoFile =
  (fileName, data) => writeFileSync(fileName, data);
// Setup clients per organization
const realestateClient = new OrganizationClient(
  config.channelName,
  config.orderer0,
  config.realestateOrg.peer,
  config.realestateOrg.ca,
  config.realestateOrg.admin
);
const traderClient = new OrganizationClient(
  config.channelName,
  config.orderer0,
  config.traderOrg.peer,
  config.traderOrg.ca,
  config.traderOrg.admin
);
const regulatorClient = new OrganizationClient(
  config.channelName,
  config.orderer0,
  config.regulatorOrg.peer,
  config.regulatorOrg.ca,
  config.regulatorOrg.admin
);

function setStatus(s) {
  status = s;

  setTimeout(() => {
    statusChangedCallbacks
      .filter(f => typeof f === 'function')
      .forEach(f => f(s));
  }, 1000);
}

module.exports.subscribeStatus = function subscribeStatus(cb) {
  if (typeof cb === 'function') {
    statusChangedCallbacks.push(cb);
  }
};

module.exports.getStatus = function getStatus() {
  return status;
};

module.exports.isReady = function isReady() {
  return status === 'ready';
};

function getAdminOrgs() {
  return Promise.all([
    realestateClient.getOrgAdmin(),
    traderClient.getOrgAdmin(),
    regulatorClient.getOrgAdmin()
  ]);
}

(async () => {
  // Login
  try {
    await Promise.all([
      realestateClient.login(),
      traderClient.login(),
      regulatorClient.login()
    ]);
  } catch (e) {
    console.log('Fatal error logging into blockchain organization clients!');
    console.log(e);
    process.exit(-1);
  }

  // Bootstrap blockchain network
  try {
    await getAdminOrgs();
    if (!(await realestateClient.checkChannelMembership())) {
      console.log('Default channel not found, attempting creation...');
      const createChannelResponse =
        await realestateClient.createChannel(config.channelConfig);
      if (createChannelResponse.status === 'SUCCESS') {
        console.log('Successfully created a new default channel.');
        console.log('Joining peers to the default channel.');
        await Promise.all([
          realestateClient.joinChannel(),
          traderClient.joinChannel(),
          regulatorClient.joinChannel()
        ]);
        // Wait for 10s for the peers to join the newly created channel
        await new Promise(resolve => {
          setTimeout(resolve, 10000);
        });
      }
    }
  } catch (e) {
    console.log('Fatal error bootstrapping the blockchain network!');
    console.log(e);
    process.exit(-1);
  }

  // Register block events
  try {
    console.log('Connecting and Registering Block Events');
    realestateClient.connectAndRegisterBlockEvent();
    traderClient.connectAndRegisterBlockEvent();
    regulatorClient.connectAndRegisterBlockEvent();
  } catch (e) {
    console.log('Fatal error register block event!');
    console.log(e);
    process.exit(-1);
  }

  // Initialize network
  // try {
  //   await Promise.all([
  //     realestateClient.initialize(),
  //     traderClient.initialize(),
  //     regulatorClient.initialize()
  //   ]);
  // } catch (e) {
  //   console.log('Fatal error initializing blockchain organization clients!');
  //   console.log(e);
  //   process.exit(-1);
  // }

  // Install chaincode on all peers
  // let installedOnRealEstateOrg, installedOnTraderOrg,
  //   installedOnRegulatorOrg;
  // try {
  //   await getAdminOrgs();
  //   installedOnRealEstateOrg = await realestateClient.checkInstalled(
  //     config.chaincodeId, config.chaincodeVersion, config.chaincodePath);
  //   installedOnTraderOrg = await traderClient.checkInstalled(
  //     config.chaincodeId, config.chaincodeVersion, config.chaincodePath);
  //   installedOnRegulatorOrg = await regulatorClient.checkInstalled(
  //     config.chaincodeId, config.chaincodeVersion, config.chaincodePath);
  // } catch (e) {
  //   console.log('Fatal error getting installation status of the chaincode!');
  //   console.log(e);
  //   process.exit(-1);
  // }

  // if (!(installedOnRealEstateOrg &&
  //   installedOnTraderOrg && installedOnRegulatorOrg)) {
  //   console.log('Chaincode is not installed, attempting installation...');

  //   // Pull chaincode environment base image
  //   try {
  //     await getAdminOrgs();
  //     const socketPath = process.env.DOCKER_SOCKET_PATH ||
  //       (process.platform === 'win32' ? '//./pipe/docker_engine' : '/var/run/docker.sock');
  //     const ccenvImage = process.env.DOCKER_CCENV_IMAGE ||
  //       'hyperledger/fabric-ccenv:latest';
  //     const listOpts = { socketPath, method: 'GET', path: '/images/json' };
  //     const pullOpts = {
  //       socketPath, method: 'POST',
  //       path: url.format({ pathname: '/images/create', query: { fromImage: ccenvImage } })
  //     };

  //     const images = await new Promise((resolve, reject) => {
  //       const req = http.request(listOpts, (response) => {
  //         let data = '';
  //         response.setEncoding('utf-8');
  //         response.on('data', chunk => { data += chunk; });
  //         response.on('end', () => { resolve(JSON.parse(data)); });
  //       });
  //       req.on('error', reject); req.end();
  //     });

  //     const imageExists = images.some(
  //       i => i.RepoTags && i.RepoTags.some(tag => tag === ccenvImage));
  //     if (!imageExists) {
  //       console.log(
  //         'Base container image not present, pulling from Docker Hub...');
  //       await new Promise((resolve, reject) => {
  //         const req = http.request(pullOpts, (response) => {
  //           response.on('data', () => { });
  //           response.on('end', () => { resolve(); });
  //         });
  //         req.on('error', reject); req.end();
  //       });
  //       console.log('Base container image downloaded.');
  //     } else {
  //       console.log('Base container image present.');
  //     }
  //   } catch (e) {
  //     console.log('Fatal error pulling docker images.');
  //     console.log(e);
  //     process.exit(-1);
  //   }

  //   // Install chaincode
  //   const installationPromises = [
  //     realestateClient.install(
  //       config.chaincodeId, config.chaincodeVersion, config.chaincodePath),
  //     traderClient.install(
  //       config.chaincodeId, config.chaincodeVersion, config.chaincodePath),
  //     regulatorClient.install(
  //       config.chaincodeId, config.chaincodeVersion, config.chaincodePath)
  //   ];
  //   try {
  //     await Promise.all(installationPromises);
  //     await new Promise(resolve => { setTimeout(resolve, 10000); });
  //     console.log('Successfully installed chaincode on the default channel.');
  //   } catch (e) {
  //     console.log('Fatal error installing chaincode on the default channel!');
  //     console.log(e);
  //     process.exit(-1);
  //   }

  //   // Instantiate chaincode on all peers
  //   // Instantiating the chaincode on a single peer should be enough (for now)
  //   try {
  //     // Initial contract types
  //     await regulatorClient.instantiate(config.chaincodeId,
  //       config.chaincodeVersion, []);
  //     console.log('Successfully instantiated chaincode on all peers.');
  //     setStatus('ready');
  //   } catch (e) {
  //     console.log('Fatal error instantiating chaincode on some(all) peers!');
  //     console.log(e);
  //     process.exit(-1);
  //   }
  // } else {
  //   console.log('Chaincode already installed on the blockchain network.');
  //   setStatus('ready');
  // }

  // init regulator 
  // try {
  //   let user = await regulatorClient._client.getUserContext("regulator", true);
  //   if (user && user.isEnrolled()) {
  //     return user;
  //   }
  //   const data = {
  //     username: "regulator",
  //     firstName: "sjfh",
  //     lastName: "snfd",
  //     identityCard: "smfds"
  //   }
  //   const successResult = await regulatorClient.invoke(config.chaincodeId, config.chaincodeVersion, 'create_publisher', data);
  //   if (successResult) {
  //     throw new Error(successResult);
  //   }
  //   else {
  //     return await regulatorClient.register(data.username);
  //   }
  // } catch (e) {
  //   console.log(`${e.message}`);
  // }
})();

// Export organization clients
module.exports = {
  realestateClient,
  traderClient,
  regulatorClient
};
