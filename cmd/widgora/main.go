// Command widgora is the Widgora server.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/widgora/widgora/internal/server"
	"github.com/widgora/widgora/internal/version"
	"github.com/widgora/widgora/internal/webui"
)

const usage = `Usage: widgora <command> [flags]

Commands:
  serve     run the HTTP server
  version   print version information

Run "widgora <command> -h" for the flags of a command.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	var err error
	switch args[0] {
	case "serve":
		err = serve(args[1:], stderr)
	case "version", "--version":
		fmt.Fprintf(stdout, "widgora %s (commit %s, built %s)\n", version.Version, version.Commit, version.Date)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
		return 2
	}
	switch {
	case errors.Is(err, flag.ErrHelp):
		return 0
	case err != nil:
		fmt.Fprintln(stderr, "widgora:", err)
		return 1
	}
	return 0
}

func serve(args []string, stderr io.Writer) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	listen := fs.String("listen", "127.0.0.1:8080", "address to listen on (host:port)")
	dataDir := fs.String("data-dir", "data", "directory for the database and other state")
	logLevel := fs.String("log-level", "info", "log level: debug, info, warn, error")
	logFormat := fs.String("log-format", "text", "log format: text or json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}

	logger, err := newLogger(stderr, *logLevel, *logFormat)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*dataDir, 0o700); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ln, err := new(net.ListenConfig).Listen(ctx, "tcp", *listen)
	if err != nil {
		return err
	}
	logger.Info("starting widgora", "version", version.Version, "commit", version.Commit, "data_dir", *dataDir)
	return server.Serve(ctx, logger, ln, server.New(logger, webui.FS()))
}

func newLogger(w io.Writer, level, format string) (*slog.Logger, error) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid --log-level %q", level)
	}
	opts := &slog.HandlerOptions{Level: lvl}
	switch format {
	case "text":
		return slog.New(slog.NewTextHandler(w, opts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, opts)), nil
	default:
		return nil, fmt.Errorf("invalid --log-format %q", format)
	}
}
