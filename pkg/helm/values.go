package helm

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	helmValues "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

func LookupValues[T any](m map[string]interface{}, key string) (v T, found bool) {
	parts := strings.Split(key, ".")
	current := m

	for i, part := range parts {
		val, ok := current[part]
		if !ok {
			return v, false
		}

		isLast := i == len(parts)-1
		if isLast {
			// For the last part, try to convert to the target type
			if typedVal, ok := val.(T); ok {
				return typedVal, true
			}
			return v, false
		}

		// Not the last part, so we expect this to be a map
		nextMap, ok := val.(map[string]interface{})
		if !ok {
			return v, false
		}
		current = nextMap
	}

	return v, false
}

func (helm *Client) mergeValues(values map[string]interface{}) (map[string]interface{}, error) {
	valuesBytes, err := yaml.Marshal(values)
	if err != nil {
		return nil, err
	}

	// Write values to a temporary file
	tempFile, err := os.CreateTemp("", "values-*.yaml")
	if err != nil {
		return nil, err
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	if _, err := tempFile.Write(valuesBytes); err != nil {
		return nil, err
	}

	// Create a values.Options with the temporary file path
	valueOpts := &helmValues.Options{
		ValueFiles: []string{tempFile.Name()},
	}

	return valueOpts.MergeValues(getter.All(helm.settings))
}
