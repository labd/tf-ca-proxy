package internal

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact/types"
	"github.com/davecgh/go-spew/spew"
)

type Store struct {
	client           *codeartifact.Client
	repositoryName   string
	repositoryDomain string
}

func NewStore() *Store {

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		panic(err)
	}

	client := codeartifact.NewFromConfig(cfg)
	return &Store{
		client:           client,
		repositoryName:   appConfig.RegistryName,
		repositoryDomain: appConfig.RegistryDomain,
	}
}

func (s *Store) listModuleVersions(ctx context.Context, r ModuleRequest) ([]ModuleVersion, error) {
	input := &codeartifact.ListPackageVersionsInput{
		Domain:     aws.String(s.repositoryDomain),
		Repository: aws.String(s.repositoryName),
		Format:     types.PackageFormatGeneric,
		Namespace:  aws.String(r.Namespace),
		Package:    aws.String(r.PackageName()),
	}

	output, err := s.client.ListPackageVersions(ctx, input)
	if err != nil {
		spew.Dump(err)
		return nil, err
	}

	versions := make([]ModuleVersion, 0, len(output.Versions))
	for _, v := range output.Versions {
		versions = append(versions, ModuleVersion{
			Version: *v.Version,
		})
	}

	return versions, nil
}

func (s *Store) getModuleVersion(ctx context.Context, r ModuleRequest) (*ModuleVersionInfo, error) {
	input := &codeartifact.DescribePackageVersionInput{
		Domain:         aws.String(s.repositoryDomain),
		Repository:     aws.String(s.repositoryName),
		Format:         types.PackageFormatGeneric,
		Namespace:      aws.String(r.Namespace),
		Package:        aws.String(r.PackageName()),
		PackageVersion: aws.String(r.Version),
	}

	output, err := s.client.DescribePackageVersion(ctx, input)
	if err != nil {
		spew.Dump(err)
		return nil, err
	}

	result := &ModuleVersionInfo{
		ID:          fmt.Sprintf("%s/%s/%s/%s", r.Namespace, r.Name, r.Provider, r.Version),
		Name:        r.Name,
		Provider:    r.Provider,
		Namespace:   r.Namespace,
		Version:     *output.PackageVersion.Version,
		PublishedAt: output.PackageVersion.PublishedTime,
	}

	return result, nil
}

func (s *Store) getModuleVersionAssets(ctx context.Context, r ModuleRequest) (string, error) {
	inputAssets := &codeartifact.ListPackageVersionAssetsInput{
		Domain:         aws.String(s.repositoryDomain),
		Repository:     aws.String(s.repositoryName),
		Format:         types.PackageFormatGeneric,
		Namespace:      aws.String(r.Namespace),
		Package:        aws.String(r.PackageName()),
		PackageVersion: aws.String(r.Version),
	}

	outputAssets, err := s.client.ListPackageVersionAssets(ctx, inputAssets)
	if err != nil {
		return "", err
	}

	if len(outputAssets.Assets) == 0 {
		return "", fmt.Errorf("No assets found")
	}

	return *outputAssets.Assets[0].Name, nil
}

func (s *Store) downloadModuleVersion(ctx context.Context, r ModuleRequest, asset string) (io.ReadCloser, error) {
	input := &codeartifact.GetPackageVersionAssetInput{
		Domain:         aws.String(s.repositoryDomain),
		Repository:     aws.String(s.repositoryName),
		Format:         types.PackageFormatGeneric,
		Namespace:      aws.String(r.Namespace),
		Package:        aws.String(r.PackageName()),
		PackageVersion: aws.String(r.Version),
		Asset:          aws.String(asset),
	}

	output, err := s.client.GetPackageVersionAsset(ctx, input)
	if err != nil {
		spew.Dump(err)
		return nil, err
	}
	return output.Asset, nil
}
