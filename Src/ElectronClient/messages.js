module.exports = {
    createDummyMessage: function(fileName) {
        var protobuf = require('protocol-buffers')

        // pass a proto file as a buffer/string or pass a parsed protobuf-schema object
        var proto = protobuf(fs.readFileSync('../Messages/message.proto'))
        var buf = proto.Message.encode({
            Email: "lukasz.pyrzyk@gmail.com",
            Photo: fs.readFileSync('/home/lukasz/Dev/iphone')
        });
        
        return buf;
    }
}