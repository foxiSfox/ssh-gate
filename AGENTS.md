# Repository Guidelines

## Project Structure & Module Organization
The Go backend lives in `backend/`. `main.go` wires the HTTP router, serves `frontend/dist`, and coordinates access to SQLite via helpers in `db/`. API logic is grouped under `handlers/`, reusable HTTP logic in `middleware/`, data contracts in `models/`, SSH key management in `ssh/`, and validation utilities in `security/`. Runtime files such as `users.db` and the aggregated `authorized_keys` stay beside the source so the container image can mount them. The Vue 3 client is under `frontend/`, with feature slices in `src/modules/`, shared UI and API hooks in `src/shared/`, and entrypoints in `main.ts` and `router.ts`. Built assets land in `frontend/dist`; rebuild rather than editing this directory.

## Build, Test, and Development Commands
Install and develop the UI with `cd frontend && npm ci && npm run dev` (Vite serves on port 5173). Use `npm run build` to regenerate `dist`, `npm run lint` for ESLint + Prettier, `npm run format` to rewrite `src/`, and `npm run type-check` for Vue TSC. For the API use `cd backend && go mod download`, `go run main.go` for local execution, and `go build ./...` to produce a binary. Run `go test ./...` before every PR; add focused `_test.go` suites as you extend handlers or middleware. The provided `Dockerfile` builds a combined image: `docker build -t ssh-gate .`.

## Coding Style & Naming Conventions
Go code must stay `gofmt`-clean (tabs, PascalCase exports, short-lived lowerCamel locals). Prefer dependency injection through constructors in `models` or `handlers` and keep package-level state minimal. Vue components use `<script setup>` with TypeScript and PascalCase filenames inside `src/modules`. Shared composables and helpers belong in `src/shared` and should export camelCase functions. Run `npm run lint` or enable the ESLint config in your editor; Prettier keeps markup to two-space indentation.

## Testing Guidelines
Backend tests belong next to the code they cover (`handlers/user_handlers_test.go`, etc.) and should exercise both happy-path and SSH failure modes using fakes over `ssh/` interfaces. Use table-driven patterns to keep coverage consistent. Frontend testing is currently manual; when adding automated coverage, follow the `*.spec.ts` naming under `src/` and document the command you introduce. Always validate UI changes with `npm run type-check` and a local Vite session.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `docs:`) as seen in the history, and squash fixups before pushing. Reference Jira or GitHub issue IDs when applicable. Each PR should include a concise summary, testing notes (`go test`, `npm run build`, screenshots for UI), and call out any database or SSH side effects. Keep PRs focused on a single feature or fix to ease review.

## Security & Configuration Tips
Never commit real SSH keys or production `users.db` files; use redacted fixtures instead. Validate that new SSH flows update both remote hosts and `backend/authorized_keys`. Store secrets in environment variables or Docker secrets, not in source control, and document any new variables in `README.md`.
