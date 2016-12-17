- Video json:

```JSON
{
    "_id" : ObjectId("58556840d48cdc00042d8565"),
    "url" : {
        "site" : "youtube",
        "link" : "asldasd"
    },
    "comment" : [ 
        {
            "_id" : ObjectId("585572983f57a100041a5713"),
            "content" : "hello",
            "time" : 1,
            "timeupdated" : ISODate("2016-12-17T17:15:03.945Z")
        }, 
        {
            "_id" : ObjectId("5855729b3f57a100041a5714"),
            "content" : "hello",
            "time" : 2,
            "timeupdated" : ISODate("2016-12-17T17:15:07.854Z")
        }
    ]
}
```

- GET /video/link?site=youtube&id=asldasd

- GET /video/id/58556840d48cdc00042d8565

- GET /video/id/58556840d48cdc00042d8565/subtitle

- POST /video/58556840d48cdc00042d8565

  - Request body JSON: `{"content":"hello","time":1}`