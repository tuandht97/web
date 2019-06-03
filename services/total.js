const { RegulatorService } = require('./regulator')
const { TraderService } = require('./trader')
const { RealEstateService } = require('./realestate')
const { regulatorClient, realestateClient, traderClient } = require('../blockchain/setup')

const regulatorService = new RegulatorService(regulatorClient);
const traderService = new TraderService(traderClient)
const realEstateService = new RealEstateService(realestateClient)

module.exports = {
    regulatorService,
    traderService,
    realEstateService
}