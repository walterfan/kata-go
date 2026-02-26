# CVE validation tests for crypto/tls

This package contains unit tests that verify fixes for known CVEs in Go's standard `crypto/tls` package. They are intended to be run with a Go toolchain that includes the fix; the test **passes** when the CVE is fixed and would **fail** on an unfixed Go version.

## CVE-2025-68121

**References:** [Go issue #77113](https://github.com/golang/go/issues/77113), [blog (session resumption pitfalls)](https://www.fanyamin.com/blog/go-cryptotls-configclone-session-resumption-pitfalls-cve-2025-68121.html)

**Issue:** `Config.Clone()` copied automatically generated session ticket keys, so two configs could share keys and allow session resumption across what should be isolated boundaries. Session resumption also did not consider full certificate chain expiry (leaf-only check).

**Fix (Go 1.25.7+):** `Clone()` no longer copies auto-generated ticket keys; only explicitly set keys are copied. Resumption checks the full certificate chain for expiry.

### How to run

```bash
# From repo root
go test cve_2025_68121_test.go -v -run TestCVE_2025_68121
```
### Interpretation

- **PASS:** The Go toolchain in use includes the fix; cloned configs do not share auto-generated session ticket keys, so the client does not resume when connecting to the cloned server.
- **FAIL (DidResume true):** The binary was built with an unfixed Go version; upgrade to Go 1.25.7 or later (or the version that backports the fix for your release line).
