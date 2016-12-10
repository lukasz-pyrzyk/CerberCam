module.exports = {
    createFromModel: function(msg) {
        var protobuf = require('protocol-buffers')

        var proto = protobuf(fs.readFileSync('../Messages/message.proto'))
        var buf = proto.Message.encode({
            Email: msg.Email,
            Photo: msg.Photo
        });
        
        return buf;
    }
}