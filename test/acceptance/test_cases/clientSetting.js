const config = require("./config");
const request = require('supertest')(config.server)

describe('client-settings', function () {
	const clientSettingObj = {
		"clientId": 1,
		"group": "sales"
	}

	it('post client-settings successfully should return 201 status', async () => {
		await request
			.post('/client-settings')
			.set(config.headers)
			.send(clientSettingObj)
			.expect(201, clientSettingObj);
	})

	it('post duplicate client-settings should return 409 status', async () => {
		await request.post('/client-settings')
			.set(config.headers)
			.send(clientSettingObj)
			.expect(409, {
				"code": "CLIENT_DUPLICATED",
				"detail": "client setting already exists",
				"status": 409,
				"title": "duplicated record",
				"type": "about:blank"
			});
	})

	it('get client-settings successfully should return 200 status', async () => {
		await request
			.get('/client-settings/' + clientSettingObj.clientId)
			.set(config.headers)
			.expect(200, clientSettingObj);
	})

	it('delete client-settings successfully should return 204 status', async () => {
		await request
			.delete('/client-settings/' + clientSettingObj.clientId)
			.set(config.headers)
			.expect(204, {});
	})

	it('post client-settings with empty body should return 400 status', async () => {
		await request
			.post('/client-settings')
			.set(config.headers)
			.send()
			.expect(400, invalidBodyError());
	})

	it('get client-settings with inexistent client id should return 404 status', async () => {
		await request
			.get('/client-settings/9')
			.set(config.headers)
			.expect(404, notFoundError(9));
	})

	it('delete client-settings with inexistent client id should return 404 status', async () => {
		await request
			.delete('/client-settings/9')
			.set(config.headers)
			.expect(404, notFoundError(9));
	})

	it('get client-settings with client id already deleted should return 404 status', async () => {
		await request
			.get('/client-settings/' + clientSettingObj.clientId)
			.set(config.headers)
			.expect(404, notFoundError(clientSettingObj.clientId));
	})

})

function notFoundError(id) {
	return {
		"code": "CLIENT_NOT_FOUND",
		"detail": `client setting not found: ${id}`,
		"status": 404,
		"title": "resource not found",
		"type": "about:blank"
	}
}

function invalidBodyError() {
	return {
		"type": "about:blank",
		"status": 400,
		"code": "INVALID_BODY",
		"title": "invalid argument",
		"detail": "you have applied a request with an invalid body. Please review the body and check the structure"
	}
}
