const { chromium } = require('playwright');
const fs = require('fs');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  const consoleMessages = [];
  page.on('console', msg => {
    consoleMessages.push({
      type: msg.type(),
      text: msg.text()
    });
    console.log(`CONSOLE [${msg.type()}]: ${msg.text()}`);
  });

  const errors = [];
  page.on('pageerror', err => {
    errors.push(err.message);
    console.log(`PAGE ERROR: ${err.message}`);
  });

  const requests = [];
  page.on('requestfailed', request => {
    requests.push({
      url: request.url(),
      error: request.failure().errorText
    });
    console.log(`REQUEST FAILED: ${request.url()} - ${request.failure().errorText}`);
  });

  try {
    console.log('Navigating to https://rustchain.org/beacon/...');
    await page.goto('https://rustchain.org/beacon/', { waitUntil: 'networkidle', timeout: 60000 });
    
    // Wait for "Loading..." to potentially disappear or stay
    console.log('Waiting 30 seconds to observe state and capture errors...');
    await page.waitForTimeout(30000);
    
    const content = await page.content();
    fs.writeFileSync('/home/administrator/.gemini/tmp/bountyos/atlas_content.html', content);
    
    await page.screenshot({ path: '/home/administrator/.gemini/tmp/bountyos/atlas_screenshot.png', fullPage: true });
    
    const results = {
      consoleMessages,
      errors,
      requests,
      timestamp: new Date().toISOString()
    };
    
    fs.writeFileSync('/home/administrator/.gemini/tmp/bountyos/atlas_debug.json', JSON.stringify(results, null, 2));
    console.log('Debug info saved to /home/administrator/.gemini/tmp/bountyos/atlas_debug.json');
    
  } catch (error) {
    console.error('Error during capture:', error);
  } finally {
    await browser.close();
  }
})();
