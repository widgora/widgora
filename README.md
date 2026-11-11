# Widgora

Widgora is an inventory and documentation hub for homelabs. It discovers hosts, VMs, containers, services,
addresses and networks, keeps the full history of every change, and keeps your own notes and journal right next
to the discovered facts.

**Status: early development.** Nothing is usable yet; expect breaking changes until the first release.

## Build from source

Requirements: Go (version in `go.mod`), Node.js (version in `.node-version`), pnpm (version in
`web/package.json`), GNU make, and [golangci-lint](https://golangci-lint.run/) for `make lint`.

```sh
make build              # web UI, then one binary in dist/
./dist/widgora serve    # http://127.0.0.1:8080
```

`widgora serve --help` lists the flags (`--listen`, `--data-dir`, `--log-level`, `--log-format`).
See [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow.

## License

[AGPL-3.0-only](LICENSE).
