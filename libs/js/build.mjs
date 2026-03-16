import * as esbuild from 'esbuild'
import { execSync } from 'child_process'

// Generate type declarations
execSync('npx tsc --declaration --declarationMap --emitDeclarationOnly --outDir dist', { stdio: 'inherit' })

// CJS bundle
await esbuild.build({
  entryPoints: ['src/index.ts'],
  bundle: true,
  platform: 'node',
  format: 'cjs',
  outfile: 'dist/index.cjs',
})

// ESM bundle
await esbuild.build({
  entryPoints: ['src/index.ts'],
  bundle: true,
  format: 'esm',
  outfile: 'dist/index.mjs',
})

// IIFE bundle for browser (exposes window.MetricCalendar)
await esbuild.build({
  entryPoints: ['src/index.ts'],
  bundle: true,
  format: 'iife',
  globalName: 'MetricCalendar',
  outfile: 'dist/metric-calendar.iife.js',
})

console.log('Build complete')
