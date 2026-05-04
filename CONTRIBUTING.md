# Contributing to AGM

Thank you for your interest in contributing. This guide covers how to add new
Kyverno policies, extend the CMBA verifier, run tests in a kind cluster, and
the code review process.

## Code of Conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).
By participating you agree to uphold its terms.

## How to Add a New Kyverno Policy

1. Create a new file under the appropriate group folder:
   - `policies/airgap/` for air-gap enforcement
   - `policies/supply-chain/` for supply chain integrity
   - `policies/cmba/` for CMBA-related enforcement
   - `policies/governance/` for general ML governance
2. Number the policy sequentially (e.g., `13-your-policy-name.yaml`).
3. Set `validationFailureAction: Enforce` for production policies. Use `Audit`
   only for policies that need a soak period.
4. Add a chainsaw test under `policies/<group>/tests/<policy-slug>/`. Every
   policy must have at least one passing and one failing test case.
5. Document the policy in `docs/04-policy-reference.md` with: what it enforces,
   what it blocks, why it is ML-specific, and how it operates without internet.
6. Run `make policy-test` locally before opening a pull request.

## How to Extend the CMBA Verifier

### Adding a new model format

1. Update the `format` enum in `cmba/crd/modelbinding-crd.yaml` to include
   the new format identifier.
2. If the new format requires a different hashing strategy (e.g., directory
   hashing for multi-file models), extend `cmba/init-container/verify.sh` and
   `cmba/sentinel/verify.go`.
3. Add an example ModelBinding in `cmba/examples/`.
4. Add test cases in `cmba/webhook/handler_test.go` and `cmba/sentinel/verify_test.go`.

### Adding a new hash algorithm

1. Add the algorithm option to the CRD spec (e.g., `spec.modelArtifact.hashAlgo`).
2. Update both verifier and sentinel to support the new algorithm.
3. Default remains SHA-256 for backwards compatibility.

## How to Test in a kind Cluster

```bash
# Create a local kind cluster
kind create cluster --name agm-test

# Install the entire platform
helm install airgap-mlops ./helm/airgap-mlops \
  -f ./helm/airgap-mlops/values-airgap.yaml

# Run policy tests
make policy-test

# Run CMBA component tests
make cmba-test

# Run end-to-end walkthrough
./examples/e2e-walkthrough/1-build-external.sh
./examples/e2e-walkthrough/2-transfer.sh
./examples/e2e-walkthrough/3-import-internal.sh
./examples/e2e-walkthrough/4-deploy-and-verify.sh

# Tear down
kind delete cluster --name agm-test
```

## Code Review Process

1. Open a pull request against `main`.
2. CI must pass: policy validation, Helm chart linting, CMBA component tests.
3. At least one maintainer review is required.
4. For policy additions: a security reviewer must explicitly approve.
5. For CMBA changes: a maintainer with Go and Kubernetes admission webhook
   experience must approve.
6. Squash-merge by default. Preserve commit messages with the rationale.

## Commit Message Format

We follow Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`.

Example:
```
feat(policies): add policy 13 for GPU memory bounds

New ClusterPolicy enforces a per-pod GPU memory ceiling on ml-serving
pods to prevent runaway inference jobs from starving co-located workloads.

Closes #42
```

## Reporting Issues

Use the issue templates in `.github/ISSUE_TEMPLATE/`. Security issues should
be reported privately per [`SECURITY.md`](SECURITY.md), not via public issues.
