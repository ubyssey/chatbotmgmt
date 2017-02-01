"use strict";

/**
 * Bootstraps the full chatbot, including the conversation and management API, with default config.
 */

const spawn = require("child_process").spawn;
const path = require("path");

const sNodeExecutable = process.execPath;
const sProjectRoot = process.cwd();

const sConversationApiBootstrapFile = path.resolve(sProjectRoot, "chatbot-conversation-api", "bootstrap.js");
const oConversationApiProcess = spawn(sNodeExecutable, [sConversationApiBootstrapFile]);

oConversationApiProcess.stdout.on("data", (data) => {
	process.stdout.write("[CONV] " + data);
});

oConversationApiProcess.stderr.on("data", (data) => {
	process.stderr.write("[CONV] " + data);
});

const sManagementApiBootstrapFile = path.resolve(sProjectRoot, "chatbot-management-api", "bootstrap.js");
const oManagementApiProcess = spawn(sNodeExecutable, [sManagementApiBootstrapFile]);

oManagementApiProcess.stdout.on("data", (data) => {
	process.stdout.write("[MGMT] " + data);
});

oManagementApiProcess.stderr.on("data", (data) => {
	process.stderr.write("[MGMT] " + data);
});
