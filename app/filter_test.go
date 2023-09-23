package app

import (
	"strings"
	"testing"
)

func TestFilter_doesFileNameMatchFilters(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		inUpdatedFile string
		inParams      Params
		want          bool
	}{
		"if filters is empty, return true": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams:      Params{},
			want:          true,
		},
		"if filters is set, but the name of updated file does not match filters, return false": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams: Params{
				Filters: []string{"*.txt", "*.md", "*.yaml"},
			},
			want: false,
		},
		"if filters is set, and the name of updated file matchs the filter, return true": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams: Params{
				Filters: []string{"*.go"},
			},
			want: true,
		},
		"if filters is set, and the name of updated file matchs the first filter, return true": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams: Params{
				Filters: []string{"*.go", "*.yaml"},
			},
			want: true,
		},
		"if filters is set, and the name of updated file matchs the last filter, return true": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams: Params{
				Filters: []string{"*.yaml", "*.go"},
			},
			want: true,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := doesFileNameMatchFilters(tc.inUpdatedFile, tc.inParams)
			if err != nil {
				t.Fatalf("expect no error, but got %+v", err)
			}

			if got != tc.want {
				t.Errorf("expect %t, but got %t", tc.want, got)
			}
		})
	}
}

func TestFilter_doesFileNameMatchFilters_error(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		inUpdatedFile string
		inParams      Params
		wantErrMsg    string
	}{
		"if bad pattern is specified to filter, returns error": {
			inUpdatedFile: "/Users/fuga/watcher/app/watcher.go",
			inParams: Params{
				Filters: []string{"\\"},
			},
			wantErrMsg: "failed to match file name",
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, gotErr := doesFileNameMatchFilters(tc.inUpdatedFile, tc.inParams)
			if gotErr == nil {
				t.Fatal("expect error, but got no error")
			}

			if !strings.Contains(gotErr.Error(), tc.wantErrMsg) {
				t.Errorf("expect %v, but got %v", tc.wantErrMsg, gotErr.Error())
			}
		})
	}
}
