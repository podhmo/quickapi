# openapi-petstore

- [apidoc.md](apidoc.md)
- [openapi.json](openapi.json)


```console
$ http :8080/pets
HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json; charset=utf-8
Date: Mon, 18 Jul 2022 16:16:04 GMT

{
    "items": []
}
```

post

```console
$ echo '{"name": "foo"}' | http :8080/pets
HTTP/1.1 200 OK
Content-Length: 20
Content-Type: application/json; charset=utf-8
Date: Mon, 18 Jul 2022 16:16:55 GMT

{
    "id": "",
    "name": "" 
}
```

validation error

```console
$ echo '{}' | http :8080/pets
HTTP/1.1 400 Bad Request
Content-Length: 553
Content-Type: application/json; charset=utf-8
Date: Mon, 18 Jul 2022 16:17:30 GMT

{
    "code": 400,
    "detail": [
        "request body has an error: doesn't match the schema: Error at \"/name\": property \"name\" is missing",
        "Schema:",
        "  {",
        "    \"additionalProperties\": false,",
        "    \"properties\": {",
        "      \"name\": {",
        "        \"type\": \"string\"",
        "      },",
        "      \"tag\": {",
        "        \"type\": \"string\"",
        "      }",
        "    },",
        "    \"required\": [",
        "      \"name\"",
        "    ],",
        "    \"type\": \"object\"",
        "  }",
        "",
        "Value:",
        "  {}",
        ""
    ],
    "error": "request body has an error: doesn't match the schema: Error at \"/name\": property \"name\" is missing"
}
```