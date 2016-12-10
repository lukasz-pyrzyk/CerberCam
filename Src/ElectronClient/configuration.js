module.exports = {
    load: function(fileName) {
        yaml = require('js-yaml');
        fs   = require('fs');

        var cfg = yaml.safeLoad(fs.readFileSync('../config.yaml', 'utf8'));
        return cfg;
    }
}