var mongoose = require('mongoose');

var Schema = mongoose.Schema;

var RealEstateSchema = new Schema(
  {
    id: { type: String, required: true, unique: true },
    name: { type: String, required: true },
    price: { type: Number, required: true },
    squareMeter: { type: Number, required: true },
    address: { type: String, required: true },
    description: { type: String, required: true },
    ownerId: { type: String, required: true },
    images: { type: [Buffer], required: true },
    amount: { type: Number, required: true },
    actice: { type: String, required: true, enum: ["New", "Publish", "Reject"], default: "New" }
  }
);

//Export model
module.exports = mongoose.model('RealEstate', RealEstateSchema);