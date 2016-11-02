"use strict";

/**
 * This function won't cause error by define non-ws protocol.
 * Instead, it overrides value with ws protocol.
 *
 * @param uri: string
 * @return string
 */
function setWsProtocol (uri) {
  var h = uri.split("://");
  var protocol = (h[1])? h[0]: "ws";
  var hostname = (h[1])? h[1]: h[0];
  if (protocol.indexOf("ws") < 0) {
    protocol = (protocol === "https")? "wss": "ws";
  }

  return protocol + "://" + hostname;
}

module.exports = {
  setWsProtocol: setWsProtocol,
}
