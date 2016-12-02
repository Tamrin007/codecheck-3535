'use strict';

var ws = new WebSocket('ws://sprint.tamrin.tech/chat');

$(function () {
  $('form').submit(function(){
    var $this = $(this);
    // ws.onopen = function() {
    //   console.log('sent message: %s', $('#m').val());
    // };
    var message = $('#m').val()
    ws.send(JSON.stringify({"text": message}));
    $('#m').val('');
    return false;
  });
  ws.onmessage = function(msg){
    console.log(msg.data);
    var resp = JSON.parse(msg.data);
    $('#messages')
      .append($('<li>')
      .append($('<span class="message">').text(resp.text)));
  };
  ws.onerror = function(err){
    console.log("err", err);
  };
  ws.onclose = function close() {
    console.log('disconnected');
  };
});
