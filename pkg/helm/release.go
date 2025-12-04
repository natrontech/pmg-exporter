package helm

import (
	"errors"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

var ErrNotFound = errors.New("not found")

const helmDriver = "secret"

type Options struct {
	RepoURL string
}

const defaultTimeout = 5 * time.Minute

func (helm *Client) Install(chart, name, version string, values map[string]interface{}, opts Options) (*release.Release, error) {
	install := action.NewInstall(helm.config)
	install.ReleaseName = name
	install.Wait = true
	install.Atomic = true
	install.Timeout = defaultTimeout
	install.RepoURL = opts.RepoURL

	if version != "" {
		install.Version = version
	}

	chartWithRegistry := chart

	if helm.ociRegistry != nil {
		chartWithRegistry = helm.ociRegistry.JoinPath(chart).String()
	}

	chartPath, err := install.ChartPathOptions.LocateChart(chartWithRegistry, helm.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	vals, err := helm.mergeValues(values)
	if err != nil {
		return nil, err
	}

	release, err := install.Run(chartRequested, vals)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (helm *Client) Uninstall(release string) error {
	uninstall := action.NewUninstall(helm.config)

	if _, err := uninstall.Run(release); err != nil {
		return err
	}

	return nil
}

func (helm *Client) Upgrade(release, chart, version string, values map[string]interface{}, opts Options) (*release.Release, error) {
	upgrade := action.NewUpgrade(helm.config)
	upgrade.Version = version
	upgrade.Wait = true
	upgrade.Atomic = true
	upgrade.Timeout = defaultTimeout

	chartWithRegistry := chart

	if helm.ociRegistry != nil {
		chartWithRegistry = helm.ociRegistry.JoinPath(chart).String()
	}

	chartPath, err := upgrade.ChartPathOptions.LocateChart(chartWithRegistry, helm.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	vals, err := helm.mergeValues(values)
	if err != nil {
		return nil, err
	}

	return upgrade.Run(release, chartRequested, vals)
}

func (helm *Client) Values(release string, allValues bool) (map[string]interface{}, error) {
	getValues := action.NewGetValues(helm.config)
	getValues.AllValues = allValues

	result, err := getValues.Run(release)

	if err != nil {
		if err.Error() == "release: not found" {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if result == nil {
		return nil, ErrNotFound
	}

	return result, nil
}

func (helm *Client) Release(name string) (*Release, error) {
	get := action.NewGet(helm.config)

	result, err := get.Run(name)

	if err != nil {
		if err.Error() == "release: not found" {
			return nil, ErrNotFound
		}

		return nil, err
	}

	if result == nil {
		return nil, ErrNotFound
	}

	return result, nil
}

func (helm *Client) Releases() ([]Release, error) {
	list := action.NewList(helm.config)
	list.All = true
	list.SetStateMask()

	results, err := list.Run()

	if err != nil {
		return nil, err
	}

	var releases []Release

	for _, rel := range results {
		releases = append(releases, *rel)
	}

	return releases, nil
}

func (helm *Client) Exists(release string) (bool, error) {
	if _, err := helm.Release(release); err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (helm *Client) Rollback(release string) error {
	rollback := action.NewRollback(helm.config)
	rollback.Wait = true

	return rollback.Run(release)
}

func (helm *Client) History(release string) ([]*release.Release, error) {
	history := action.NewHistory(helm.config)
	return history.Run(release)
}

func (helm *Client) Status(release string) (*release.Release, error) {
	status := action.NewStatus(helm.config)
	return status.Run(release)
}
