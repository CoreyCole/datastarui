# Verification

Use this as the DatastarUI verification entrypoint for QRSPI `/q-verify`.

## E2E story testing

Read `docs/e2e-story-testing.md` for the full Go Story E2E contract, config shape, server setup, artifact locations, visual review, and goldens commands.

## Standard commands

```bash
go test ./e2e/...
go test ./components/... -run Test -count=1
templ generate && go build
```

## Browser story verification

DatastarUI component stories need the demo app at `http://localhost:4242`. Prefer the existing Docker/live-reload server:

```bash
just up
just docker-tail app
curl -f http://localhost:4242/components/select
```

Run targeted browser stories with the DatastarUI E2E config:

```bash
just e2e --config datastarui-e2e.yml --base-url http://localhost:4242 --no-restart --story select --viewport desktop-full
```

For DSUI component-page coverage, run the current required story set:

```bash
for story in dropdown select dialog sheet datepicker; do
  just e2e --config datastarui-e2e.yml --base-url http://localhost:4242 --no-restart --story "$story" --viewport desktop-full
done
```

Record `.e2e-runs/<run-id>` artifact paths in the QRSPI `verify.md`.
