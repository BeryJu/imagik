# Imagik

Pyazo, but fast. Imagik is a fast and lightweight fileserver, that

- lets you access files by path
- lets you access files by their Hash (MD5/SHA1/SHA2/SHA512)
- lets you upload files with a very simple API

## Running

### Docker

Run the container like this:

```
docker run -p 8000:8000 -v "whatever directory you want to share":/share beryju/imagik
```

Now you can access imagik on http://localhost:8000

### Binary

Download a binary from [GitHub](https://github.com/BeryJu/imagik/releases) and run it:

```
./imagik
```

Now you can access imagik on http://localhost:8000

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.

## API

### GET `/<path>`

Retrieve file stored at `path`.

### GET `/<path>?meta`

Retrieve metadata for file stored at `path`.

### PUT `/<path>`

**Requires authentication**

Accepts file uploads from the HTTP Request body, like using `curl --data "@/path/to/filename"`.
Returns a JSON object with all the hashes,
```json
{
    "SHA128":"acd5aeeb3c8d1cf580a59bc3e125d249ecdd0eda",
    "SHA256":"e6b104c1420af07013b4378ddacaaa3938259422f07d5d47f7ea114cf9de80cf",
    "SHA512":"10c08e2134fb953f891c2a3655f3744c0321fa72aefdf6bff000eff0a3f7882a008fff477dfec9aa22519ad17fb0fafd602caf3773cb848a5250131fdf8559ab",
    "SHA512Short":"10c08e2134fb953f",
    "MD5":"7e97fa079923fcdb39eb39b480729f36"
}
```

### GET `/api/pub/health/liveness`

Healthcheck endpoint, which returns a 201 Response as soon as imagik is running.

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

## Migrating from pyazo

If you didn't use Collections in pyazo, you can simple re-use the same Media folder for imagik, and all URLs will continue to work.

If you did use Collections, use the script below, to mirror your Collection Structure into Filesystem folders, which are used by imagik.

```python
# Execute this in your pyazo installation directory
# docker-compose exec server ./manage.py shell
# Then paste the contents below into the shell.
# This will output the commands required to move the files
# into folders.
from pyazo.core.models import *
for c in Collection.objects.all():
    print(f"mkdir {c.name}")
    for o in c.object_set.all():
        rel_path = o.file.path.replace('/app/media/', '')
        print(f"mv {rel_path} {c.name}/{rel_path}")
```
