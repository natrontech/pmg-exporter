# CLAUDE.md

## Updating Python version and dependencies

When updating to a new Python version or upgrading dependencies, touch all of the following locations:

### 1. Update dependencies in pyproject.toml

Edit [pyproject.toml](pyproject.toml) and bump pinned versions under `dependencies`:

```toml
dependencies = [
    "prometheus-client==0.x.y",
    "proxmoxer==2.x.y",
    "python-dotenv==1.x.y",
    "requests==2.x.y",
]
```

### 2. Update the Python version constraint in pyproject.toml

Edit [pyproject.toml](pyproject.toml) and bump `requires-python`:

```toml
requires-python = ">=3.11"
```

Also update the `classifiers` list to reflect the supported versions:

```toml
classifiers = [
    "Programming Language :: Python :: 3.11",
    ...
]
```

### 3. Update hardcoded Python version in workflows

Two workflows hardcode the Python version and must be updated manually:

- [.github/workflows/release.yml](.github/workflows/release.yml) — `python-version: "3.12"`
- [.github/workflows/bandit.yml](.github/workflows/bandit.yml) — `python-version: "3.12"`

The following workflow uses `python-version: '3.x'` and picks up any available version automatically — no changes needed:

- [.github/workflows/python-lint.yml](.github/workflows/python-lint.yml)

### 4. Update GitHub Actions versions

All workflow files under [.github/workflows/](.github/workflows/) pin actions by commit SHA with a tag comment, e.g.:

```yaml
uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
```

To update, get the latest tag and its commit SHA for each action:

```bash
# Get latest tag
gh release view --repo <owner>/<repo> --json tagName -q '.tagName'

# Get commit SHA for that tag
gh api repos/<owner>/<repo>/commits/<tag> --jq '.sha'
```

Then update the SHA and the tag comment in the workflow file.

**Exception — `slsa-framework/slsa-github-generator`** must be referenced by tag, not SHA (see [upstream docs](https://github.com/slsa-framework/slsa-github-generator/?tab=readme-ov-file#referencing-slsa-builders-and-generators)):

```yaml
uses: slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0
```

Actions used across the workflows:

| Action | Workflow(s) |
|---|---|
| `actions/checkout` | all |
| `actions/dependency-review-action` | dependency-review.yml |
| `actions/setup-python` | release.yml, python-lint.yml, bandit.yml |
| `actions/upload-artifact` | release.yml, scorecard.yml |
| `actions/download-artifact` | release.yml |
| `anchore/sbom-action` | release.yml |
| `docker/build-push-action` | release.yml |
| `docker/login-action` | release.yml |
| `docker/metadata-action` | release.yml |
| `docker/setup-buildx-action` | release.yml |
| `github/codeql-action` | codeql.yml, scorecard.yml |
| `google/osv-scanner-action` | osv-scan.yml |
| `ossf/scorecard-action` | scorecard.yml |
| `sigstore/cosign-installer` | release.yml, release-verification.yml |
| `slsa-framework/slsa-github-generator` | release.yml (**tag only**) |
| `slsa-framework/slsa-verifier` | release-verification.yml |
| `softprops/action-gh-release` | release.yml |

### 5. Update pre-commit hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) pins the `rev` of each hook repository. Update all revisions to their latest tags:

```bash
pre-commit autoupdate
```

This updates the `rev` fields for all repos in [.pre-commit-config.yaml](.pre-commit-config.yaml):
- `pre-commit/pre-commit-hooks`
- `gitleaks/gitleaks`
- `psf/black`
- `pycqa/flake8`
- `RobertCraigie/pyright-python`

### 6. Verify

```bash
pip install hatch
hatch build
```
