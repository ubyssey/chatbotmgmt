"use strict";

/*
 * Environment configuration file for Chatbot Conversation API.
 */

const oEnvironments = {
	production: {
		port: 80,
		address: "0.0.0.0"
	},
	development: {
		port: 3000,
		address: "127.0.0.1"
	}
}

module.exports = oEnvironments;
