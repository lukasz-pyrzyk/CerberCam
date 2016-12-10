module.exports = {
    send: function(config, msg) {
        var amqp = require('amqplib/callback_api');
        var protobuf = require('protocol-buffers')

        console.log("connecting to the queue %s", config.host)
        amqp.connect(config.host, function(err, conn) {
            if(err != null) console.log(err)
            else console.log("Connected successfully")

            conn.createChannel(function(err, ch) {
                if(err != null) console.log(err)
                else console.log("Channel selected successfully")

                var q = config.requests;
                console.log("Selected queue %s", q)

                ch.assertQueue(q, {durable: false});
                ch.sendToQueue(q, new Buffer('Hello World!'));
                console.log("Message has been sent");
            });
        });
    }
}