# Verification

Use this as the DatastarUI verification entrypoint for QRSPI `/q-verify`.

## E2E story testing

Read `docs/e2e-story-testing.md` for the full Go Story E2E contract, config shape, managed server setup, artifact locations, visual review, and goldens commands.

## Standard commands

```bash
go test ./e2e/...
go test ./components/... -run Test -count=1
templ generate && go build
```

## Browser story verification

Normal changed-only path:

```bash
just e2e
```

This starts a supervised demo server on a free local port, waits for configured readiness, runs changed component/story jobs vs `main`, writes `.e2e-runs/<run-id>/manifest.json`, `summary.json`, `index.html`, `server.log`, and per-job artifacts, then cleans up the server.

Full suite:

```bash
just e2e --all --jobs 2
```

Targeted story:

```bash
just e2e --story select-component --viewport desktop-full
```

External server escape hatch:

```bash
just up
curl -f http://localhost:4242/components/select
just e2e --base-url http://localhost:4242 --no-restart --story select --viewport desktop-full
```

Record `.e2e-runs/<run-id>` artifact paths in the QRSPI `verify.md`.
