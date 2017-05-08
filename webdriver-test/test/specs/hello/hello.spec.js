const assert = require('assert');

describe('Hello world test', function() {
	it('should get google page and its title', function() {
		browser.url('https://google.com');
		let title = browser.getTitle();
		assert.equal(title, 'Google');
	});
});
