# Security Policy

## Supported Versions

This repository is a bootstrap architecture template. Security fixes are applied to the default branch.

| Version | Supported |
| ------- | --------- |
| `main`  | Yes       |

## Reporting a Vulnerability

Please do not report security vulnerabilities in public issues.

Report suspected vulnerabilities by using GitHub private vulnerability reporting if enabled for this repository, or by contacting the repository maintainers directly.

Include:

- Affected files, commands, or workflows
- Reproduction steps or proof of concept
- Expected impact
- Suggested remediation, if known

Maintainers should acknowledge reports within 5 business days and provide a remediation plan or status update as soon as practical.

## GitHub Actions Security

Workflows should follow least-privilege defaults:

- Set `permissions:` explicitly.
- Avoid `pull_request_target` unless privileged behavior is required and reviewed.
- Do not print secrets or sensitive environment data.
- Prefer short-lived credentials and OIDC over long-lived cloud secrets.
- Keep Actions dependencies updated with Dependabot.

## Secret Scanning

Gitleaks scans the working tree in the pre-commit hook and scans repository
history in CI. Run `task secrets` before committing when the hook is not
installed. Run `task secrets-history` for an initial audit or when responding
to a suspected credential exposure.

If a secret is detected, remove it, revoke or rotate the credential, and assess
whether the Git history also needs remediation before continuing.
