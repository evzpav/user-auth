const config = require("./config");
const request = require('supertest')(config.server)

describe('schema', function () {
    const expected = {
        "schemas": [
            {
                "_id": "1",
                "name": "cfo",
                "provider": "evzpav",
                "subject": "evzpav-cfo-source",
                "source": "",
                "schema": {
                    "properties": {
                        "data": {
                            "properties": {
                                "isMatriz": {
                                    "type": "boolean"
                                },
                                "partes": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "required": [
                                            "advogados"
                                        ],
                                        "properties": {
                                            "advogados": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "required": [
                                                        "tipoAdvogado",
                                                        "nomeAdvogado",
                                                        "cpf",
                                                        "cnpj",
                                                        "oab",
                                                        "dataEntrada",
                                                        "dataExclusao"
                                                    ],
                                                    "properties": {
                                                        "tipoAdvogado": {
                                                            "type": "string"
                                                        },
                                                        "nomeAdvogado": {
                                                            "type": "string"
                                                        },
                                                        "cpf": {
                                                            "type": "string"
                                                        },
                                                        "cnpj": {
                                                            "type": "string"
                                                        },
                                                        "oab": {
                                                            "type": "string"
                                                        },
                                                        "dataEntrada": {
                                                            "type": "string"
                                                        },
                                                        "dataExclusao": {
                                                            "type": "string"
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            },
                            "required": [
                                "isMatriz"
                            ],
                            "type": "object"
                        }
                    },
                    "required": [
                        "data"
                    ],
                    "type": "object"
                }
            }
        ]
    }
    it('test schemas', async () => {
        await request
			.get('/clients/1/providers/evzpav/schemas/cfo/types/source')
			.set(config.headers)
			.expect(200, expected);
    })
})
