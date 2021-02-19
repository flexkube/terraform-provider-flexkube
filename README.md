# Flexkube Terraform Provider [![Build Status](https://travis-ci.com/flexkube/terraform-provider-flexkube.svg?branch=master)](https://travis-ci.com/flexkube/terraform-provider-flexkube) [![Maintainability](https://api.codeclimate.com/v1/badges/bc27e5bc32a9b40efaa0/maintainability)](https://codeclimate.com/github/flexkube/terraform-provider-flexkube/maintainability) [![codecov](https://codecov.io/gh/flexkube/terraform-provider-flexkube/branch/master/graph/badge.svg)](https://codecov.io/gh/flexkube/terraform-provider-flexkube) [![Go Report Card](https://goreportcard.com/badge/github.com/flexkube/terraform-provider-flexkube)](https://goreportcard.com/report/github.com/flexkube/terraform-provider-flexkube)

The Flexkube provider allows to create and manage Kubernetes cluster components using [libflexkube](https://github.com/flexkube/libflexkube). With this provider, you can create containers like etcd or kubelet on remote machines over SSH using Docker container runtime.

This provider also provides `flexkube_helm_release` resource, so you can use it to manage cluster-essential workloads like CNI plugins or [CoreDNS](https://coredns.io/).

## Table of contents
* [User documentation](#user-documentation)
* [Building and testing](#building-and-testing)
* [Authors](#authors)

## User documentation

For user documentation, see [Terraform Registry](https://registry.terraform.io/providers/flexkube/flexkube/latest/docs).

## Building

For local builds, run `make` which will build the binary, run unit tests and linter.

## Releasing

This project use `goreleaser` for releasing. To release new version, follow the following steps:

* Add a changelog for new release to CHANGELOG.md file.

* Update `docs/index.md` file to update recommended version in the example.

* Tag new release on desired git commit, using example command:

  ```sh
  git tag -a v0.4.7 -s -m "Release v0.4.7"
  ```

* Push the tag to GitHub
  ```sh
  git push origin v0.4.7
  ```

* Run `goreleser` to create a GitHub Release:
  ```sh
  GITHUB_TOKEN=githubtoken GPG_FINGERPRINT=gpgfingerprint goreleaser release --release-notes <(go run github.com/rcmachado/changelog show 0.4.7)
  goreleaser
  ```

* Go to newly create [GitHub release](https://github.com/flexkube/terraform-provider-flexkube/releases/tag/v0.4.7), verify that the changelog
  and artefacts looks correct and publish it.

## Authors

* **Mateusz Gozdek** - *Initial work* - [invidian](https://github.com/invidian)
