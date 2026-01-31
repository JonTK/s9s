# s9s-web Integration Guide

This document provides step-by-step instructions for integrating the s9s documentation into the s9s-web repository via git submodule.

## Overview

The s9s documentation now lives in this repository (`s9s/docs/`) and will be consumed by s9s-web via a git submodule. This ensures documentation updates with code changes in a single PR.

## Steps to Integrate

### 1. Add s9s as Submodule in s9s-web

In the s9s-web repository:

```bash
cd /path/to/s9s-web

# Add s9s as a submodule
git submodule add https://github.com/jontk/s9s.git s9s-docs

# The docs will be accessible at: s9s-web/s9s-docs/docs/
```

### 2. Update s9s-web Documentation Path

Edit `s9s-web/lib/docs.ts`:

```typescript
import fs from 'fs'
import path from 'path'
import matter from 'gray-matter'

// UPDATE THIS PATH to read from submodule
const DOCS_PATH = path.join(process.cwd(), 's9s-docs', 'docs')

// Rest of the file stays the same
export async function getDocBySlug(slug: string[]) {
  const realSlug = slug.join('/')
  const fullPath = path.join(DOCS_PATH, `${realSlug}.md`)

  // ... existing code
}

// ... rest of file
```

### 3. Update Navigation Structure

Edit `s9s-web/lib/constants.ts` to match new documentation structure:

```typescript
export const DOCS_NAVIGATION = [
  {
    title: 'Getting Started',
    items: [
      { title: 'Overview', slug: 'index' },
      { title: 'Installation', slug: 'getting-started/installation' },
      { title: 'Quick Start', slug: 'getting-started/quickstart' },
      { title: 'Configuration', slug: 'getting-started/configuration' },
    ],
  },
  {
    title: 'User Guide',
    items: [
      { title: 'Navigation', slug: 'user-guide/navigation' },
      { title: 'Keyboard Shortcuts', slug: 'user-guide/keyboard-shortcuts' },
      {
        title: 'Views',
        slug: 'user-guide/views/index',
        children: [
          { title: 'Dashboard', slug: 'user-guide/views/dashboard' },
          { title: 'Jobs', slug: 'user-guide/views/jobs' },
          { title: 'Nodes', slug: 'user-guide/views/nodes' },
          { title: 'Partitions', slug: 'user-guide/views/partitions' },
          { title: 'Users', slug: 'user-guide/views/users' },
          { title: 'Accounts', slug: 'user-guide/views/accounts' },
          { title: 'QoS', slug: 'user-guide/views/qos' },
          { title: 'Reservations', slug: 'user-guide/views/reservations' },
          { title: 'Health', slug: 'user-guide/views/health' },
        ],
      },
      { title: 'Job Management', slug: 'user-guide/job-management' },
      { title: 'Node Operations', slug: 'user-guide/node-operations' },
      { title: 'Batch Operations', slug: 'user-guide/batch-operations' },
      { title: 'Filtering & Search', slug: 'user-guide/filtering' },
      { title: 'Export & Reporting', slug: 'user-guide/export' },
    ],
  },
  {
    title: 'Guides',
    items: [
      { title: 'SSH Integration', slug: 'guides/ssh-integration' },
      { title: 'Job Streaming', slug: 'guides/job-streaming' },
      { title: 'Mock Mode', slug: 'guides/mock-mode' },
      { title: 'Troubleshooting', slug: 'guides/troubleshooting' },
    ],
  },
  {
    title: 'Reference',
    items: [
      { title: 'Commands', slug: 'reference/commands' },
      { title: 'Configuration', slug: 'reference/configuration' },
      { title: 'API', slug: 'reference/api' },
    ],
  },
  {
    title: 'Plugins',
    items: [
      { title: 'Overview', slug: 'plugins/overview' },
      { title: 'Development', slug: 'plugins/development' },
      { title: 'Observability', slug: 'plugins/observability' },
    ],
  },
  {
    title: 'Development',
    items: [
      { title: 'Setup', slug: 'development/setup' },
      { title: 'Architecture', slug: 'development/architecture' },
      { title: 'Contributing', slug: 'development/contributing' },
      { title: 'Testing', slug: 'development/testing' },
      { title: 'Linting', slug: 'development/linting' },
      { title: 'CI/CD', slug: 'development/ci-cd' },
    ],
  },
  {
    title: 'Enterprise',
    items: [
      { title: 'Features', slug: 'enterprise/features' },
    ],
  },
  {
    title: 'About',
    items: [
      { title: 'Changelog', slug: 'about/changelog' },
      { title: 'License', slug: 'about/license' },
    ],
  },
] as const
```

### 4. Update Next.js Dynamic Routes

Edit `s9s-web/app/docs/[...slug]/page.tsx` to handle nested paths:

```typescript
import { getDocBySlug, getAllDocs } from '@/lib/docs'
import { notFound } from 'next/navigation'
import { MDXRemote } from 'next-mdx-remote/rsc'

export async function generateStaticParams() {
  const docs = await getAllDocs()
  return docs.map((doc) => ({
    slug: doc.slug.split('/'),
  }))
}

export default async function DocPage({ params }: { params: { slug: string[] } }) {
  const doc = await getDocBySlug(params.slug || ['index'])

  if (!doc) {
    notFound()
  }

  return (
    <div className="prose prose-slate dark:prose-invert max-w-none">
      <h1>{doc.title}</h1>
      <MDXRemote source={doc.content} />
    </div>
  )
}
```

### 5. Handle Demo GIF Assets

#### Option A: Copy During Build (Recommended for Production)

Edit `s9s-web/package.json`:

```json
{
  "scripts": {
    "prebuild": "git submodule update --init --recursive && mkdir -p public/assets/demos && cp -r s9s-docs/demos/output/*.gif public/assets/demos/",
    "build": "next build",
    "dev": "next dev"
  }
}
```

#### Option B: Symlink (For Development)

```bash
cd s9s-web/public
mkdir -p assets
ln -s ../../s9s-docs/demos/output assets/demos
```

### 6. Update Build Pipeline

Edit `.github/workflows/deploy.yml` (or equivalent CI/CD config):

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout with submodules
        uses: actions/checkout@v4
        with:
          submodules: recursive  # This is important!

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build

      - name: Deploy
        # ... your deployment steps
```

### 7. Update Image Configuration

Edit `s9s-web/next.config.js` to allow demo GIF serving:

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: [],
    unoptimized: true, // GIFs don't need optimization
  },
  // ... other config
}

module.exports = nextConfig
```

### 8. Delete Old Documentation

Remove the old `s9s-web/docs/` folder:

```bash
cd s9s-web
rm -rf docs/
git rm -rf docs/
git commit -m "Remove old docs - now using s9s submodule"
```

### 9. Update .gitmodules

The `.gitmodules` file should be automatically created:

```ini
[submodule "s9s-docs"]
    path = s9s-docs
    url = https://github.com/jontk/s9s.git
```

### 10. Initialize Submodule for Other Developers

Other developers will need to initialize the submodule after cloning:

```bash
git clone https://github.com/jontk/s9s-web.git
cd s9s-web
git submodule update --init --recursive
```

Add this to `s9s-web/README.md`:

````markdown
## Development Setup

1. Clone the repository with submodules:
```bash
git clone --recurse-submodules https://github.com/jontk/s9s-web.git
cd s9s-web
```

Or if you already cloned it:
```bash
git submodule update --init --recursive
```

2. Install dependencies:
```bash
npm install
```

3. Run development server:
```bash
npm run dev
```
````

## Testing the Integration

### 1. Test Locally

```bash
cd s9s-web
npm run dev
```

Visit http://localhost:3000/docs and verify:
- [ ] All documentation pages load
- [ ] Demo GIFs display correctly
- [ ] Navigation works
- [ ] Links between docs work
- [ ] Search finds content
- [ ] Mobile responsive layout works

### 2. Test Build

```bash
npm run build
npm start
```

Check production build for:
- [ ] No build errors
- [ ] All pages generate correctly
- [ ] Assets are copied/available
- [ ] Performance is acceptable

## Updating Documentation

When s9s documentation is updated:

```bash
cd s9s-web
git submodule update --remote s9s-docs
git add s9s-docs
git commit -m "Update documentation to latest"
git push
```

## Troubleshooting

### Submodule Not Initialized

```bash
git submodule update --init --recursive
```

### Demo GIFs Not Loading

Check that:
1. `public/assets/demos/` exists
2. GIF files are present
3. Paths in markdown use `/assets/demos/` (absolute from public)

### Documentation Pages 404

Verify:
1. `DOCS_PATH` in `lib/docs.ts` points to `s9s-docs/docs`
2. Slug paths in `lib/constants.ts` match file structure
3. Files exist at expected paths in submodule

### Stale Documentation

```bash
cd s9s-web
git submodule update --remote --merge
```

## Benefits of This Approach

1. **Single Source of Truth**: Documentation lives with the code
2. **Atomic Updates**: Docs and code update together in PRs
3. **Version Syncing**: s9s-web automatically gets latest docs
4. **No Duplication**: One documentation source, multiple consumers
5. **Better DX**: Developers work in one repo for code + docs

## Migration Checklist

- [ ] Add s9s submodule to s9s-web
- [ ] Update `lib/docs.ts` DOCS_PATH
- [ ] Update `lib/constants.ts` navigation
- [ ] Update `app/docs/[...slug]/page.tsx` for nested routes
- [ ] Update `package.json` prebuild script
- [ ] Update CI/CD to checkout with submodules
- [ ] Update `next.config.js` for images
- [ ] Delete old `s9s-web/docs/` folder
- [ ] Update README with submodule instructions
- [ ] Test locally
- [ ] Test build
- [ ] Deploy to staging
- [ ] Verify on staging
- [ ] Deploy to production
- [ ] Update documentation links in other repos

## Support

If you encounter issues during integration, please:
1. Check this guide's troubleshooting section
2. Verify submodule is initialized and up to date
3. Check build logs for errors
4. Open an issue on GitHub with reproduction steps
