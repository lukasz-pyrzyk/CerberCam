const rabbitmq = require('../queue')
const proto = require('../messages')
const protobuf = require('protocol-buffers')

module.exports = {
    sendPhoto: function(cfg, msg) {
        var protoMessage = proto.createFrom(msg);
        rabbitmq.send(cfg.queue, protoMessage)
    }
}