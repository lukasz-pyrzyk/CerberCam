const rabbitmq = require('../queue')
const proto = require('../modules/proto')

module.exports = {
    sendPhoto: function(cfg, msg) {
        var protoMessage = proto.createFromModel(msg);
        rabbitmq.send(cfg.queue, protoMessage)
    }
}