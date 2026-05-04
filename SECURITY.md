# Security Policy

## Supported Versions

| Version | Supported |
|---|---|
| 0.1.x | Yes (current) |

## Reporting a Vulnerability

If you discover a security vulnerability in AGM, please report it privately.
Do **not** open a public GitHub issue.

**Preferred reporting channel**:
- GitHub private vulnerability reporting: use the "Report a vulnerability"
  button on the [Security tab](https://github.com/Tumuluri007/airgap-mlops/security)
  of this repository.

**What to include in your report**:
- A clear description of the vulnerability and its impact.
- Steps to reproduce, including any required configuration.
- The component affected (Kyverno policy ID, CMBA webhook/sentinel/init,
  Helm chart, or transfer-bundler).
- The version of AGM you tested against.
- Any suggested mitigation, if you have one.

## Response Timeline

- **Acknowledgement**: within 5 business days.
- **Triage and severity assessment**: within 10 business days.
- **Fix or mitigation**: timeline depends on severity.
  - Critical: target resolution within 14 days.
  - High: target resolution within 30 days.
  - Medium: target resolution within 90 days.
  - Low: addressed in the next minor release.

## Disclosure Policy

We follow coordinated disclosure. Once a fix is available, we will publish a
security advisory on the repository. Reporters will be credited unless they
request anonymity.

## Scope

In scope:
- All code in this repository.
- The Helm chart and its templates.
- The CMBA admission webhook, init container, and sentinel sidecar.
- The Kyverno policy library.
- The transfer-bundler shell scripts.

Out of scope:
- Vulnerabilities in upstream dependencies (Kyverno, Sigstore, KServe, MLflow,
  Harbor, Argo Workflows) should be reported to those projects directly.
- Vulnerabilities in third-party Helm dependencies should be reported to the
  respective maintainers.

## Air-Gap Threat Model Notes

AGM is designed for air-gapped clusters where the threat surface differs from
internet-connected MLOps platforms. Of particular interest:
- **Transfer bundle tampering**: physical media can be modified between the
  external pipeline and internal pipeline. The signed manifest is the primary
  defense; signing key compromise is in scope.
- **Model file replacement on shared volumes**: this is the failure mode CMBA
  exists to detect. Bypass techniques (e.g., disabling the sentinel sidecar)
  are in scope.
- **Insider threats**: an operator with cluster admin rights can disable
  policies and CMBA. Defense-in-depth recommendations are documented in
  [`docs/05-cmba-internals.md`](docs/05-cmba-internals.md).
