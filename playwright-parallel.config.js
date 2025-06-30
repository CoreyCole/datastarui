// @ts-check
const { defineConfig, devices } = require('@playwright/test');

/**
 * Optimized Playwright config for parallel execution
 * Use this for faster test runs when you want maximum parallelization
 */
module.exports = defineConfig({
  testDir: './tests',
  
  /* Maximum parallelization */
  fullyParallel: true,
  
  /* Use more workers for faster execution */
  workers: process.env.CI ? 4 : '75%',
  
  /* Shorter timeouts for faster feedback */
  timeout: 30000,
  expect: { timeout: 10000 },
  
  /* Fail fast in CI */
  forbidOnly: !!process.env.CI,
  
  /* Minimal retries for speed */
  retries: process.env.CI ? 1 : 0,
  
  /* Optimized reporter for parallel runs */
  reporter: process.env.CI ? [
    ['json', { outputFile: 'test-results/results.json' }],
    ['github']
  ] : [
    ['list'],
    ['html', { outputFolder: 'test-results/html-report', open: 'never' }]
  ],
  
  use: {
    baseURL: 'http://localhost:4242',
    
    /* Faster navigation */
    navigationTimeout: 15000,
    actionTimeout: 10000,
    
    /* Minimal tracing for speed */
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    
    /* Optimize browser context */
    launchOptions: {
      args: [
        '--disable-web-security',
        '--disable-features=TranslateUI',
        '--no-first-run',
        '--disable-default-apps'
      ]
    }
  },

  /* Run only on Chromium for maximum speed */
  projects: [
    {
      name: 'chromium-fast',
      use: { 
        ...devices['Desktop Chrome'],
        /* Optimize for speed */
        viewport: { width: 1280, height: 720 },
        ignoreHTTPSErrors: true,
      },
    },
    
    /* Add Firefox and Safari only if needed */
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
  ],

  webServer: {
    command: 'just watch',
    url: 'http://localhost:4242',
    reuseExistingServer: !process.env.CI,
    timeout: 30000,
    /* Don't wait for server output in parallel mode */
    stdout: 'pipe',
    stderr: 'pipe',
  },
});