# Contributing

All contributions are welcome! If you find a bug or have a feature request, please open an issue or submit a pull request.

Please note that we have a [Code of Conduct](./CODE_OF_CONDUCT.md), please follow it in all your interactions with the project.

## How to Contribute

You can make a contribution by following these steps:

  1. Fork this repository, and develop your changes on that fork.
  2. Commit your changes
  3. Submit a [pull request](#pull-requests) from your fork to this project.

Before you start, read through the requirements below.  

### Commits

Please make your commit messages meaningful. We recommend creating commit messages according to [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

### Commit Signature Verification

Each commit's signature must be verified.

  * [About commit signature verification](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/about-commit-signature-verification)

### Pull Requests

All contributions, including those made by project members, need to be reviewed. We use GitHub pull requests for this purpose. See [GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on how to use pull requests. See the requirements above for PR on this project.

### Major new features

If a major new feature is added, there should be new tests for it. If there are no tests, the PR will not be merged.

### Versioning

Versions follow [Semantic Versioning](https://semver.org/) terminology and are expressed as `x.y.z`:

- where `x` is the major version
- `y` is the minor version
- and `z` is the patch version

## Code convention

## Pre-Commit

Please install [pre-commit](https://pre-commit.com/) to enforce some pre-commit checks.
After cloning the repository, you will need to install the hook script manually:

```bash
pre-commit install
```

## Python type checking
We use [pyright](https://pyright.org/) to check Python types.
If you are using VSCode, you can install the [Pyright extension](https://marketplace.visualstudio.com/items?itemName=ms-pyright.pyright) and set it to `strict` mode for better type checking.

## Python dependencies
Use standard library tools whenever possible. Only use third-party tools when absolutely necessary. Your PR will be rejected if it introduces unnecessary dependencies.

Use [pip](https://pip.pypa.io/en/stable/) to manage Python dependencies.
We recommend using [venv](https://docs.python.org/3/library/venv.html) to create a virtual environment for development.
You can create an editable install of the package with:
```bash
python -m pip install -e .
```

## Python code style
We use [black](https://black.readthedocs.io/en/stable/) to format Python code.
You can run it manually with:
```bash
black .
```
We use [flake8](https://flake8.pycqa.org/en/latest/) to check Python code style.
You can run it manually with:
```bash
flake8 .
```

## Reporting Issues
If you encounter any bugs or have feature requests, please open an issue on GitHub. When reporting a bug, please include as much detail as possible to help us understand and reproduce the issue.
When submitting a feature request, please provide a clear description of the desired functionality and its potential benefits.
