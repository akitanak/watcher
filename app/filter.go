package app

import (
	"fmt"
	"path/filepath"
)

func doesFileNameMatchFilters(updatedFile string, params Params) (bool, error) {
	if len(params.Filters) == 0 {
		return true, nil
	}

	for _, filter := range params.Filters {
		matched, err := filepath.Match(filter, filepath.Base(updatedFile))
		if err != nil {
			return false, fmt.Errorf("failed to match file name: %w", err)
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}
