const webdriverio = require('webdriverio');
const options = {
	desiredCapabilities: {
		browserName: 'chrome'
	}
};

webdriverio
    .remote(options)
    .init()
    .url('https://www.google.com')
    .getTitle().then(function(title) {
        console.log('Title was: ' + title);
    })
    .end();
