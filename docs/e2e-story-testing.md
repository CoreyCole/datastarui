# E2E Story Testing

DatastarUI owns the reusable Go Story E2E library for templ + DatastarUI apps. Authored Go tests are the source of truth; YAML is run/environment wiring only: app name, base URL, run package, artifacts, server command, and viewports. App behavior belongs in Go helpers.

## Config

`datastarui-e2e.yml` is discovered from the current directory upward, or set explicitly with `E2E_CONFIG`.

```yaml
app: datastarui
base_url: http://localhost:4242
run_package: ./components/...
artifacts_dir: .e2e-runs
server:
  command: just build-local
  skip_when_base_url_set: true
viewports:
  - desktop-full
```

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
just e2e --config datastarui-e2e.yml --base-url http://localhost:4242 --no-restart --story select-component
scripts/datastarui.sh e2e review --run .e2e-runs/<run-id> --plan-dir <plan-dir>
scripts/datastarui.sh e2e goldens compare --run .e2e-runs/<run-id> --plan-dir <plan-dir>
scripts/datastarui.sh e2e goldens accept --run .e2e-runs/<run-id> --human-approved --golden-root ./testdata/goldens
```

`just e2e` calls `scripts/datastarui.sh`, which rebuilds the stable launcher `bin/datastarui` only when launcher sources change. The launcher builds `bin/datastarui-runtime-<hash>` when CLI/E2E sources change, then execs that runtime.

## Running the demo server for browser E2E

DatastarUI component stories need the demo app at `http://localhost:4242`. Do not run `go run main.go`; use the existing Docker/live-reload server or an explicitly managed process.

Preferred local path:

```bash
just up
just docker-tail app
curl -f http://localhost:4242/components/select
just e2e --config datastarui-e2e.yml --base-url http://localhost:4242 --no-restart --story select-component --viewport desktop-full
```

If Docker is unavailable and a human explicitly wants a local one-off process, build first and run the compiled binary in a supervised shell/tmux outside normal verification:

```bash
just build-local
./datastarui
```

Stop the one-off process after testing. Do not leave unmanaged DatastarUI servers running as verification evidence.

## Vamos-managed app processes

Long term, feature-branch browser verification should be owned by Vamos workspace management, not ad hoc ports. Vamos should start app child processes, record their checkout, branch, commit, port, public URL, and latest E2E review index, then expose those links from the main `/workspaces` detail page. Until DatastarUI is registered as a Vamos-managed app, record the demo URL/port and run artifact path manually in verification notes.
