"use strict";

/*
 * Environment configuration file for Chatbot Management API.
 */

const oEnvironments = {
	production: {
		port: 80,
		address: "0.0.0.0"
	},
	development: {
		port: 3001,
		address: "127.0.0.1"
	}
}

module.exports = oEnvironments;
