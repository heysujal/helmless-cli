package downloader

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

type DependencyManager interface {
	Update(chart *chart.Chart) error
}

type Options struct {
	Fs       afero.Fs
	Settings *cli.EnvSettings
}

type dependencyManager struct {
	fs       afero.Fs
	settings *cli.EnvSettings
}

func New() *dependencyManager {
	return NewWithOptions(Options{})
}

func NewWithOptions(opts Options) *dependencyManager {
	if opts.Settings == nil {
		opts.Settings = cli.New()
	}
	if opts.Fs == nil {
		opts.Fs = afero.NewOsFs()
	}
	return &dependencyManager{
		fs:       opts.Fs,
		settings: opts.Settings,
	}
}

func (m *dependencyManager) Update(chart *chart.Chart) error {
	fmt.Fprintf(os.Stderr, "Updating dependencies for %s...\n", chart.Name())

	man := &downloader.Manager{
		Out:              os.Stderr, // Show download progress
		ChartPath:        chart.ChartFullPath(),
		Getters:          getter.All(m.settings),
		RepositoryConfig: m.settings.RepositoryConfig,
		RepositoryCache:  m.settings.RepositoryCache,
		Debug:            m.settings.Debug,
	}

	if err := man.Update(); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Dependencies updated successfully\n")
	return nil
}
