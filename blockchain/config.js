const { readFileSync } = require('fs');
const { resolve } = require('path');

const basePath = resolve(__dirname, '../certs');
const readCryptoFile =
  filename => readFileSync(resolve(basePath, filename)).toString();
const config = {
  channelName: 'default',
  channelConfig: readFileSync(resolve(__dirname, '../channel.tx')),
  chaincodeId: 'trit_chaincode',
  chaincodeVersion: 'v2',
  chaincodePath: 'trit_chaincode',
  orderer0: {
    hostname: 'orderer0',
    url: 'grpcs://34.87.116.245:7050',
    pem: readCryptoFile('ordererOrg.pem')
  },
  realestateOrg: {
    peer: {
      hostname: 'realestate-peer',
      url: 'grpcs://34.87.84.124:7051',
      eventHubUrl: 'grpcs://34.87.84.124:7053',
      pem: readCryptoFile('realestateOrg.pem')
    },
    ca: {
      name: 'realestate-org',
      hostname: 'realestate-ca',
      url: 'grpcs://34.87.84.124:7054',
      mspId: 'RealEstateOrgMSP'
    },
    admin: {
      key: readCryptoFile('Admin@realestate-org-key.pem'),
      cert: readCryptoFile('Admin@realestate-org-cert.pem')
    }
  },
  regulatorOrg: {
    peer: {
      hostname: 'regulator-peer',
      url: 'grpcs://34.87.89.73:10051',
      eventHubUrl: 'grpcs://34.87.89.73:10053',
      pem: readCryptoFile('regulatorOrg.pem')
    },
    ca: {
      name: 'regulator-org',
      hostname: 'regulator-ca',
      url: 'grpcs://34.87.89.73:10054',
      mspId: 'RegulatorOrgMSP'
    },
    admin: {
      key: readCryptoFile('Admin@regulator-org-key.pem'),
      cert: readCryptoFile('Admin@regulator-org-cert.pem')
    }
  },
  traderOrg: {
    peer: {
      hostname: 'trader-peer',
      url: 'grpcs://35.247.134.54:9051',
      pem: readCryptoFile('traderOrg.pem'),
      eventHubUrl: 'grpcs://35.247.134.54:9053',
    },
    ca: {
      name: 'trader-org',
      hostname: 'trader-ca',
      url: 'grpcs://35.247.134.54:9054',
      mspId: 'TraderOrgMSP'
    },
    admin: {
      key: readCryptoFile('Admin@trader-org-key.pem'),
      cert: readCryptoFile('Admin@trader-org-cert.pem')
    }
  }
};

// if (true) {
//   config.orderer0.url = 'grpcs://localhost:7050';

//   config.realestateOrg.peer.url = 'grpcs://localhost:7051';
//   config.traderOrg.peer.url = 'grpcs://localhost:9051';
//   config.regulatorOrg.peer.url = 'grpcs://localhost:10051';

//   config.realestateOrg.peer.eventHubUrl = 'grpcs://localhost:7053';
//   config.traderOrg.peer.eventHubUrl = 'grpcs://localhost:9053';
//   config.regulatorOrg.peer.eventHubUrl = 'grpcs://localhost:10053';

//   config.realestateOrg.ca.url = 'grpcs://localhost:7054';
//   config.traderOrg.ca.url = 'grpcs://localhost:9054';
//   config.regulatorOrg.ca.url = 'grpcs://localhost:10054';
// }

module.exports = config;