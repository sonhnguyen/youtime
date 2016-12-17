- Video json:

```JSON
{
  "id": "585557a9ed4f1fa336b9e1fd",
  "url": {
    "site": "youtube",
    "link": "asldasd"
  },
  "comment": [
    {
      "content": "hello",
      "time": 1,
      "datecreated": "2016-12-17T22:35:22.471+07:00"
    },
    {
      "content": "hello",
      "time": 2,
      "datecreated": "2016-12-17T22:35:03.873+07:00"
    }
  ]
}
```

- GET /video/link?site=youtube&link=asldasd

- GET /video/id/585557a9ed4f1fa336b9e1fd

- POST /video/id/585557a9ed4f1fa336b9e1fd

  - Request body JSON: `{"content":"hello","time":1}`