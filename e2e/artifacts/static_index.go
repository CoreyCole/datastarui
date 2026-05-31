package artifacts

import (
	"html/template"
	"os"
	"path/filepath"
)

type StaticIndexOptions struct{ CDNURL string }

func DefaultDatastarCDNURL() string {
	return "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.1/bundles/datastar.js"
}

func WriteStaticIndex(manifest RunManifest, runDir string, opts StaticIndexOptions) (string, error) {
	if opts.CDNURL == "" {
		opts.CDNURL = DefaultDatastarCDNURL()
	}
	path := filepath.Join(runDir, "index.html")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := staticIndexTemplate.Execute(file, struct {
		Manifest RunManifest
		CDNURL   string
	}{Manifest: manifest, CDNURL: opts.CDNURL}); err != nil {
		return "", err
	}
	return path, nil
}

var staticIndexTemplate = template.Must(template.New("static-index").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>DatastarUI E2E run {{ .Manifest.ID }}</title>
  <script type="module" src="{{ .CDNURL }}"></script>
</head>
<body>
  <main>
    <p>DatastarUI E2E verification artifact</p>
    <h1>Run {{ .Manifest.ID }}</h1>
    <dl>
      <dt>App</dt><dd>{{ .Manifest.App }}</dd>
      <dt>Base URL</dt><dd>{{ .Manifest.BaseURL }}</dd>
      <dt>Viewports</dt><dd>{{ .Manifest.ViewportFilter }}</dd>
    </dl>
    <ul>
    {{ range .Manifest.Artifacts }}
      <li data-kind="{{ .Kind }}"><a href="{{ .Path }}">{{ .FeatureSlug }}/{{ .ScenarioSlug }}/{{ .Viewport }} {{ .Label }}</a></li>
    {{ end }}
    </ul>
  </main>
</body>
</html>
`))
