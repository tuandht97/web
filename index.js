var express = require('express')
const fileUpload = require('express-fileupload');
const morgan = require('morgan');
process.env.SECRET = "my secret"
var bodyParser = require('body-parser')
var port = 3000
var app = express()
var realestate = require('./routers/realestate.router')
var regulator = require('./routers/regulator.router')
var trader = require('./routers/trader.router')
var auth = require('./routers/auth.router')

//Set up mongoose connection
var mongoose = require('mongoose');
var mongoDB = "mongodb://localhost:27017/estate";
mongoose.connect(mongoDB, { useNewUrlParser: true });
var db = mongoose.connection;
db.on('error', console.error.bind(console, 'MongoDB connection error:'));

if (process.env.NODE_ENV == "test") {
    // use morgan to log at command line
    app.use(morgan('combined')); //'combined' outputs the Apache style LOGs
}

app.use(fileUpload());
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({
    extended: false
}));

app.use('/uploads', express.static('uploads'));
app.use('/api/realestate', realestate)
app.use('/api/regulator', regulator)
app.use('/api/trader', trader)
app.use('/api/auth', auth)

module.exports = app;
app.listen(port, () => console.log(`Example app listening on port ${port}!`))