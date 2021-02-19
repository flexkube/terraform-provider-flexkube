# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0] - 2021-02-18

### Added

- `flexkube_kubelet` has now `extra_args` argument for extra kubelet arguments.

### Changed

- Updated to `libflexkube` version `v0.5.1`.
- e2e tests are now using `containerd` as container runtime.

## [0.4.1] - 2020-09-20

### Added

- `flexkube_containers` configuration can have now `env` field specified to set environment variables
  for containers.
- Added configuration for running E2E and local tests.

### Changed

- `libflexkube` has been updated to [v0.4.3](https://github.com/flexkube/libflexkube/releases/tag/v0.4.3).
- Migrated to [Terraform SDK v2](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html).
- `flexkube_pki` resource has been rewritten to follow the patterns from other resources. It should no longer report
  inconsistent plan issues and should no longer show ambiguous changes.
- All resources will now produce less verbose diff when running `terraform plan` and will show that fields will have
  new values computed instead. That should greatly reduce number of issues with inconsistent plan, at the cost of
  changes visibility.

### Fixed

- `flexkube_helm_release` now validates fields `values` and `kubeconfig` to make sure they contain valid YAML
  formatted content, so it does not corrupt the Terraform state, which requires manual interaction to recover from.
- Added missing documentation to `flexkube_containers` resource for `host` block.
- Adding and removing controller nodes should now work without interruption. Previously, it has been broken because of a bug
  in `libflexkube` and `flexkube_pki` and `flexkube_etcd_cluster` resources complaining about the inconsistent plan.

## [0.4.0] - 2020-08-31

### Added

- Initial release based on [libflexkube v0.4.0](https://github.com/flexkube/libflexkube/releases/tag/v0.4.0).

### Changed

- flexkube_containers: rename 'container' to 'host_configured_container'.

[0.5.0]: https://github.com/flexkube/terraform-provider-flexkube/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/flexkube/terraform-provider-flexkube/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/flexkube/terraform-provider-flexkube/releases/tag/v0.4.0
