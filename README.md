### TODO

- Webinterface
    - Browser of folders and files
    - Move files between folders
    -



# GoPyazo

Pyazo, but fast. GoPyazo is a fast and lightweight fileserver, that

- lets you access files by path
- lets you access files by their Hash (MD5/SHA1/SHA2/SHA512)
- lets you upload files with a very simple API

## Running

```
Usage:
  gopyazo [directory to serve] [flags]

Flags:
      --cache-enabled             Enable in-memory cache
      --cache-eviction int        Time after which entry can be evicted (in minutes) (default 10)
      --cache-max-item-size int   Maximum Item size to cache (in bytes) (default 500)
      --cache-max-size int        Maximum Cache size in MB (0 disables the limit)
      --debug                     Enable debug-mode
      --exif-purge-gps            Purge GPS-Related EXIF metadata (default true)
  -h, --help                      help for pixie
      --silent                    Enable silent mode (no access logs)
      --spa-mode                  Enable SPA-mode (redirect all requests to missing files to /index.html
```

### Docker

Run the container like this:

```
docker run -v "whatever directory you want to share":/data -w /data beryju/pixie:latest-amd64
```

Now you can access pixie on http://localhost:8080

### Binary

Download a binary from [GitHub](https://github.com/BeryJu/pixie/releases) and run it:

```
./pixie /data
```

Now you can access pixie on http://localhost:8080

## Configuration

By default, a gallery is shown for every folder. To prevent this, create an empty `index.html` file in the folder.

## API

### `/-/ping`

Healthcheck endpoint, which returns `pong` with a 200 Status Code. Useful for Kubernetes Deployments and general Monitoring.

### `/<directory>/?json`

Lists directory Contents as JSON, used by the Gallery Page to load all files to be displayed.

### `/<file.ext>?meta`

Return file's metadata, in the following format:

```json
{
    "name": "Canon_DIGITAL_IXUS_400.jpg",
    "size": 9198,
    "content_type": "image/jpeg",
    "exif": {
        "ApertureValue": "",
        "ColorSpace": "",
        "ComponentsConfiguration": "",
        "CompressedBitsPerPixel": "",
        "CustomRendered": "",
        "DateTime": "2008:07:31 17:15:01",
        "DateTimeDigitized": "2004:08:27 13:52:55",
        "DateTimeOriginal": "2004:08:27 13:52:55",
        "DigitalZoomRatio": "",
        "ExifIFDPointer": "",
        "ExifVersion": "",
        "ExposureBiasValue": "",
        "ExposureMode": "",
        "ExposureTime": "",
        "FNumber": "",
        "FileSource": "",
        "Flash": "",
        "FlashpixVersion": "",
        "FocalLength": "",
        "FocalPlaneResolutionUnit": "",
        "FocalPlaneXResolution": "",
        "FocalPlaneYResolution": "",
        "InteroperabilityIFDPointer": "",
        "InteroperabilityIndex": "R98",
        "Make": "Canon",
        "MakerNote": "",
        "MaxApertureValue": "",
        "MeteringMode": "",
        "Model": "Canon DIGITAL IXUS 400",
        "PixelXDimension": "",
        "PixelYDimension": "",
        "ResolutionUnit": "",
        "SceneCaptureType": "",
        "SensingMethod": "",
        "ShutterSpeedValue": "",
        "Software": "GIMP 2.4.5",
        "ThumbJPEGInterchangeFormat": "",
        "ThumbJPEGInterchangeFormatLength": "",
        "WhiteBalance": "",
        "XResolution": "",
        "YCbCrPositioning": "",
        "YResolution": ""
    }
}
```
