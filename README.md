
`Dockerfile`:
```
FROM scratch

ARG BINARY

COPY ${BINARY} /${BINARY}
```

Initialize `buildkit`:

```
export DOCKER_BUILDKIT=1
```

Build images:

**Note**: Randomly generate data to avoid common SHAs or any other weird issue:
```
dd if=/dev/urandom of=ota-esp32.bin count=7 bs=1M
dd if=/dev/urandom of=ota-esp32s2.bin count=4 bs=1M
```

```
docker buildx build --platform custom/esp32 -t harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0-esp32 --build-arg BINARY=ota-esp32.bin . --push --provenance false
```

```
docker buildx build --platform custom/esp32 -t harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0-esp32s2 --build-arg BINARY=ota-esp32s2.bin . --push --provenance false
```

Build and push manifest:

```
docker manifest create harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0 \
	--amend harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0-esp32s2 
	--amend harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0-esp32
```

Build tool to inspect metadata / fetch binaries:

```
go build
```

```
$ ./oci-esp32 -image harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0 --arch esp32 --os custom
Image Digest: sha256:b004348524d5009509ca4029cd60dcd52f20a7c74b45b5bac23e7abaab6bf2b3
Architecture: esp32
OS: custom
Created: 2024-10-10T13:34:14Z
Environment Variables:
  PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
Layer 1 - Size: 12586901 bytes, DiffID: sha256:47ffc196cb73e9d4792cdebd9aa5962a8ee54021168a487216a39bf04254e75b
Processing Layer 1...
Extracted: extracted_files/ota-esp32.bin
Extraction complete!
```

```
$ ./oci-esp32 -image harbor.nbfc.io/nubificus/test/esp32-firmware:1.1.0 --arch esp32s2 --os custom
Image Digest: sha256:3d7c9aa1022e706315247ead7ef4a80d9e66cc6bb5f5747db1e5eef9287b31fe
Architecture: esp32s2
OS: custom
Created: 2024-10-10T13:55:29Z
Environment Variables:
  PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
Layer 1 - Size: 12586903 bytes, DiffID: sha256:51c9ff27879e8a241d248083b65064140e8018827d75b361105e38db7f43aa9c
Processing Layer 1...
Extracted: extracted_files/ota-esp32s2.bin
Extraction complete!

