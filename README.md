# go-livepeer-bitexact-verifier

Verification code based on bitexact comparison of video transcoding.

## Building

```
cd go-livepeer-bitexact-verifier
./install_ffmpeg.sh
go build -v .
```

## Docker

```
cd go-livepeer-bitexact-verifier
# Build Docker image
docker build -f verifier .
# Run Docker application
docker run -e ARG0=<SEGMENT_DATA_IPFS_HASH> -e ARG1=<TRANSCODING_OPTIONS> verifier
```
