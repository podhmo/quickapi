# todo

```console
$ http :8080/todos
HTTP/1.1 200 OK
Content-Length: 89
Content-Type: application/json; charset=utf-8
Date: Mon, 18 Jul 2022 16:15:19 GMT

{
    "items": [
        {
            "bool": false,
            "id": 1,
            "title": "hello"
        },
        {
            "bool": false,
            "id": 3,
            "title": "byebye"
        }
    ]
}
```