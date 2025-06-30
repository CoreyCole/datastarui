# Playwright Testing Setup

This document explains how to use Playwright for browser automation testing with the DatastarUI project.

## Quick Start

### 1. Start the Development Environment
```bash
# Start the main app and Tailwind watcher
just up

# In a separate terminal, start Playwright container
just playwright-up
```

### 2. Run Tests
```bash
# Run all tests
just playwright-test

# Run tests with browser UI (for debugging)
just playwright-test-headed

# Open Playwright UI for interactive testing
just playwright-ui
```

### 3. Interactive Browser Session
```bash
# Open shell in Playwright container
just playwright-shell

# Inside the container, you can run:
npx playwright test                    # Run tests
npx playwright test --debug           # Debug mode
npx playwright codegen http://app:4242 # Record new tests
```

## Container Architecture

The Playwright setup runs in a separate Docker container that:
- Has access to the DatastarUI app via Docker networking (`http://app:4242`)
- Shares the project files for test development
- Runs browsers in headless mode by default
- Persists browser downloads and cache

## Configuration

### Browser Configuration
- **Chromium**: Desktop Chrome simulation
- **Firefox**: Desktop Firefox simulation  
- **WebKit**: Desktop Safari simulation

### Test Settings
- **Base URL**: `http://app:4242` (in container)
- **Parallel**: Tests run in parallel for speed
- **Screenshots**: Captured on test failure
- **Videos**: Recorded on test failure
- **Traces**: Collected on retry

## Writing Tests

### Basic Test Structure
```javascript
const { test, expect } = require('@playwright/test');

test('component behavior', async ({ page }) => {
  await page.goto('/components/button');
  
  // Test component interactions
  await page.click('[data-testid="primary-button"]');
  
  // Verify Datastar signal updates
  const signal = await page.locator('[data-signals]').getAttribute('data-signals');
  expect(signal).toContain('clicked');
});
```

### Testing Datastar Components

When testing Datastar components, focus on:

1. **Signal State**: Check `data-signals` attributes
2. **DOM Updates**: Verify reactive updates work
3. **Server Communication**: Test SSE and fetch requests
4. **Accessibility**: Ensure ARIA attributes are correct

### Example Datastar Test
```javascript
test('button component with signals', async ({ page }) => {
  await page.goto('/components/button');
  
  // Check initial signal state
  const buttonSignals = await page.locator('[data-signals]').first();
  const initialState = await buttonSignals.getAttribute('data-signals');
  
  // Interact with component
  await page.click('button[data-on-click]');
  
  // Verify signal update
  await page.waitForTimeout(100);
  const updatedState = await buttonSignals.getAttribute('data-signals');
  expect(updatedState).not.toBe(initialState);
});
```

## Debugging

### Common Issues

1. **Container Connectivity**: Ensure the app container is running first
2. **Port Conflicts**: Use `http://app:4242` inside container, not `localhost:4242`
3. **Signal Timing**: Add `waitForTimeout()` after Datastar interactions

### Debug Commands
```bash
# Check container logs
docker logs datastarui-local-playwright-1

# Check app accessibility from Playwright container
just playwright-shell
curl http://app:4242

# View browser console errors
just playwright-test --headed
```

### Playwright UI Mode

For visual debugging:
```bash
just playwright-ui
```

Access the UI at `http://localhost:8080` to:
- Record new tests interactively
- Debug existing tests step-by-step
- Inspect DOM and network requests
- View screenshots and videos

## CI/CD Integration

The setup is designed to work in CI environments:
- Uses `CI=true` environment variable
- Runs in headless mode
- Outputs test results in JSON format
- Captures artifacts on failure

## File Structure

```
datastarui/
├── Dockerfile.playwright          # Playwright container setup
├── playwright.config.js           # Test configuration
├── playwright-parallel.config.js  # Parallel test config
├── tests/                         # Test files
│   └── basic.spec.js              # Basic connectivity tests
└── test-results/                  # Test output
    ├── html-report/               # HTML test reports
    └── results.json               # JSON test results
```

## Next Steps

1. **Component Coverage**: Add tests for each component
2. **Visual Regression**: Implement screenshot comparison
3. **Accessibility**: Add automated a11y testing
4. **Performance**: Add performance testing with Lighthouse
5. **Cross-browser**: Ensure tests pass on all browsers