# AWS CodeArtifact proxy for Terraform
This is a relatively simple project which allows you to use AWS CodeArtifact as
registry for Terraform moduels acting as a bridge. It exposes an API compatible
with the registry format.


# Installation
An example terraform deployment setup is available in the `deployment/` folder.
There are three environment variables which need to be set:
  - `REGISTRY_NAME`: the name of your CodeArtifact registry
  - `REGISTRY_DOMAIN`: the domain of your CodeArtifact registry
  - `AUTH_TOKENS`: a comma separate list of tokens for authentication


# Contributing
Contributions are welcome. Please open an issue or a PR.
