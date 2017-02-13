"use strict";

/*
 * Configure and prepare the Conversation API app for use.
 */

const express = require("express");
const oApp = express();

// load middleware
const bodyParser = require("body-parser");
oApp.use(bodyParser.json());

// load routes
const oRoutes = require("./config/routes");
oRoutes.configureRoutes(oApp);

module.exports = oApp;
