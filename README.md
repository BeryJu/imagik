# GoPyazo

Pyazo, but fast. GoPyazo is a fast and lightweight fileserver, that

- lets you access files by path
- lets you access files by their Hash (MD5/SHA1/SHA2/SHA512)
- lets you upload files with a very simple API

## Running

### Docker

Run the container like this:

```
docker run -v "whatever directory you want to share":/data -w /data beryju/gopyazo:latest-amd64
```

Now you can access gopyazo on http://localhost:8080

### Binary

Download a binary from [GitHub](https://github.com/BeryJu/gopyazo/releases) and run it:

```
./gopyazo server
```

Now you can access gopyazo on http://localhost:8080

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.

## API

### `/api/pub/health/liveness`

Healthcheck endpoint, which returns a 201 Response as soon as gopyazo is running.

### `/api/pub/health/readiness`

Healthcheck Readiness probe, which returns a 201 after the Hash Map has been populated, otherwise a 500.

### `/api/priv/list`

**Requires authentication**

List contents of a directory. Accepts a query parameter `pathOffset`, which is appended to the root directory.

### `/api/priv/move`

**Requires authentication**

Move a file. Requires two query parameters, `from` and `to`, which are relative to the root directory.
