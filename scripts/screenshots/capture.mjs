#!/usr/bin/env node
// Captures the four catalog screens for the wiki guide.
// Opens the Vite app with ?screenshot=1 so Wails bindings are stubbed.
// No Go process and no Wails window.

import { spawn } from 'node:child_process'
import { createServer } from 'node:net'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright'

const here = dirname(fileURLToPath(import.meta.url))
const repo = join(here, '..', '..')
const frontend = join(repo, 'frontend')
const outDir = join(repo, 'docs', 'wiki', 'images')
const width = 1100
const height = 520

const shots = [
  {
    file: 'screen-games.png',
    url: '/?screenshot=1',
    wait: 'Four Souls',
  },
  {
    file: 'screen-collections.png',
    open: 'Four Souls',
    url: '/games/four_souls',
    wait: 'Base Game V2',
  },
  {
    file: 'screen-decks.png',
    open: 'Base Game V2',
    url: '/games/four_souls/collections/base_game',
    wait: 'Character',
  },
  {
    file: 'screen-cards.png',
    open: 'Character',
    url: '/games/four_souls/collections/base_game/decks/character',
    wait: 'Isaac',
  },
]

function unusedPort() {
  return new Promise((resolve, reject) => {
    const server = createServer()
    server.listen(0, '127.0.0.1', () => {
      const { port } = server.address()
      server.close(err => (err ? reject(err) : resolve(port)))
    })
    server.on('error', reject)
  })
}

function startVite(port) {
  const child = spawn('npx', ['vite', '--host', '127.0.0.1', '--port', String(port), '--strictPort'], {
    cwd: frontend,
    stdio: ['ignore', 'pipe', 'pipe'],
    env: { ...process.env, BROWSER: 'none' },
  })
  let output = ''
  child.stdout.on('data', chunk => {
    output += chunk.toString()
  })
  child.stderr.on('data', chunk => {
    output += chunk.toString()
  })

  const ready = (async () => {
    const start = Date.now()
    while (Date.now() - start < 30000) {
      if (child.exitCode !== null) {
        throw new Error(`vite exited ${child.exitCode}:\n${output}`)
      }
      try {
        const res = await fetch(`http://127.0.0.1:${port}/`)
        if (res.ok || res.status === 404) {
          return
        }
      } catch {
        // Not ready yet.
      }
      await new Promise(resolve => setTimeout(resolve, 100))
    }
    throw new Error(`vite did not start within 30s:\n${output}`)
  })()

  return { child, ready }
}

async function imagesReady(page) {
  await page.waitForFunction(() => {
    const images = [...document.querySelectorAll('.card img')]
    return images.length > 0 && images.every(img => img.complete && img.naturalWidth > 0)
  })
  // The list fades in. Capture after that transition, or the cards look washed out.
  await page.waitForFunction(() => {
    const list = document.querySelector('.main__list')
    return list && getComputedStyle(list).opacity === '1'
  })
}

function cardNamed(page, name) {
  return page.locator('.card').filter({
    has: page.locator('.label__name', { hasText: new RegExp(`^${name}$`) }),
  })
}

async function main() {
  const port = await unusedPort()
  const vite = startVite(port)
  const browser = await chromium.launch()
  try {
    await vite.ready
    const page = await browser.newPage({
      viewport: { width, height },
      deviceScaleFactor: 1,
    })
    await page.goto(`http://127.0.0.1:${port}/?screenshot=1`, { waitUntil: 'networkidle' })
    await page.locator('.main').waitFor()

    for (const shot of shots) {
      if (shot.open) {
        await cardNamed(page, shot.open).locator('img').click()
        await page.waitForURL(url => url.pathname === shot.url, { waitUntil: 'commit' })
      }
      await cardNamed(page, shot.wait).waitFor()
      await imagesReady(page)
      const dest = join(outDir, shot.file)
      await page.locator('.main').screenshot({ path: dest, type: 'png' })
      console.log('wrote', dest)
    }
  } finally {
    await Promise.race([
      browser.close().catch(() => {}),
      new Promise(resolve => setTimeout(resolve, 2000)),
    ])
    if (vite.child && vite.child.exitCode === null) {
      vite.child.kill('SIGKILL')
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch(err => {
    console.error(err)
    process.exit(1)
  })
