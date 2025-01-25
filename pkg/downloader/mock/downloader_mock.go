package downloader

import (
	"fmt"

	"github.com/helmless/helmless-cli/pkg/downloader"
	"helm.sh/helm/v3/pkg/chart"
)

type mockDependencyManager struct{}

func NewMockDependencyManager() downloader.DependencyManager {
	return &mockDependencyManager{}
}

func (m *mockDependencyManager) Update(chart *chart.Chart) error {
	fmt.Println("Updating dependencies for", chart.Name())
	return nil
}
