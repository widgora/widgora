# Contributing to Widgora

Thanks for your interest. Widgora is in early development, so larger changes are best discussed in an issue first.

## Build and test

You need Go (version in `go.mod`), Node.js (`.node-version`), pnpm (`packageManager` in `web/package.json`), GNU
make and [golangci-lint](https://golangci-lint.run/) v2.

| Command | What it does |
|---|---|
| `make dev` | Go server on 127.0.0.1:8080 plus the Vite dev server with hot reload |
| `make build` | Builds the web UI and embeds it into one binary, `dist/widgora` |
| `make test` | Go tests with the race detector, web unit tests |
| `make lint` | golangci-lint (including gofumpt), Prettier, ESLint, svelte-check |
| `make fmt` | Formats Go and web code |
| `make vuln` | govulncheck |
| `make licenses` | Checks that every shipped dependency has an allowed license |

CI runs the same targets; a pull request needs all of them green.

## Pull requests

- One topic per pull request, against `main`. Pull requests are squash-merged.
- Title: `area: what changed`, for example `ui: show the host name in the page title`.
- Describe what changed and why, and how you tested it (the template has the sections).
- User-visible changes get an entry under `Unreleased` in `CHANGELOG.md`.
- Updates must not need a reinstall or manual migration.

## Dependencies

Every dependency must be open source and compatible with AGPL-3.0-only; `.licenses-allowed` lists the accepted
licenses and `make licenses` enforces it. Web packages that end up in the UI bundle belong in `dependencies` (they
are license-checked), build tools in `devDependencies`.

## License

Widgora is licensed under [AGPL-3.0-only](LICENSE). By contributing you agree that your contribution is licensed
under the same terms.

Security problems: see [SECURITY.md](SECURITY.md), never a public issue.
