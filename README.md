# SGX Attestation App
This is a sample application to demostrate the SGX enclave attestation with Intel® Trust Authority using [trustauthority-client](https://github.com/intel/trustauthority-client-for-go/)

## Instructions for Ubuntu
Use the below command to install all the dependencies necessary to run this sample app.

```sh
$ sudo ./install.sh
```

## Build
Once the above install script is done, use the below steps to build and run the app

```sh
make -C enclave
```
```sh
CGO_CFLAGS_ALLOW="-f.*" /usr/local/go/bin/go build
```
```sh
LD_LIBRARY_PATH=enclave/ ./sgxexample --config config.json
```

## Config Definition
```json
{
    "trustauthority_url": "<trustauthority url>",
    "trustauthority_api_url": "<trustauthority attestation api url>",
    "trustauthority_api_key": "<trustauthority attestation api key>"
}
```
