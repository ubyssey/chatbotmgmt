"use strict";

/*
 * Route definitions for the Chatbot Management API.
 *
 * This module exports an object with a single function, configureRoutes,
 * which accepts an express application, and configures routes on it.
 * When the configureRoutes method returns, all publicly routable controllers
 * should be mounted.
 */

const oSessions = require("./../controllers/sessions");

function configureRoutes(oApp) {
	// Mount web-facing controllers
	oApp.use("/sessions", oSessions);
}

module.exports.configureRoutes = configureRoutes;