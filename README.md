# Media Uploader

This is a simple media uploader, allowing you to upload a media on a server
(file, memory, s3, ...) with a few metadata such as a name or some tags.

## How to build

### From Source

```bash
git clone https://github.com/Taluu/media-api.git
cd media-api
go build -o build/media-api ./app
build/media-api
```

Available flags are `-domain` and `-port`, with respective default values
being `localhost` and `8080`.

### From binary release

Binaries should be released on the Releases page on the github repo.

## Running tests

Tests are packaged in the repository, to run them, just run the following
command :

```bash
go test ./pkg/...
```

You can also check the code coverage with usual flags from the `go test` and
`go tool` commands :

```bash
 go test ./pkg/... -coverprofile /tmp/coverage.out -covermode count
 go tool cover -html=/tmp/coverage.out -o /tmp/cover.html
```

## Endpoints

Note : all example are hitting as if the domain is `localhost` and the port
`8080`, and the key `test`.

### Creating a media

This is one special endpoint, as instead of sending a plain json body as the
other endpoints, you have to send a multipart/form-data with a `POST /medias`
request as follows :

```bash
curl -X POST -H "Content-Type: multipart/form-data" -F "media=@/path/to/file.ext" -F "data={\"tags\": [\"foo\", \"bar\"]};type=application/json" http://localhost:8080/medias
```

You should then get a 201 json response like the following :

```json
{
  "id": "121a7a2c-5777-40e8-8c27-425c3777f378",
  "name": "file.ext",
  "file": "http://localhost:8080/viewer/121a7a2c-5777-40e8-8c27-425c3777f378",
  "tags": ["foo", "bar"]
}
```

The `data` field is optionnal, while being a json object (encoded as a string) ;
it can contain a `name` string field, and a `tags` string array. In the example
above, only the `tags` property was provided, resulting in the filename used as
a name. Another example with only a `name` :

```bash
curl -X POST -H "Content-Type: multipart/form-data" -F "media=@/path/to/file.ext" -F "data={\"name\": \"my media\"};type=application/json" http://localhost:8080/medias
```

You should then get a 201 json response like the following :

```json
{
  "id": "1986600e-d65c-4c04-b2df-2cca4299ff62",
  "name": "my media",
  "file": "http://localhost:8080/viewer/1986600e-d65c-4c04-b2df-2cca4299ff62",
  "tags": []
}
```

Note that provided tags in the request will be created if they do not already
exist.

### Creating a tag

To create a new tag, just send the following json to the `POST /tags` endpoint :

```json
{
  "name": "foo"
}
```

So with a curl command :

```bash
curl -X POST http://localhost:8080/tags -H "Content-type: application/json" -d "{\"name\": \"foo\"}"
```

You will then have a 201, with the following response :

```json
{
  "name": "foo"
}
```

You will have a 400 if the json body is malformed or no name are provided.

### Listing available tags

You can get all the tags currently available and registered with a call to the
`GET /tags` endpoint :

```bash
curl http://localhost:8080/tags -H "Content-type: application/json"
```

You will then receive a 200 response with the following content :

```json
{
  "tags": [{ "name": "foo" }, { "name": "bar" }]
}
```

### Searching a media by a tag

You can search all medias that are tagged with a specific tag by sending a
request to the `GET /medias/{tagName}` endpoint. For example, with a `foo` tag :

```bash
curl http://localhost:8080/medias/foo -H "Content-type: application/json"
```

You will then get a 200 response returning the list of medias that match the
request :

```json
{
  "medias": [
    {
      "id": "121a7a2c-5777-40e8-8c27-425c3777f378",
      "name": "file.ext",
      "file": "http://localhost:8080/viewer/121a7a2c-5777-40e8-8c27-425c3777f378",
      "tags": ["foo", "bar"]
    }
  ]
}
```

If the tag doesn't exist or no medias are associated with it, it will still
return a 200 but with an empty `medias` array.

### Downloading a media

Even if this was not asked in the test, I added an endpoint to be able to
download an uploaded media on the `GET /viewer/{mediaID}` endpoint :

```bash
curl http://localhost:8080/viewer/121a7a2c-5777-40e8-8c27-425c3777f378 -H "Content-type: application/json" --output /tmp/file.ext
```

You will then receive the file content as a 200 response. If the media is not
found, or for some reasons its corresponding file can't be found, you will then
have a 404 with the following json body :

```json
{
  "code": 404,
  "error": "media not found"
}
```

## Feedback

Overall, this exercice was particularly enjoyable, as I didn't work with file
upload in a long time (the challenge was mostly on how to pass some metadata
with the file to upload, and also how to curl it the proper way).

I also made the deliberate choice to go the simple option with not bothering on
a few things that should be taken care of in a prod environment, such as

- Checking on media upload the media type and its size, putting some restrictions
  to avoid pitfalls such as upload a malicious file or way too big files as
  "medias".
- Having a proper storage. Currently all this program does is storing in memory
  the data, which means that if the program is stopped and re-executed, the data
  that was there will be gone. This is easily fixable by implementing the proper
  interfaces `media.MediaRepository`, `media.TagRegistry` and
  `media.MediaUploader` ; a `file uploader` is available as an example, even if
  it is not used in the application, and all this setup is within the
  `app/main.go` file.

For the tests, I only made unit tests ; for integration tests, I usually like
to work with cucumber (or equivalent), but I do not have the time to put it
there (and I also think it's a bit out of scope for this test).

Note that I tried to approach the exercise with a DDD approach, even though I
had only one domain to bother with, which was the medias handling.
But still, I decided to add the intefaces to try to keep an open mind on how
to extend what I did there, such as stated above with using another uploader
rather the one that stores everything in memory (it could be writing on disk
as provided as an example in this repository, a s3 uploader instead, ... pick
your poison !)

You will also notice that I split what "launches" the application (what is in
the `app` directory at the root) so that it's the control tower than knows
how to instanciate providers, listen to the server, ... for a future
improvment, something like a DIC could be added so no modifications or almost
is needed when adding or removing things. But again, I felt this was out of
scope for this test.

I also added a middleware, purely as an example to log requests. I'm pretty
sure there are other and more meaningful ways of achieving that, but for
simplicity sake, "it works" (tm)
