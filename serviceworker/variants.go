package serviceworker

// DefaultIcons returns complete favicon and PWA icons configuration
func DefaultIcons() []Icon {
	return []Icon{
		// Android PWA icons
		{
			Src:     "/favicons/android-chrome-192x192.png",
			Sizes:   "192x192",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/android-chrome-512x512.png",
			Sizes:   "512x512",
			Type:    "image/png",
			Purpose: "any",
		},
		// Android maskable icons
		{
			Src:     "/favicons/maskable-192x192.png",
			Sizes:   "192x192",
			Type:    "image/png",
			Purpose: "maskable",
		},
		{
			Src:     "/favicons/maskable-512x512.png",
			Sizes:   "512x512",
			Type:    "image/png",
			Purpose: "maskable",
		},
		// iOS icons
		{
			Src:     "/favicons/apple-touch-icon.png",
			Sizes:   "180x180",
			Type:    "image/png",
			Purpose: "any",
		},
		// Standard favicons
		{
			Src:     "/favicons/favicon.ico",
			Sizes:   "48x48",
			Type:    "image/x-icon",
			Purpose: "any",
		},
		{
			Src:     "/favicons/favicon-32x32.png",
			Sizes:   "32x32",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/favicon-16x16.png",
			Sizes:   "16x16",
			Type:    "image/png",
			Purpose: "any",
		},
		// Safari pinned tab
		{
			Src:     "/favicons/safari-pinned-tab.svg",
			Sizes:   "any",
			Type:    "image/svg+xml",
			Purpose: "maskable",
		},
		// Microsoft tiles
		{
			Src:     "/favicons/mstile-70x70.png",
			Sizes:   "70x70",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/mstile-144x144.png",
			Sizes:   "144x144",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/mstile-150x150.png",
			Sizes:   "150x150",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/mstile-310x150.png",
			Sizes:   "310x150",
			Type:    "image/png",
			Purpose: "any",
		},
		{
			Src:     "/favicons/mstile-310x310.png",
			Sizes:   "310x310",
			Type:    "image/png",
			Purpose: "any",
		},
	}
}

// DefaultManifest returns standard PWA manifest configuration
func DefaultManifest() ManifestArgs {
	return ManifestArgs{
		Name:            "DatastarUI",
		ShortName:       "DatastarUI",
		Description:     "Go/templ port of shadcn/ui with Datastar reactivity",
		StartURL:        "/",
		Scope:           "/",
		ThemeColor:      "#000000",
		BackgroundColor: "#ffffff",
		Display:         "standalone",
		Orientation:     "any",
		Icons:           DefaultIcons(),
		Categories:      []string{"development", "utilities"},
	}
}

// DefaultRegister returns standard service worker registration
func DefaultRegister() RegisterArgs {
	return RegisterArgs{
		Enabled:        true,
		Scope:          "/",
		UpdateViaCache: "imports",
		SkipWaiting:    true,
		ClientsClaim:   true,
	}
}

// DefaultCache returns standard cache configuration
func DefaultCache() CacheArgs {
	return CacheArgs{
		Name:       "datastarui",
		Version:    "1.0.0",
		MaxAge:     86400, // 24 hours
		MaxEntries: 50,
		Strategy:   "network-first",
	}
}

// DefaultAssets returns complete favicon assets to precache
func DefaultAssets() []string {
	return []string{
		"/css/build.css",
		"/favicons/favicon.ico",
		"/favicons/apple-touch-icon.png",
		"/favicons/android-chrome-192x192.png",
		"/favicons/android-chrome-512x512.png",
		"/favicons/maskable-192x192.png",
		"/favicons/maskable-512x512.png",
		"/favicons/favicon-32x32.png",
		"/favicons/favicon-16x16.png",
		"/favicons/safari-pinned-tab.svg",
		"/favicons/mstile-70x70.png",
		"/favicons/mstile-144x144.png",
		"/favicons/mstile-150x150.png",
		"/favicons/mstile-310x150.png",
		"/favicons/mstile-310x310.png",
		"/static/img/logo.svg",
	}
}

// DefaultRoutes returns standard routes to cache
func DefaultRoutes() []string {
	return []string{
		"/",
		"/components/",
		"/components/button",
		"/components/card",
		"/components/sheet",
		"/components/dialog",
		"/components/dropdown",
		"/components/select",
		"/components/popover",
		"/components/calendar",
		"/components/datepicker",
		"/components/checkbox",
		"/components/form",
		"/components/input",
		"/components/label",
		"/components/breadcrumb",
		"/components/tabs",
		"/docs",
	}
}

// DefaultServiceWorker returns complete service worker configuration
func DefaultServiceWorker() ServiceWorkerArgs {
	return ServiceWorkerArgs{
		Manifest:        DefaultManifest(),
		Register:        DefaultRegister(),
		Cache:           DefaultCache(),
		Assets:          DefaultAssets(),
		Routes:          DefaultRoutes(),
		OfflineFallback: "/offline",
	}
}
