const config = require("./config");
const request = require('supertest')(config.server)
const assert = require("assert");

describe('references', function () {

	const referenceId = "empresa.veiculos"
	const label = {
		"description": "descrição em português",
		"language": "pt_BR",
		"label": "label em português",
		"shortDescription": "descrição curta em português",
		"group": "default",
		referenceId
	}
	const reference = {
		"label": {
			"en_US": "label in english",
			"pt_BR": "label em português",
		},
		"description": {
			"en_US": "description in english",
			"pt_BR": "descrição em português",
		},
		"shortDescription": {
			"en_US": "short description in english",
			"pt_BR": "descrição curta em português",
		}
	}

	const referenceWithLabel = {
		"label": {
			"pt_BR": "label em português alterado",
		},
	}

	const referenceModifiedWithLabel = {
		referenceId,
		...referenceWithLabel
	}

	const customGroup = "sales"


	const referenceWithId = {
		referenceId,
		...reference
	}

	const referenceWithGroup = {
		"group": "default",
		...referenceWithId
	}

	const labelGroup = {
		"group": customGroup,
		"label": {
			"pt_BR": customGroup + " label em português",
			"en_US": customGroup + " label in english"
		},
		"description": {
			"pt_BR": customGroup + " descrição em português",
			"en_US": customGroup + " description in english"
		},
		"shortDescription": {
			"pt_BR": customGroup + " descrição curta em português",
			"en_US": customGroup + " short description in english"
		}
	}

	const labelGroupWithRefId = {
		referenceId,
		...labelGroup
	}


	//REFERENCES
	it('post references successfully should return 201 status', async () => {
		await request
			.post('/references')
			.set(config.headers)
			.send(referenceWithId)
			.expect(201, referenceWithGroup);
	})

	//REFERENCES
	it('post duplicate references should return 409 status', async () => {
		await request
			.post('/references')
			.set(config.headers)
			.send(referenceWithId)
			.expect(409, duplicateError("reference"));
	})

	//REFERENCES
	it('post references with empty body should return 400 status', async () => {
		await request
			.post('/references')
			.set(config.headers)
			.send()
			.expect(400, invalidBodyError());
	})

	//REFERENCES
	it('get references successfully should return 200 status', async () => {
		await request
			.get('/references')
			.set(config.headers)
			.expect(200, { "references": [referenceWithGroup] });
	})

	//REFERENCES
	it('get references with lang query param successfully should return 200 status', async () => {
		const language = "pt_BR";

		await request
			.get(`/references?lang=${language}`)
			.set(config.headers)
			.expect(200)
			.expect((res) => {
				const body = res.body;
				const refLabel = body.references[0]
				assert.equal(refLabel.label, label.label)
				assert.equal(refLabel.description, label.description)
				assert.equal(refLabel.shortDescription, label.shortDescription)
				assert.equal(refLabel.language, language)
			});
	})

	//REFERENCES
	it('get references with client query param successfully should return 200 status', async () => {
		const client = "2";

		await request
			.get(`/references?client=${client}`)
			.set(config.headers)
			.expect(200)
			.expect(200, { "references": [referenceWithGroup] });
	})

	//REFERENCES
	it('put references successfully should return 200 status', async () => {
		await request
			.put('/references/' + referenceId)
			.set(config.headers)
			.send(reference)
			.expect(200, referenceWithId);
	})

	it('put references successfully removing description and short descripton should return 200 status', async () => {
		await request
			.put('/references/' + referenceId)
			.set(config.headers)
			.send(referenceWithLabel)
			.expect(200, referenceModifiedWithLabel);
	})

	//LABEL-GROUP
	it('post a label-group successfully should return 201 status', async () => {
		await request
			.post('/label-groups')
			.set(config.headers)
			.send({ "id": customGroup })
			.expect(201, { "id": customGroup });
	})

	//CUSTOM-LABEL-GROUP
	it('post custom group labels from known referenceID should return 201 status', async () => {
		await request
			.post(`/references/${referenceId}/label-groups`)
			.set(config.headers)
			.send(labelGroup)
			.expect(201, labelGroupWithRefId);
	})

	//REFERENCES
	it('get references with group query param successfully should return 200 status', async () => {
		await request
			.get(`/references?group=${customGroup}`)
			.set(config.headers)
			.expect(200)
			.expect({ "references": [labelGroupWithRefId] })
	})

	//CUSTOM-LABEL-GROUP
	it('post duplicate custom group labels should return 409 status', async () => {
		await request
			.post(`/references/${referenceId}/label-groups`)
			.set(config.headers)
			.send(labelGroup)
			.expect(409, duplicateError("reference"));
	})

	//CUSTOM-LABEL-GROUP
	it('get custom group labels from known referenceID should return 200 status', async () => {
		await request
			.get(`/references/${referenceId}/label-groups`)
			.set(config.headers)
			.expect(200, { "references": [labelGroupWithRefId] });
	})

	//CUSTOM-LABEL-GROUP
	it('put custom group labels from known referenceID should return 200 status', async () => {
		await request
			.put(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.send(labelGroup)
			.expect(200, labelGroupWithRefId);
	})

	//CUSTOM-LABEL-GROUP
	it('put custom group labels from inexistent referenceID should return 200 status', async () => {
		const referenceId = "inexistentRefId"
		await request
			.put(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.send(labelGroup)
			.expect(404, notFoundError("reference", "referenceID", referenceId));
	})

	//CUSTOM-LABEL-GROUP
	it('put custom group labels from inexistent customGroup should return 200 status', async () => {
		const customGroup = "inexistentGroup"
		await request
			.put(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.send(labelGroup)
			.expect(404, notFoundError("label group", "groupID", customGroup));
	})


	//CUSTOM-LABEL-GROUP
	it('delete custom group labels from known referenceID should return 204 status', async () => {
		await request
			.delete(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.expect(204, {});
	})

	//CUSTOM-LABEL-GROUP
	it('delete custom group labels from inexistent referenceId should return 404 status', async () => {
		const referenceId = "inexistentRefId"
		await request
			.delete(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.expect(404, notFoundError("reference", "referenceID", referenceId));
	})

	//CUSTOM-LABEL-GROUP
	it('delete custom group labels from inexistent customGroup should return 404 status', async () => {
		const customGroup = "inexistentGroup"
		await request
			.delete(`/references/${referenceId}/label-groups/${customGroup}`)
			.set(config.headers)
			.expect(404, notFoundError("reference", "referenceID", referenceId));
	})

	//REFERENCES
	it('delete references successfully should return 204 status', async () => {
		await request
			.delete('/references/' + referenceId)
			.set(config.headers)
			.expect(204, {});
	})

	//REFERENCES
	it('delete references from inexistent referenceId should return 404 status', async () => {
		const referenceId = "inexistentRefId"
		await request
			.delete('/references/' + referenceId)
			.set(config.headers)
			.expect(404, notFoundError("reference", "referenceID", referenceId));
	})

	//REFERENCES
	it('get references after deleting successfully should return 404 status', async () => {
		await request
			.get('/references/' + referenceId)
			.set(config.headers)
			.expect(404, notFoundError("reference", "referenceID", referenceId));
	})

	//LABEL-GROUP
	it('delete a label-group successfully should return 204 status', async () => {
		await request
			.delete(`/label-groups/${customGroup}`)
			.set(config.headers)
			.expect(204, {});
	})

})

function notFoundError(attribute, key, value) {
	let code = attribute

	if (code.indexOf(" ") > 0) {
		code = code.replace(" ", "_")
	}

	let notFoundObj = {
		"code": code.toUpperCase() + "_NOT_FOUND",
		"detail": `${attribute} not found: ${value}`,
		"status": 404,
		"title": "resource not found",
		"type": "about:blank",
		"arguments": {}
	}
	notFoundObj.arguments[key] = value + "";

	return notFoundObj;
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

function duplicateError(attribute) {
	return {
		"code": attribute.toUpperCase() + "_DUPLICATED",
		"detail": attribute + " already exists",
		"status": 409,
		"title": "duplicated record",
		"type": "about:blank"
	}
}
