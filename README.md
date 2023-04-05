# AWS CodeArtifact proxy for Terraform
This is a relatively simple project which allows you to use AWS CodeArtifact as
registry for Terraform modules acting as a bridge. It exposes an API compatible
with the registry format.


# Installation
An example terraform deployment setup is available in the `deployment/` folder.
There are three environment variables which need to be set:
  - `REGISTRY_NAME`: the name of your CodeArtifact registry
  - `REGISTRY_DOMAIN`: the domain of your CodeArtifact registry
  - `AUTH_TOKENS`: a comma separate list of tokens for authentication

# Uploading new terraform modules
Since this is using AWS CodeArtifact uploading new modules happens with the aws
cli. See the [documentation](https://docs.aws.amazon.com/codeartifact/latest/ug/publishing-using-generic-packages.html)
fore more info.

Example:
```
aws codeartifact publish-package-version \
  --domain your-domain \
  --repository your-repository \
  --namespace your-department \
  --format generic \
  --package terraform-aws-vpc \
  --package-version 1.2.1 \
  --asset-name terraform-aws-vpc-1.2.1.zip \
  --asset-content terraform-aws-vpc-1.2.1.zip \
  --asset-sha256 $(sha256sum asset.tar.gz | awk '{print $1;}')
```

# Contributing
Contributions are welcome. Please open an issue or a PR.
