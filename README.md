[![Build](https://github.com/jidicula/vci-check/actions/workflows/build.yml/badge.svg)](https://github.com/jidicula/vci-check/actions/workflows/build.yml) [![Latest Release](https://github.com/jidicula/vci-check/actions/workflows/release-draft.yml/badge.svg)](https://github.com/jidicula/vci-check/actions/workflows/release-draft.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/jidicula/vci-check)](https://goreportcard.com/report/github.com/jidicula/vci-check) [![Go Reference](https://pkg.go.dev/badge/github.com/jidicula/vci-check.svg)](https://pkg.go.dev/github.com/jidicula/vci-check)

# vci-check

`vci-check` is a simple web API for checking if a vaccine credential issuer is indeed a trusted issuer. It accepts a GET query at the root endpoint with the issuer URL as a query parameter:

`curl host/?iss=https://totallylegitissuer/creds`

The full list of trusted issuers is available in the [VCI Directory](https://github.com/the-commons-project/vci-directory).

## Deploying

You can deploy it via Docker CLI with `docker run -p 8080:8080 ghcr.io/jidicula/vci-check:latest`

# vci-check/checker

This repo also provides the Go package `checker`, containing types corresponding to the VCI directory and functions for working with them.
