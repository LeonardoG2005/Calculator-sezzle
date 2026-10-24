# Calculator

Full stack calculator application with a React frontend and a Go REST backend.

## Project structure

```text
backend/    Go API and calculator domain
frontend/   React and TypeScript application
```

## Package manager options

The frontend supports either npm or pnpm. Choose one for local development and
use the matching lockfile:

```bash
cd frontend

# npm
npm install
npm run dev

# pnpm
pnpm install
pnpm run dev
```

The two package managers expose the same scripts. With pnpm, the first install
may ask for permission to run esbuild's build script. Approve only `esbuild`
with `pnpm approve-builds`.

The root [Makefile](Makefile) provides shortcuts for building, testing, linting,
coverage, and running the services. Its frontend commands use npm by default.
Use `make PACKAGE_MANAGER=pnpm test-frontend` or
`make PACKAGE_MANAGER=pnpm build-frontend` to use pnpm.
