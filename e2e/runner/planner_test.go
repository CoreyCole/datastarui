package runner

import (
	"reflect"
	"testing"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func TestPlanJobsChangedFiles(t *testing.T) {
	cfg := appconfig.Config{RunPackage: "./components/..."}
	tests := []struct {
		name    string
		changed []ChangedFile
		want    []E2EJob
	}{
		{
			name:    "component source maps to component package",
			changed: []ChangedFile{{Status: "M", Path: "components/select/select.templ"}},
			want:    []E2EJob{{ID: "select", Package: "./components/select", Component: "select", Reason: "changed components/select/select.templ"}},
		},
		{
			name:    "generated component source maps to component package",
			changed: []ChangedFile{{Status: "M", Path: "components/select/select_templ.go"}},
			want:    []E2EJob{{ID: "select", Package: "./components/select", Component: "select", Reason: "changed components/select/select_templ.go"}},
		},
		{
			name:    "component page maps to component package",
			changed: []ChangedFile{{Status: "M", Path: "pages/components/selectpage/select_page.templ"}},
			want:    []E2EJob{{ID: "select", Package: "./components/select", Component: "select", Reason: "changed pages/components/selectpage/select_page.templ"}},
		},
		{
			name:    "shared e2e path maps to full package",
			changed: []ChangedFile{{Status: "M", Path: "e2e/runtime/scenario.go"}},
			want:    []E2EJob{{ID: "all", Package: "./components/...", Reason: "shared or unclassified change"}},
		},
		{
			name:    "command path maps to full package",
			changed: []ChangedFile{{Status: "M", Path: "cmd/datastarui/e2e.go"}},
			want:    []E2EJob{{ID: "all", Package: "./components/...", Reason: "shared or unclassified change"}},
		},
		{
			name:    "layout path maps to full package",
			changed: []ChangedFile{{Status: "M", Path: "layouts/sidebar.go"}},
			want:    []E2EJob{{ID: "all", Package: "./components/...", Reason: "shared or unclassified change"}},
		},
		{
			name:    "unknown path maps to full package",
			changed: []ChangedFile{{Status: "M", Path: "README.md"}},
			want:    []E2EJob{{ID: "all", Package: "./components/...", Reason: "shared or unclassified change"}},
		},
		{
			name:    "rename maps old and new components",
			changed: []ChangedFile{{Status: "R", OldPath: "components/select/foo.go", Path: "components/dropdown/foo.go"}},
			want: []E2EJob{
				{ID: "dropdown", Package: "./components/dropdown", Component: "dropdown", Reason: "changed components/dropdown/foo.go"},
				{ID: "select", Package: "./components/select", Component: "select", Reason: "renamed from components/select/foo.go"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PlanJobs(cfg, RunOptions{}, tt.changed)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("jobs = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPlanJobsAllAndExplicitFilters(t *testing.T) {
	cfg := appconfig.Config{RunPackage: "./components/..."}
	all, err := PlanJobs(cfg, RunOptions{All: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(all, []E2EJob{{ID: "all", Package: "./components/...", Reason: "all requested"}}) {
		t.Fatalf("all jobs = %#v", all)
	}

	filtered, err := PlanJobs(cfg, RunOptions{Story: "select component opens options", Scenario: "desktop full"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []E2EJob{{ID: "all", Package: "./components/...", RunPattern: "SelectComponentOpensOptions.*DesktopFull", Reason: "explicit story/scenario filter"}}
	if !reflect.DeepEqual(filtered, want) {
		t.Fatalf("filtered jobs = %#v, want %#v", filtered, want)
	}
}

func TestParseNameStatus(t *testing.T) {
	input := "M\tcomponents/select/select.templ\nR100\tcomponents/select/old.go\tcomponents/dropdown/new.go\n"
	got := parseNameStatus(input)
	want := []ChangedFile{
		{Status: "M", Path: "components/select/select.templ"},
		{Status: "R", OldPath: "components/select/old.go", Path: "components/dropdown/new.go"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parse = %#v, want %#v", got, want)
	}
}

func TestParseUntrackedAndDedupeChangedFiles(t *testing.T) {
	got := dedupeChangedFiles(append(parseUntracked("b.go\na.go\n"), ChangedFile{Status: "A", Path: "a.go"}))
	want := []ChangedFile{{Status: "A", Path: "a.go"}, {Status: "A", Path: "b.go"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("dedupe = %#v, want %#v", got, want)
	}
}
