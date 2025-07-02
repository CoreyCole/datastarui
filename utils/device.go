package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// IsMobile detects if the request is from a mobile device based on User-Agent
func IsMobile(c echo.Context) bool {
	userAgent := strings.ToLower(c.Request().Header.Get("User-Agent"))
	
	// Common mobile device indicators
	mobileKeywords := []string{
		"mobile",
		"android",
		"iphone",
		"ipod",
		"ipad",
		"windows phone",
		"blackberry",
		"opera mini",
		"opera mobi",
	}
	
	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}
	
	return false
}

// DeviceType represents the type of device making the request
type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
)

// GetDeviceType returns a more granular device type
func GetDeviceType(c echo.Context) DeviceType {
	userAgent := strings.ToLower(c.Request().Header.Get("User-Agent"))
	
	// Check for tablet first (iPads can contain "mobile" in UA)
	if strings.Contains(userAgent, "ipad") || strings.Contains(userAgent, "tablet") {
		return DeviceTablet
	}
	
	// Then check for mobile
	if IsMobile(c) {
		return DeviceMobile
	}
	
	return DeviceDesktop
}