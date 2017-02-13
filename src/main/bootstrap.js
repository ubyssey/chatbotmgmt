"use strict";

/**
 * Bootstraps the full chatbot, including the conversation and management API, with default config.
 */

const path = require("path");

const sConversationApiAppFile = path.resolve(__dirname, "conversation-api", "app.js");
const sConversationApiEnvironmentsFile = path.resolve(__dirname, "conversation-api", "config", "environments.js");
const oConversationApi = require(sConversationApiAppFile);
const oConversationApiConfig = require(sConversationApiEnvironmentsFile)[process.env.NODE_ENV];

if (oConversationApiConfig === undefined || oConversationApiConfig.port === undefined || oConversationApiConfig.address === undefined) {
    throw "[FATAL] Required conversation API configuration not found! Check config/environments.js and NODE_ENV environment variable.";
}

const sManagementApiAppFile = path.resolve(__dirname, "management-api", "app.js");
const sManagementApiEnvironmentsFile = path.resolve(__dirname, "management-api", "config", "environments.js");
const oManagementApi = require(sManagementApiAppFile);
const oManagementApiConfig = require(sManagementApiEnvironmentsFile)[process.env.NODE_ENV];

if (oManagementApiConfig === undefined || oManagementApiConfig.port === undefined || oManagementApiConfig.address === undefined) {
    throw "[FATAL] Required management API configuration not found! Check config/environments.js and NODE_ENV environment variable.";
}

oConversationApi.listen(oConversationApiConfig.port, oConversationApiConfig.address);
oManagementApi.listen(oManagementApiConfig.port, oManagementApiConfig.address);

console.log("[INFO] In " + process.env.NODE_ENV + " mode");
console.log("[INFO] Conversation API listening on port " + oConversationApiConfig.port + " bound to address " + oConversationApiConfig.address);
console.log("[INFO] Management API listening on port " + oManagementApiConfig.port + " bound to address " + oManagementApiConfig.address);
