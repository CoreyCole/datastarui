package serviceworker

import (
	"strings"
	"text/template"
)

// GenerateServiceWorker generates the service worker JavaScript
func GenerateServiceWorker(args ServiceWorkerArgs) (string, error) {
	swTemplate := `const CACHE_NAME = '{{.Cache.Name}}-v{{.Cache.Version}}';
const PRECACHE_URLS = [{{range $i, $asset := .Assets}}{{if $i}}, {{end}}"{{$asset}}"{{end}}];
const ROUTE_URLS = [{{range $i, $route := .Routes}}{{if $i}}, {{end}}"{{$route}}"{{end}}];
const OFFLINE_FALLBACK = '{{.OfflineFallback}}';
const STRATEGY = '{{.Cache.Strategy}}';

self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME)
			.then((cache) => cache.addAll(PRECACHE_URLS))
			.then(() => self.skipWaiting())
	);
});

self.addEventListener('activate', (event) => {
	event.waitUntil(
		caches.keys().then((cacheNames) =>
			Promise.all(
				cacheNames
					.filter(name => name !== CACHE_NAME)
					.map(name => caches.delete(name))
			)
		)
	);
	return self.clients.claim();
});

self.addEventListener('fetch', (event) => {
	const url = new URL(event.request.url);
	
	if (event.request.method !== 'GET' || url.protocol === 'chrome-extension:') {
		return;
	}

	if (PRECACHE_URLS.some(asset => url.pathname.includes(asset))) {
		event.respondWith(
			caches.match(event.request)
				.then(response => response || fetch(event.request))
		);
		return;
	}

	if (ROUTE_URLS.some(route => url.pathname.startsWith(route))) {
		if (STRATEGY === 'network-first') {
			event.respondWith(
				fetch(event.request)
					.then(response => {
						const responseClone = response.clone();
						caches.open(CACHE_NAME).then(cache => {
							cache.put(event.request, responseClone);
						});
						return response;
					})
					.catch(() => caches.match(event.request)
						.then(response => response || caches.match(OFFLINE_FALLBACK)))
			);
		} else if (STRATEGY === 'cache-first') {
			event.respondWith(
				caches.match(event.request)
					.then(response => response || fetch(event.request))
			);
		}
		return;
	}
});`

	tmpl, err := template.New("sw").Parse(swTemplate)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, args)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

// GenerateManifest generates the PWA manifest JSON
func GenerateManifest(args ManifestArgs) (string, error) {
	manifestTemplate := `{
		"name": "{{.Name}}",
		"short_name": "{{.ShortName}}",
		"description": "{{.Description}}",
		"start_url": "{{.StartURL}}",
		"scope": "{{.Scope}}",
		"theme_color": "{{.ThemeColor}}",
		"background_color": "{{.BackgroundColor}}",
		"display": "{{.Display}}",
		"orientation": "{{.Orientation}}",
		"icons": [{{range $i, $icon := .Icons}}{{if $i}}, {{end}}{
			"src": "{{.Src}}",
			"sizes": "{{.Sizes}}",
			"type": "{{.Type}}",
			"purpose": "{{.Purpose}}"
		}{{end}}],
		"categories": [{{range $i, $cat := .Categories}}{{if $i}}, {{end}}"{{$cat}}"{{end}}]
	}`

	tmpl, err := template.New("manifest").Parse(manifestTemplate)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, args)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
