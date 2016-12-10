module.exports = {
    load: function(fileName) {
        console.log("loading configuration")
        yaml = require('js-yaml')
        fs   = require('fs')

        var cfg
        try {
            cfg = yaml.safeLoad(fs.readFileSync(fileName, 'utf8'))
            console.log("Configuration loaded")
        } catch (error) {
            console.log("Unable to load configuration")
            console.log(error)
        }
        
        console.log(cfg)
        return cfg;
    }
}