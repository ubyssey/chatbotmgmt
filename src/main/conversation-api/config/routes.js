"use strict";

/*
 * Route definitions for the Chatbot Conversation API.
 *
 * This module exports an object with a single function, configureRoutes,
 * which accepts an express application, and configures routes on it.
 * When the configureRoutes method returns, all publicly routable controllers
 * should be mounted.
 */

const oWebhook = require("./../controllers/webhook");
const oMetrics = require("./../controllers/metrics");

function configureRoutes(oApp) {
	// Mount web-facing controllers
	oApp.use("/webhook", oWebhook);
	oApp.use("/metrics", oMetrics);
}

module.exports.configureRoutes = configureRoutes;