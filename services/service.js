'use strict';

const config = require('../blockchain/config');
var RealEstate = require('../models/realestate');

const Service = class Service {
    constructor(client) {
        this._client = client;
    }
    async getShareHolder(data) {
        try {
            let successResult = await this.invoke('get_shareholder', data);
            return successResult;
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async getRealEstate(id) {
        try {
            return await RealEstate.findOne({ id: id })
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async deleteRealEstate(id) {
        try {
            return await RealEstate.deleteOne({ id });
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listPublishRealEstate() {
        try {
            return await RealEstate.find({ actice: "Publish" });
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async getPublisher(data) {
        try {
            let successResult = await this.invoke('get_publisher', data);
            return successResult;
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async login(data) {
        try {
            await this._client.loginUser(data);
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listTransferContractByBuyer(data) {
        try {
            return await this.invoke('list_transfer_contract_by_buyer', data);
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listTransferContractBySeller(data) {
        try {
            return await this.invoke('list_transfer_contract_by_seller', data);
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listPublishContract(data) {
        try {
            return await this.invoke('list_publish_contract');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listShareHolder(data) {
        try {
            return await this.invoke('list_shareholder');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listPublisher(data) {
        try {
            return await this.invoke('list_publisher');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listTransferContract(data) {
        try {
            return await this.invoke('list_transfer_contract');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listCode(data) {
        try {
            return await this.invoke('list_real_estate');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async listSellTritAdvertisingContract(data) {
        try {
            return await this.invoke('list_sell_trit_advertising_contract');
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    invoke(fcn, ...args) {
        return this._client.invoke(
            config.chaincodeId, config.chaincodeVersion, fcn, ...args);
    }
}

module.exports.Service = Service