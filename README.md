# plexrenamer

Rename Japan EPG based files to Plex Media Library compliant file structure

## Build and Publish

```
docker build -t toruta39/plexrenamer:latest .
```

## Test

```
go test
```

## Run

```
docker run --rm -v "$(pwd)/artifacts:/media" toruta39/plexrenamer:latest -dryrun
```