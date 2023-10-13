# CertManager Component Operator

[![REUSE status](https://api.reuse.software/badge/github.com/SAP/cert-manager-cop)](https://api.reuse.software/info/github.com/SAP/cert-manager-cop)

## About this project

Component Operator for [https://cert-manager.io](https://cert-manager.io).

After installing this component operator into a Kubernetes cluster, cert-manager may be installed by deploying

```yaml
apiVersion: operator.kyma-project.io/v1alpha1
kind: CertManager
metadata:
  name: cert-manager
# spec:
  # optional spec attributes
```

In `spec`, all values of the [upstream Helm chart](https://github.com/cert-manager/cert-manager/tree/master/deploy/charts/cert-manager) are allowed. Caveats:
- the component operator does not perform any validation, it just passes the provided spec as values to the helm chart
- the supported/allowed spec format might change, when the included upstream chart changes
- the helm chart has an option to skip deployment of custom resource definitions; the component operator will forcefully overwrite the according switch (`Values.installCRDs`) to be always true
- deploying multiple `CertManager` resources in a cluster will not work.

In addition, the following attributes can be supplied in `spec`:
- `namespace`: target namespace for the cert-manager (if not specified, the namespace of the owning `CertManager` resource will be used)
- `name`: target name for the deployed cert-manager (if not specified, generated resources will be prefixed with the name of the owning `CertManager` resource)
- `additionalResources`: array of additional resource manifests that will be deployed along with the cert-manager.

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/cert-manager-cop/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2023 SAP SE or an SAP affiliate company and cert-manager-cop contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/cert-manager-cop).
