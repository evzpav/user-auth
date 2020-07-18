const config = require("./config");
const request = require('supertest')(config.server)

describe('label-group', function () {
	const customGroupObj = { "id": "sales" };

	it('post a label-group successfully should return 201 status', async () => {
		await request
			.post('/label-groups')
			.set(config.headers)
			.send(customGroupObj)
			.expect(201, customGroupObj);
	})

	it('post duplicate label-group should return 409 status', async () => {
		await request
			.post('/label-groups')
			.set(config.headers)
			.send(customGroupObj)
			.expect(409, {
				"code": "LABEL_GROUP_DUPLICATED",
				"detail": "label group already exists",
				"status": 409,
				"title": "duplicated record",
				"type": "about:blank"
			});
	})

	it('post a label-group with empty body should return 400 status', async () => {
		await request
			.post('/label-groups')
			.set(config.headers)
			.send()
			.expect(400, invalidBodyError());
	})

	it('get label-groups successfully should return 200 status', async () => {
		await request
			.get('/label-groups')
			.set(config.headers)
			.expect(200, {
				labelGroups: [customGroupObj]
			});
	})

	it('delete a label-group successfully should return 204 status', async () => {
		await request
			.delete('/label-groups/' + customGroupObj.id)
			.set(config.headers)
			.expect(204, {});
	})

	it('get label-groups inexistent should return 200 status with empty array', async () => {
		await request
			.get('/label-groups')
			.set(config.headers)
			.expect(200, {
				labelGroups: []
			});
	})

	it('delete a label-group with inexistent client id should return 404 status', async () => {
		await request
			.delete('/label-groups/9')
			.set(config.headers)
			.expect(404, notFoundError(9));
	})

	it('get label-groups already deleted should return 200 status with empty array', async () => {
		await request
			.get('/label-groups')
			.set(config.headers)
			.expect(200, {
				labelGroups: []
			});
	})

})

function notFoundError(id) {
	return {
		"code": "LABEL_GROUP_NOT_FOUND",
		"detail": `label group not found: ${id}`,
		"status": 404,
		"title": "resource not found",
		"type": "about:blank",
		"arguments": {
			"groupID": id + ""
		}
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
