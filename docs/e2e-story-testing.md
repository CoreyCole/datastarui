# E2E Story Testing

DatastarUI owns the reusable Go Story E2E library for templ + DatastarUI apps. Authored Go tests are the source of truth; YAML only wires app runtime details such as base URL, pages, selectors, artifacts, and server command.

## Config

`vamos-e2e.yaml` is discovered from the current directory upward, or set explicitly with `E2E_CONFIG`.

```yaml
app: datastarui
base_url: http://localhost:4242
run_package: ./tests/e2e
artifacts_dir: .e2e-runs
server:
  command: just build-local
  skip_when_base_url_set: true
auth:
  mode: none
preflight:
  mode: none
pages:
  SelectComponent: /components/select
```

Common environment overrides:

- `E2E_CONFIG` — explicit config file.
- `E2E_BASE_URL` — app URL under test.
- `E2E_ARTIFACTS_DIR` — screenshot/HTML artifact root.
- `E2E_VIEWPORTS` — comma-separated viewport names.
- `E2E_RUN_BROWSER=1` — allow browser tests to run.
- `E2E_CAPTURE_SUCCESS=1` — capture artifacts on passing scenarios too.

## Story API

Write normal Go tests under `tests/e2e`:

```go
func TestSelectComponent(t *testing.T) {
	spec.Feature("Select component").
		Scenario("opens options").
		Given(spec.OpenPage("SelectComponent")).
		When(spec.Click(spec.SelectorAlias("select.trigger"))).
		Then(
			spec.Visible(spec.SelectorAlias("select.content")),
			spec.ConsoleClean(),
		).
		Run(t)
}
```

Prefer semantic locators when they are clear: `Role`, `Text`, `Label`, `TestID`, and `CSS`. Use `SelectorAlias` for repeated or component-specific selectors that are easier to keep in config.

## Commands

List component proof tests without launching a browser:

```bash
go test ./e2e/... ./tests/e2e -list Test
```

Run browser tests when the DatastarUI app is already available at `base_url`:

```bash
E2E_BASE_URL=http://localhost:4242 E2E_RUN_BROWSER=1 go test ./tests/e2e -run 'Test(SelectComponent|DatePickerComponent|TabsComponent)' -count=1
```

Do not run `go run main.go` for local verification. Use the existing Docker/live-reload environment or configured `just build-local` command.
