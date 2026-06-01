package runner

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/coreycole/datastarui/e2e/appconfig"
)

func PlanJobs(cfg appconfig.Config, opts RunOptions, changed []ChangedFile) ([]E2EJob, error) {
	if opts.All {
		return []E2EJob{fullJob(cfg, "all requested")}, nil
	}
	if strings.TrimSpace(opts.Story) != "" || strings.TrimSpace(opts.Scenario) != "" {
		return []E2EJob{filteredJob(cfg, opts)}, nil
	}

	components := map[string]string{}
	full := len(changed) == 0
	for _, file := range changed {
		component, shared, ok := ClassifyDatastarUIPath(file.Path)
		if file.OldPath != "" {
			oldComponent, oldShared, oldOK := ClassifyDatastarUIPath(file.OldPath)
			shared = shared || oldShared
			ok = ok && oldOK
			if oldComponent != "" {
				components[oldComponent] = "renamed from " + file.OldPath
			}
		}
		if shared || !ok {
			full = true
			continue
		}
		if component != "" {
			components[component] = "changed " + file.Path
		}
	}
	if full || len(components) == 0 {
		return []E2EJob{fullJob(cfg, "shared or unclassified change")}, nil
	}

	jobs := make([]E2EJob, 0, len(components))
	for component, reason := range components {
		jobs = append(jobs, E2EJob{
			ID:        component,
			Package:   ComponentPackage(component),
			Component: component,
			Reason:    reason,
		})
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].ID < jobs[j].ID })
	return jobs, nil
}

func fullJob(cfg appconfig.Config, reason string) E2EJob {
	pkg := strings.TrimSpace(cfg.RunPackage)
	if pkg == "" {
		pkg = "./tests/e2e"
	}
	return E2EJob{ID: "all", Package: pkg, Reason: reason}
}

func filteredJob(cfg appconfig.Config, opts RunOptions) E2EJob {
	job := fullJob(cfg, "explicit story/scenario filter")
	pattern := SlugToTestFragment(opts.Story)
	if strings.TrimSpace(opts.Scenario) != "" {
		pattern += ".*" + SlugToTestFragment(opts.Scenario)
	}
	job.RunPattern = pattern
	return job
}

func ClassifyDatastarUIPath(path string) (component string, shared bool, ok bool) {
	p := filepath.ToSlash(strings.TrimSpace(path))
	if p == "" {
		return "", false, false
	}
	if strings.HasPrefix(p, "components/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 2 && parts[1] != "" {
			return parts[1], false, true
		}
	}
	if strings.HasPrefix(p, "pages/components/") {
		parts := strings.Split(p, "/")
		if len(parts) >= 3 && parts[2] != "" {
			name := strings.TrimSuffix(parts[2], "page")
			if name != "" && name != parts[2] {
				return name, false, true
			}
		}
	}

	for _, prefix := range []string{"e2e/", "cmd/datastarui/", "scripts/", "layouts/", "static/"} {
		if strings.HasPrefix(p, prefix) {
			return "", true, true
		}
	}
	if sharedFiles[p] {
		return "", true, true
	}
	return "", false, false
}

var sharedFiles = map[string]bool{
	"datastarui-e2e.yml": true,
	"justfile":           true,
	"main.go":            true,
	"go.mod":             true,
	"go.sum":             true,
}

func ComponentPackage(component string) string {
	return "./components/" + component
}

func SlugToTestFragment(slug string) string {
	parts := strings.FieldsFunc(slug, func(r rune) bool {
		return r == '-' || r == '_' || r == ' ' || r == '/' || r == '.'
	})
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}
