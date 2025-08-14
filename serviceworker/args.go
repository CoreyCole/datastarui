package serviceworker

// ManifestArgs defines the configuration for the PWA manifest
type ManifestArgs struct {
	Name            string
	ShortName       string
	Description     string
	StartURL        string
	Scope           string
	ThemeColor      string
	BackgroundColor string
	Display         string // standalone, minimal-ui, browser
	Orientation     string // portrait, landscape, any
	Icons           []Icon
	Categories      []string
}

// Icon defines app icons for different sizes
type Icon struct {
	Src     string
	Sizes   string
	Type    string
	Purpose string // any, maskable, monochrome
}

// RegisterArgs defines service worker registration options
type RegisterArgs struct {
	Enabled        bool
	Scope          string
	UpdateViaCache string // imports, all, none
	SkipWaiting    bool
	ClientsClaim   bool
}

// CacheArgs defines caching strategy configuration
type CacheArgs struct {
	Name       string
	Version    string
	MaxAge     int // seconds
	MaxEntries int
	Strategy   string // network-first, cache-first, stale-while-revalidate
}

// ServiceWorkerArgs defines the complete service worker configuration
type ServiceWorkerArgs struct {
	Manifest        ManifestArgs
	Register        RegisterArgs
	Cache           CacheArgs
	Assets          []string // files to precache
	Routes          []string // routes to cache
	OfflineFallback string   // fallback page for offline
}
