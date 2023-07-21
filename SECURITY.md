# Security Policy

We take security vulnerabilities seriously (and so should you!)

Our policy on reported vulnerabilities (see below on how to report) is that we will
respond to the reporter of a vulnerability within two (2) business days of receiving
the report and notify the reporter whether and when a remediation will be committed.

When a remediation for a security vulnerability is committed, we will cut a tagged
release of `ghw` and include in the release notes for that tagged release a description
of the vulnerability and a discussion of how it was remediated, along with a note
urging users to update to that fixed version.

## Reporting a Vulnerability

While `ghw` does have automated Github Dependabot alerts about security vulnerabilities
in `ghw`'s dependencies, there is always a chance that a vulnerability in a dependency
goes undetected by Dependabot. If you are aware of a vulnerability either in `ghw` or
one of its dependencies, please do not hesitate to reach out to `ghw` maintainers via
email or Slack. **Do not discuss vulnerabilities in a public forum**.

`ghw`'s primary maintainer is Jay Pipes, who can be found on the Kubernetes Slack
community as `@jaypipes` and reached via email at jaypipes at gmail dot com.
