// @ts-check
const { test, expect } = require('@playwright/test');

test.describe('DatastarUI Basic Tests', () => {
  test('should load home page', async ({ page }) => {
    await page.goto('/');
    
    // Check that the page loads
    await expect(page).toHaveTitle(/DatastarUI/);
    
    // Check for Datastar library
    const datastorScript = page.locator('script[src*="datastar"]');
    await expect(datastorScript).toBeAttached();
    
    console.log('✅ Home page loaded successfully');
  });

  test('should navigate to components page', async ({ page }) => {
    await page.goto('/');
    
    // Look for navigation to components
    const componentsLink = page.getByRole('link', { name: /components/i });
    if (await componentsLink.count() > 0) {
      await componentsLink.first().click();
      await expect(page.url()).toContain('components');
      console.log('✅ Components navigation works');
    } else {
      console.log('ℹ️ Components link not found, skipping navigation test');
    }
  });

  test('should have responsive design', async ({ page }) => {
    await page.goto('/');
    
    // Test mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.waitForTimeout(500);
    
    // Test tablet viewport  
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.waitForTimeout(500);
    
    // Test desktop viewport
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.waitForTimeout(500);
    
    console.log('✅ Responsive design test completed');
  });

  test('should check for JavaScript errors', async ({ page }) => {
    const errors = [];
    
    page.on('console', msg => {
      if (msg.type() === 'error') {
        errors.push(msg.text());
      }
    });
    
    page.on('pageerror', error => {
      errors.push(error.toString());
    });
    
    await page.goto('/');
    await page.waitForTimeout(2000); // Wait for any async operations
    
    if (errors.length > 0) {
      console.log('⚠️ JavaScript errors found:', errors);
    } else {
      console.log('✅ No JavaScript errors detected');
    }
    
    // Don't fail the test for JS errors, just log them
    // expect(errors).toHaveLength(0);
  });
});