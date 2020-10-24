# GoPyazo

Pyazo, but fast. GoPyazo is a fast and lightweight fileserver, that

- lets you access files by path
- lets you access files by their Hash (MD5/SHA1/SHA2/SHA512)
- lets you upload files with a very simple API

## Running

### Docker

Run the container like this:

```
docker run -p 8000:8000 -v "whatever directory you want to share":/share beryju/gopyazo server
```

Now you can access gopyazo on http://localhost:8000

### Binary

Download a binary from [GitHub](https://github.com/BeryJu/gopyazo/releases) and run it:

```
./gopyazo server
```

Now you can access gopyazo on http://localhost:8000

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.

## API

### GET `/<path>`

Retrieve file stored at `path`.

### GET `/<path>?meta`

Retrieve metadata for file stored at `path`.

### PUT `/<path>`

**Requires authentication**

Accepts file uploads from the HTTP Request body, like using `curl --data "@/path/to/filename"`

### GET `/api/pub/health/liveness`

Healthcheck endpoint, which returns a 201 Response as soon as gopyazo is running.

### GET `/api/pub/health/readiness`

Healthcheck Readiness probe, which returns a 201 after the Hash Map has been populated, otherwise a 500.

### GET `/api/priv/list[?pathOffset=]`

**Requires authentication**

List contents of a directory. Accepts a query parameter `pathOffset`, which is appended to the root directory.

### POST `/api/priv/move?to=&from=`

**Requires authentication**

Move a file. Requires two query parameters, `from` and `to`, which are relative to the root directory.

### POST `/api/priv/upload`

**Requires authentication**

Accepts Multipart-Form Encoded files and uploads them to the respective path from the form relative to the root directory.
