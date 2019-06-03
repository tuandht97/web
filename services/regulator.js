'use strict';

const { Service } = require('./service');
const uuidV4 = require('uuid/v4');
var RealEstate = require('../models/realestate');

const RegulatorService = class RegulatorService extends Service {

    async listNewRealEstate() {
        try {
            return await RealEstate.find({ "actice": "New" });
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async publishRealEstate(id) {
        try {
            let doc = await RealEstate.findOne({ id });
            if (!doc) throw new Error("This real estate hasn't been created");
            if (doc.actice != "New") throw new Error("Can't change this real estate");
            doc.actice = "Publish";
            let data = {
                id: doc.id,
                price: doc.price,
                squareMeter: doc.squareMeter,
                address: doc.address,
                amount: doc.amount,
                ownerId: doc.ownerId,
                description: doc.description
            };
            await doc.save();

            data = Object.assign({ uuid: uuidV4() }, data);
            const successResult = await this.invoke('create_and_publish_real_estate', data);
            if (successResult) {
                throw new Error(successResult);
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async rejectRealEstate(id) {
        try {
            let doc = await RealEstate.findOne({ id });
            if (!doc) throw "Not found this real estate";
            if (doc.actice !== "New") throw "Can't change this real estate";
            doc.actice = "Reject";
            await doc.save();
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async createPublisher(data) {
        try {
            await this._client.getOrgAdmin();
            const successResult = await this.invoke('create_publisher', data);
            if (successResult) {
                throw new Error(successResult);
            }
            else {
                return await this._client.register(data.username);
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async payIn(data) {
        try {
            const successResult = await this.invoke('pay_in', data);
            if (successResult) {
                throw (new Error(successResult))
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async logHistory(data) {
        try {
            let result = await this.invoke('get_history_for_ufo', data);
            // result.map(async (element) => {
            //     let tx_id = element.TxId;
            //     let ts = await this._client.queryTransaction(tx_id);
            // });
            return result;
        } catch (e) {
            throw (`${e.message}`);
        }
    }

}

module.exports.RegulatorService = RegulatorService
