"use strict";

/*
 * Bootstraps the Chatbot Conversation API.
 * 
 * The Conversation API is the component of the project that actually
 * handles communication with Facebook.
 * 
 * Looks for configuration in the following environment variables:
 * - NODE_ENV: the configuration to use when running the application.
 *				Valid values are "development" and "production".
 */

const oApp = require("./app");

// load environment
const oEnvironment = require("./config/environments")[process.env.NODE_ENV];

if (oEnvironment === undefined || oEnvironment.port === undefined || oEnvironment.address === undefined) {
	throw "[FATAL] Required configuration not found! Check config/environments.js and NODE_ENV environment variable.";
}

oApp.listen(oEnvironment.port, oEnvironment.address);
console.log("[INFO] In " + process.env.NODE_ENV + " mode");
console.log("[INFO] Listening on port " + oEnvironment.port + " bound to address " + oEnvironment.address);
