'use strict';

const { TraderService } = require('./trader');
var RealEstate = require('../models/realestate');

const RealEstateService = class RealEstateService extends TraderService {
    async createRealEstate(data) {
        try {
            const realestate = new RealEstate(data);
            await realestate.save()
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listRealEstateByOwner(ownerId) {
        try {
            return await RealEstate.find({ ownerId });
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async editRealEstate(data) {
        try {
            let doc = await RealEstate.findOne({ id: data.id });
            if (!doc) throw "Not found this real estate";
            if (doc.actice !== "New") throw "Can't change this real estate";
            await RealEstate.updateOne({ id: data.id }, data)
        } catch (e) {
            throw (`${e.message}`)
        }
    }
}

module.exports.RealEstateService = RealEstateService