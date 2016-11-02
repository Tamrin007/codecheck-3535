"use strict";

const ws = require("ws");
const setWsProtocol = require("./utils").setWsProtocol;
const assert = require("chai").assert;
const env = require("../config/env.json");

let host_url = process.env.HOST_URL || `ws://${env.hostname}:${env.port}${env.path}`;
const HOST_URL = setWsProtocol(host_url);

describe("broadcast", () => {
  it ("should broadcast messages to all users", done => {
    let ua = new ws(HOST_URL);
    let ub = new ws(HOST_URL);
    let complete = () => {
      ua.close(); ub.close();
      done();
    };
    let message = JSON.stringify({ text: "Hi! I'm a user a!" });
    ua.on("open", () => ua.send(message));
    ub.on("open", () => true);

    let cnt = 0;
    ua.on("message", resp => {
      let msg = JSON.parse(resp);
      assert.equal(msg.success, true);
      assert.equal(msg.type, "message");
      assert.equal(msg.text, message);
      if (++cnt >= 2) { complete(); }
    });
    ua.on("message", resp => {
      let msg = JSON.parse(resp);
      assert.equal(msg.success, true);
      assert.equal(msg.type, "message");
      assert.equal(msg.text, message);
      if (++cnt >= 2) { complete(); }
    });
  });
});
