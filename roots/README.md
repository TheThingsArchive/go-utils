# Roots

Go can't get the system roots on Windows, so we need to add them manually. For this, we use the Root CAs that are bundled with Mozilla Firefox through the `ca-certificate` package in Alpine Linux.

## Updating the certificates

_This requires `docker`._

```
go generate
```
