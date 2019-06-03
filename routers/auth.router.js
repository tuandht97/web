'use strict';
const express = require('express')
const router = express.Router()
var fs = require('fs')
var jwt = require('jsonwebtoken');
process.env.SECRET = "my secret"
var path = require('path');
const { regulatorService, traderService, realEstateService } = require('../services/total')

router.post('/register', async (req, res, next) => {
    let {
        username,
        firstName,
        lastName,
        identityCard,
        org
    } = req.body;
    if (typeof username !== 'string'
        || typeof firstName !== 'string'
        || typeof lastName !== 'string'
        || typeof identityCard !== 'string'
        || typeof org !== 'string') {
        res.json({ error: 'Invalid request.' });
        return;
    }

    try {
        if (org == "RealEstate") {
            await realEstateService.createShareHolder({
                username,
                firstName,
                lastName,
                identityCard
            })
            const url = path.join(__dirname, '../realestate-peer', username);
            res.download(url, function (err) {
                if (err) {
                    throw err;
                } else {
                    fs.unlinkSync(url);
                }
            });
            return
        } else if (org == "Trader") {
            await traderService.createShareHolder({
                username,
                firstName,
                lastName,
                identityCard
            })
            const url = path.join(__dirname, '../trader-peer', username);
            res.download(url, function (err) {
                if (err) {
                    throw err;
                } else {
                    fs.unlinkSync(url);
                }
            });
            return
        }  else if (org == "Regulator") {
            await regulatorService.createPublisher({
                username,
                firstName,
                lastName,
                identityCard
            })
            const url = path.join(__dirname, '../regulator-peer', username);
            res.download(url, function (err) {
                if (err) {
                    throw err;
                } else {
                    fs.unlinkSync(url);
                }
            });
            return
        } else {
            return res.json({ error: "Incorrect Org" });
        }
    } catch (error) {
        res.json({ error });
    }
})

router.post('/login', async (req, res, next) => {
    try {
        if (req.files) {
            const data = req.files.key.data.toString('utf8');
            const item = JSON.parse(data)
            const name = item.name;
            const mspid = item.mspid;
            const signingIdentity = item.signingIdentity;
            if (mspid == "RealEstateOrgMSP") {
                const url = path.join(__dirname, '../realestate-peer', name);
                req.files.key.mv(url)
                await realEstateService.login(name)
            } else if (mspid == "TraderOrgMSP") {
                const url = path.join(__dirname, '../trader-peer', name);
                req.files.key.mv(url)
                await traderService.login(name)
            } else if (mspid == "RegulatorOrgMSP") {
                const url = path.join(__dirname, '../regulator-peer', name);
                req.files.key.mv(url)
                await regulatorService.login(name)
            } else {
                throw "Incorrect Org"
            }
            var payload = {
                name,
                mspid,
                signingIdentity
            }
            var token = jwt.sign(payload, process.env.SECRET, {
                expiresIn: 86400 // expires in 24 hours
            });

            res.json({
                token: token
            });
        }
        else {
            throw "File not upload";
        }

    } catch (error) {
        res.json({ error })
    }
});

router.post('/logout', async (req, res, next) => {
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
            }
        });

    } else {
        // if there is no token
        // return an error
        return res.status(403).send({
            error: 'No token provided.'
        });

    }


    try {
        let name = req.decoded.name
        let mspid = req.decoded.mspid;
        if (mspid == "RealEstateOrgMSP") {
            const url = path.join(__dirname, '../realestate-peer', name);
            fs.unlinkSync(url);
        } else if (mspid == "TraderOrgMSP") {
            const url = path.join(__dirname, '../trader-peer', name);
            fs.unlinkSync(url);
        } else if (mspid == "RegulatorOrgMSP") {
            const url = path.join(__dirname, '../regulator-peer', name);
            fs.unlinkSync(url);
        }
        res.json({ result: "success" })
    } catch (error) {
        res.json({ error })
    }
});

router.get('/exchange', async (req, res) => {
    try {
      const result = await regulatorService.listSellTritAdvertisingContract();
      res.json({ result });
    } catch (error) {
      res.json({ error });
    }
  });


module.exports = router