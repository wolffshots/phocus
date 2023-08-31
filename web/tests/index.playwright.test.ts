import { expect, test } from '@playwright/test';

test('expect save button to be visible', async ({ page }) => {
        await page.goto('/');
        await expect(page.getByRole('button', { name: 'Save' })).toBeVisible();
});

test('expect loading... to be visible', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('cell', { name: 'Loading...' })).toBeVisible();
});