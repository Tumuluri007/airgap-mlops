---
name: Bug Report
about: Report a bug in AGM
title: "[bug] "
labels: bug
---

## Description

A clear, concise description of the bug.

## Component

- [ ] Reference architecture / Helm chart
- [ ] Kyverno policy (which: ___)
- [ ] CMBA admission webhook
- [ ] CMBA init container (cmba-verify)
- [ ] CMBA sentinel sidecar
- [ ] Transfer bundler scripts
- [ ] CI/CD workflows
- [ ] Documentation

## Reproduction Steps

1. ...
2. ...
3. ...

## Expected Behavior

What you expected to happen.

## Actual Behavior

What actually happened, including any error messages.

## Environment

- AGM version:
- Kubernetes version:
- Helm version:
- OS:
- Cluster type (kind, minikube, production):

## Logs

```
kubectl logs -n cmba-system -l app=cmba-webhook --tail=100
```

## Additional Context

Anything else relevant: configuration, related changes, screenshots.
