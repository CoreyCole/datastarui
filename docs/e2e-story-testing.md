# E2E Story Testing

DatastarUI owns the reusable Go Story E2E library for templ + DatastarUI apps. Authored Go tests are the source of truth; YAML is run/environment wiring only: app name, base URL, run package, artifacts, managed server command, readiness, port env, and viewports. App behavior belongs in Go helpers.

## Config

`datastarui-e2e.yml` is discovered from the current directory upward, or set explicitly with `E2E_CONFIG`.

```yaml
app: datastarui
base_url: http://localhost:4242
run_package: ./components/...
artifacts_dir: .e2e-runs
server:
  command: just build-local
  managed_command: ./datastarui
  skip_when_base_url_set: true
  readiness_path: /components/select
  readiness_timeout: 30s
  port_env: PORT
viewports:
  - desktop-full
```

Consumer contract for another DSUI-style app:

- `server.command` is synchronous setup/build.
- `server.managed_command` is the long-running app command supervised by the runner.
- `server.port_env` receives the runner-assigned port; the app must listen on it.
- `server.readiness_path` is probed before browser stories run.
- Tests receive `E2E_BASE_URL`, `E2E_ARTIFACTS_DIR`, `E2E_VIEWPORTS`, and `E2E_RUN_BROWSER=1`.
- Changed-only default uses DatastarUI's built-in `components/<name>` and `pages/components/<name>page` classifier. Other layouts should use `--all` or explicit `--story` until a custom classifier seam exists.

## Story API

Write component stories beside the component package, for example `components/select/select_component_e2e_test.go`, with an external `_test` package when possible. Use the flat Story builder:

```go
func TestSelectComponent(t *testing.T) {
    spec.Story(t, "select component opens options").
        Visit(spec.Path("/components/select")).
        Do(spec.Click(spec.CSS("[data-slot='select-trigger']"))).
        Expect(spec.ExpectStep(spec.Visible(spec.CSS("[data-slot='select-content']")))).
        Expect(spec.ExpectStep(spec.ConsoleClean())).
        Run()
}
```

Prefer app/component packages that expose typed page, fixture, actor, and expectation objects when selectors repeat. Avoid YAML page/selector registries for app behavior.

Demo pages should document usage and component behavior, not manual QA scripts. When a page needs regression coverage, add a Go Story and keep page copy product/demo oriented.

## Commands

```bash
go test ./e2e/... ./components/... -list Test
just e2e
just e2e --all --jobs 2
just e2e --story select-component --viewport desktop-full
scripts/datastarui.sh e2e review --run .e2e-runs/<run-id> --plan-dir <plan-dir>
scripts/datastarui.sh e2e goldens compare --run .e2e-runs/<run-id> --plan-dir <plan-dir>
scripts/datastarui.sh e2e goldens accept --run .e2e-runs/<run-id> --human-approved --golden-root ./testdata/goldens
```

`just e2e` calls `scripts/datastarui.sh`, which rebuilds the stable launcher `bin/datastarui` only when launcher sources change. The launcher builds `bin/datastarui-runtime-<hash>` when CLI/E2E sources change, then execs that runtime.

## Runner-owned demo server

The normal path does not require manual `just up`:

```bash
just e2e
```

The runner allocates a free localhost port, runs `server.command`, starts `server.managed_command` with `server.port_env`, waits for `server.readiness_path`, runs jobs, writes manifest/summary/index/server logs, and cleans the process up on success or failure.

External server escape hatch:

```bash
just up
curl -f http://localhost:4242/components/select
just e2e --base-url http://localhost:4242 --no-restart --story select --viewport desktop-full
```

Do not leave unmanaged DatastarUI servers running as verification evidence. Managed `e2e run` child processes are allowed because the runner owns cleanup.
