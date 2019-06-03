'use strict';
const express = require('express')
const router = express.Router();
var path = require('path');
const { realEstateService } = require('../services/total')
var jwt = require('jsonwebtoken');

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
        if (req.decoded.mspid && req.decoded.mspid == "RealEstateOrgMSP") {
          try {
            await realEstateService.login(req.decoded.name);
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

// real estate
router.post('/create-real-estate', async (req, res, next) => {
  let {
    id,
    name,
    price,
    squareMeter,
    address,
    amount,
    description
  } = req.body;
  if (typeof id !== 'string'
    // || typeof price !== 'number'
    // || typeof squareMeter !== 'number'
    || typeof address !== 'string'
    || typeof description !== 'string'
    // || typeof amount !== 'number'
  ) {
    res.json({ error: 'Invalid request.' });
    return;
  }

  let images = [];

  if (req.files) {
    for (var i = 0; i < req.files.images.length; i++) {
      var file = req.files.images[i];
      if (file.mimetype != 'image/png' && file.mimetype != 'image/jpeg' && file.mimetype != 'image/jpg') {
        res.json({ error: 'Images incorrect format.' });
        return;
      }
      images.push(file.data);
    }
  }

  try {
    await realEstateService.createRealEstate({
      id,
      name,
      price: Number(price),
      squareMeter: Number(squareMeter),
      address,
      amount: Number(amount),
      ownerId: req.decoded.name,
      images,
      description
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-own-real-estate', async (req, res, next) => {
  let ownerId = req.decoded.name;
  try {
    let result = await realEstateService.listRealEstateByOwner(ownerId);
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.put('/edit-real-estate', async (req, res) => {
  let {
    id,
    name,
    price,
    squareMeter,
    address,
    amount,
    description
  } = req.body;
  if (typeof id !== 'string'
    // || typeof price !== 'number'
    // || typeof squareMeter !== 'number'
    || typeof address !== 'string'
    || typeof description !== 'string'
    // || typeof amount !== 'number'
  ) {
    res.json({ error: 'Invalid request.' });
    return;
  }

  let images = [];
  if (req.files) {
    for (var i = 0; i < req.files.images.length; i++) {
      var file = req.files.images[i];
      if (file.mimetype != 'image/png' && file.mimetype != 'image/jpeg' && file.mimetype != 'image/jpg') {
        res.json({ error: 'Images incorrect format.' });
        return;
      }
      images.push(file.data);
    }
  }

  try {
    await realEstateService.editRealEstate({
      id,
      name,
      price: Number(price),
      squareMeter: Number(squareMeter),
      address,
      amount: Number(amount),
      ownerId: req.decoded.name,
      images,
      description
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error });
  }
});

router.delete('/real-estate/:id', async (req, res) => {
  try {
    let id = req.params.id;
    await realEstateService.deleteRealEstate(id);
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error });
  }
});


// real estate and trader

router.post('/create-advertising', async (req, res, next) => {
  let {
    tritId,
    amount,
    price
  } = req.body;
  if (typeof tritId !== 'string'
    || typeof amount !== 'number'
    || typeof price !== 'number') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    let result = await realEstateService.createAdvertising({
      tritId,
      amount,
      price
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error })
  }
});
router.post('/change-price-advertising', async (req, res, next) => {
  let {
    price,
    sellTritAdvertisingContractId
  } = req.body;
  if (typeof price !== 'number'
    || typeof sellTritAdvertisingContractId !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await realEstateService.changePriceAdvertising({
      price,
      sellTritAdvertisingContractId
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error })
  }
});

router.post('/create-transfer-contract-for-buyer', async (req, res, next) => {
  let {
    amount,
    sellTritAdvertisingContractId
  } = req.body;
  if (typeof amount !== 'number'
    || typeof sellTritAdvertisingContractId !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    let result = await realEstateService.createTransferContractForBuyer({
      amount,
      sellTritAdvertisingContractId
    });
    res.json({ result });
  } catch (error) {
    res.json({ error })
  }
});

router.post('/payin-by-shareholder', async (req, res, next) => {
  let {
    amount
  } = req.body;
  if (typeof amount !== 'number') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await realEstateService.payInByShareHolder({
      amount
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error })
  }
});

router.post('/delete-advertising', async (req, res, next) => {
  let {
    id
  } = req.body;
  if (typeof id !== 'string') {
    res.json({ error: 'Invalid request.' });
    return;
  }

  try {
    await realEstateService.deleteAdvertising({
      id
    });
    res.json({ result: "success" });
  } catch (error) {
    res.json({ error })
  }
});

// common
router.post('/get-shareholder', async (req, res, next) => {
  try {
    let {
      username
    } = req.body;
    if (typeof username !== 'string') {
      res.json({ error: 'Invalid request.' });
      return;
    }
    let result = await realEstateService.getShareHolder({
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
    let result = await realEstateService.getPublisher({
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
    const result = await realEstateService.listTransferContractByBuyer({
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
    const result = await realEstateService.listTransferContractBySeller({
      username
    });
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});


router.get('/list-publish-contract', async (req, res) => {
  try {
    const result = await realEstateService.listPublishContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-shareholder', async (req, res) => {
  try {
    const result = await realEstateService.listShareHolder();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-publisher', async (req, res) => {
  try {
    const result = await realEstateService.listPublisher();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-transfer-contract', async (req, res) => {
  try {
    const result = await realEstateService.listTransferContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-code', async (req, res) => {
  try {
    const result = await realEstateService.listCode();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-advertising', async (req, res) => {
  try {
    const result = await realEstateService.listSellTritAdvertisingContract();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/get-real-estate/:id', async (req, res) => {
  try {
    let id = req.params.id;
    const result = await realEstateService.getRealEstate(id);
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

router.get('/list-publish-real-estate', async (req, res) => {
  try {
    const result = await realEstateService.listPublishRealEstate();
    res.json({ result });
  } catch (error) {
    res.json({ error });
  }
});

module.exports = router
