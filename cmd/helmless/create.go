package helmless

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/helmless/helmless-cli/pkg/downloader"
	"github.com/manifoldco/promptui"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

type createOptions struct {
	name string
	fs   afero.Fs
	dm   downloader.DependencyManager
}

func newCreateCmd() *cobra.Command {
	return newCreateCmdWithOptions(createOptions{
		fs: afero.NewOsFs(),
		dm: downloader.New(),
	})
}

func newCreateCmdWithOptions(opts createOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new Helm chart",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				prompt := promptui.Prompt{
					Label: "Chart name",
				}
				name, err := prompt.Run()
				if err != nil {
					return fmt.Errorf("failed to get chart name: %w", err)
				}
				opts.name = name
			} else {
				opts.name = args[0]
			}

			return opts.run()
		},
	}

	return cmd
}

func (o *createOptions) run() error {
	fmt.Printf("Creating new Helm chart '%s'...\n", o.name)

	// Create chart directory
	if err := o.fs.MkdirAll(o.name, 0755); err != nil {
		return fmt.Errorf("failed to create chart directory: %w", err)
	}

	chartDir := o.name

	fmt.Println("⚙️  Generating chart files...")
	// Create Chart.yaml
	metadata := &chart.Metadata{
		APIVersion: "v2",
		Name:       o.name,
		Version:    "0.1.0",
		AppVersion: "0.1.0",
		Type:       "application",
		Dependencies: []*chart.Dependency{
			{
				Name:       "google-cloudrun-service",
				Version:    "0.1.0",
				Repository: "oci://ghcr.io/helmless",
				Alias:      "app",
			},
		},
	}

	err := chartutil.CreateFrom(metadata, ".", "bootstrap-chart")
	if err != nil {
		return fmt.Errorf("failed to create chart: %w", err)
	}

	chart, err := loader.Load(chartDir)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	if err := o.dm.Update(chart); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}

	// Update values.yaml with the app name
	valuesPath := filepath.Join(o.name, "values.yaml")
	values, err := afero.ReadFile(o.fs, valuesPath)
	if err != nil {
		return fmt.Errorf("failed to read values.yaml: %w", err)
	}

	// Replace app name in values.yaml
	newValues := strings.Replace(string(values), "my-helmless-app", o.name, 1)

	if err := afero.WriteFile(o.fs, valuesPath, []byte(newValues), 0644); err != nil {
		return fmt.Errorf("failed to write values.yaml: %w", err)
	}

	fmt.Printf("✨ Successfully created '%s' chart\n", o.name)
	return nil
}
