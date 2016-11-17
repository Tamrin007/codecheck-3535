"use strict";

const ws = require("ws");
const setWsProtocol = require("./utils").setWsProtocol;
const assert = require("chai").assert;
const env = require("../config/env.json");

let host_url = process.env.HOST_URL || `ws://${env.hostname}:${env.port}${env.path}`;
const HOST_URL = setWsProtocol(host_url);

describe("bot ping", () => {
  it ("should return `pong` when mentioned with `bot`", done => {
    let user = new ws(HOST_URL);
    let message = { text: "bot ping" };
    user.on("open", () => {
      return user.send(JSON.stringify(message));
    });
    let cnt = 0;
    user.on("message", resp => {
      let msg = JSON.parse(resp);

      assert.equal(msg.success, true);
      switch (msg.type) {
        case "message":
          assert.equal(msg.text, message.text); break;
        case "bot":
          assert.equal(msg.text, "pong"); break;
        default:
          assert.equal(msg.type, "message or bot", "expected 'message' or 'bot'");
          break;
      }
      if (++cnt >= 2) {
        user.close();
        done();
      }
    });
  });

  it ("should return `pong` when mentioned with `@bot`", done => {
    let user = new ws(HOST_URL);
    let message = { text: "@bot ping" };
    user.on("open", () => {
      return user.send(JSON.stringify(message));
    });
    let cnt = 0;
    user.on("message", resp => {
      let msg = JSON.parse(resp);

      assert.equal(msg.success, true);
      switch (msg.type) {
        case "message":
          assert.equal(msg.text, message.text); break;
        case "bot":
          assert.equal(msg.text, "pong"); break;
        default:
          assert.fail(msg.type, "message or bot", "'message' or 'bot' expected");
          break;
      }
      if (++cnt >= 2) {
        user.close();
        done();
      }
    });
  });

  it ("should recognize with a text start with `bot:`", done => {
    let user = new ws(HOST_URL);
    let message = { text: "bot:ping" };
    user.on("open", () => {
      return user.send(JSON.stringify(message));
    });

    let cnt = 0;
    user.on("message", resp => {
      let msg = JSON.parse(resp);

      assert.equal(msg.success, true);
      switch (msg.type) {
        case "message":
          assert.equal(msg.text, message.text); break;
        case "bot":
          assert.equal(msg.text, "pong"); break;
        default:
          assert.fail(msg.type, "message or bot", "'message' or 'bot' expected");
          break;
      }
      if (++cnt >= 2) {
        user.close();
        done();
      }
    });
  });

});
