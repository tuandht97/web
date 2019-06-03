'use strict';
const express = require('express')
const router = express.Router();
var { regulatorService } = require('../services/total');
var jwt = require('jsonwebtoken');
var path = require('path');
var fs = require('fs')


//Authentication
router.all('*', async (req, res, next) => {
  // check header or url parameters or post parameters for token
  var token = req.headers.authorization;

  // decode token
  if (token) {

    // verifies secret and checks exp
    jwt.verify(token, process.env.SECRET, async (err, decoded) => {
      if (err) {
        return res.json({ error: 'Failed to authenticate token.' });
      } else {
        // if everything is good, save to request for use in other routes
        req.decoded = decoded;
        if (req.decoded.mspid && req.decoded.mspid == "RegulatorOrgMSP") {
          try {
            await regulatorService.login(req.decoded.name);
          } catch (error) {
            return res.json({ error })
          }
          next();
        } else {
          return res.json({ error: "You don't belong this Org" })
        }
      }
    });

  } else {

    // if there is no token
    // return an error
    return res.status(403).send({
      error: 'No token provided.'
    });

  }

});

//regulator
router.get('/list-new-real-estate', async (req, res, next) => {
  try {
    let result = await regulatorService.listNewRealEstate();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/create-publish-contract', async (req, res, next) => {
  let {
    id
  } = req.body;
  if (typeof id !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await regulatorService.publishRealEstate(id);
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/reject-real-estate', async (req, res, next) => {
  let {
    id
  } = req.body;
  if (typeof id !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await regulatorService.rejectRealEstate(id);
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/payin', async (req, res, next) => {
  let {
    amount,
    receiver
  } = req.body;
  if (typeof amount !== 'number'
    || typeof receiver !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await regulatorService.payIn({
      amount,
      receiver
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error })
  }
});

router.post('/create-publisher', async (req, res, next) => {
  let {
    username,
    firstName,
    lastName,
    identityCard
  } = req.body;
  if (typeof username !== 'string'
    || typeof firstName !== 'string'
    || typeof lastName !== 'string'
    || typeof identityCard !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await regulatorService.createPublisher({
      username,
      firstName,
      lastName,
      identityCard
    });
    const url = path.join(__dirname, '../regulator-peer', username);
    res.download(url, function (err) {
      if (err) {
        throw err;
      } else {
        fs.unlinkSync(url);
      }
    });
    return
  } catch (error) {
    res.json({ error });
  }
});

router.post('/log-history', async (req, res, next) => {
  let {
    assetId,
    assetType
  } = req.body;
  if (typeof assetId !== 'string'
    || typeof assetType !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    let result = await regulatorService.logHistory({
      assetId,
      assetType
    });
    res.json({ result })
  } catch (error) {
    res.json({ error });
  }
});

//commom
router.post('/get-shareholder', async (req, res, next) => {
  try {
    let {
      username
    } = req.body;
    if (typeof username !== 'string') {
      res.json({ error: 'Invalid request.' });
      return;
    }
    let result = await regulatorService.getShareHolder({
      username
    });
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/get-publisher', async (req, res, next) => {
  try {
    let {
      username
    } = req.body;
    if (typeof username !== 'string') {
      res.json({ error: 'Invalid request.' });
      return;
    }
    let result = await regulatorService.getPublisher({
      username
    });
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/list-transfer-contract-by-buyer', async (req, res) => {
  try {
    let {
      username
    } = req.body;
    if (typeof username !== 'string') {
      res.json({ error: 'Invalid request.' });
      return;
    }
    const result = await regulatorService.listTransferContractByBuyer({
      username
    });
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.post('/list-transfer-contract-by-seller', async (req, res) => {
  try {
    let {
      username
    } = req.body;
    if (typeof username !== 'string') {
      res.json({ error: 'Invalid request.' });
      return;
    }
    const result = await regulatorService.listTransferContractBySeller({
      username
    });
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-publish-contract', async (req, res) => {
  try {
    const result = await regulatorService.listPublishContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-shareholder', async (req, res) => {
  try {
    const result = await regulatorService.listShareHolder();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-publisher', async (req, res) => {
  try {
    const result = await regulatorService.listPublisher();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-transfer-contract', async (req, res) => {
  try {
    const result = await regulatorService.listTransferContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-code', async (req, res) => {
  try {
    const result = await regulatorService.listCode();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-advertising', async (req, res) => {
  try {
    const result = await regulatorService.listSellTritAdvertisingContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/get-real-estate/:id', async (req, res) => {
  try {
    let id = req.params.id;
    const result = await regulatorService.getRealEstate(id);
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-publish-real-estate', async (req, res) => {
  try {
    const result = await regulatorService.listPublishRealEstate();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

module.exports = router