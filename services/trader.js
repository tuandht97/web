'use strict';

const { Service } = require('./service');
const uuidV4 = require('uuid/v4');

const TraderService = class TraderService extends Service {
    async createShareHolder(data) {
        try {
            await this._client.getOrgAdmin();
            const successResult = await this.invoke('create_shareholder', data);
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

    async createAdvertising(data) {
        try {
            data = Object.assign({ uuid: uuidV4() }, data);
            const successResult = await this.invoke('create_sell_trit_advertising_contract', data);
            if (successResult) {
                throw (new Error(successResult))
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async changePriceAdvertising(data) {
        try {
            const successResult = await this.invoke('change_price_sell_trit_advertising_contract', data);
            if (successResult) {
                throw (new Error(successResult))
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async createTransferContractForBuyer(data) {
        try {
            data = Object.assign({ uuid: uuidV4() }, data);
            const successResult = await this.invoke('create_transfer_contract_for_buyer', data);
            if (successResult) {
                throw new Error(successResult);
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async payInByShareHolder(data) {
        try {
            const successResult = await this.invoke('pay_in_by_shareholder', data);
            if (successResult) {
                throw new Error(successResult);
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }

    async deleteAdvertising(data) {
        try {
            const successResult = await this.invoke('delete_sell_trit_advertising_contract', data);
            if (successResult) {
                throw (new Error(successResult))
            }
        } catch (e) {
            throw (`${e.message}`);
        }
    }
}

module.exports.TraderService = TraderService
