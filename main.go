package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	app := &cli.App{
		Name:  "updir",
		Usage: "serve static files of a directory",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bind",
				Aliases: []string{"b"},
				Value:   "0.0.0.0",
				Usage:   "bind to this address",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8899,
				Usage:   "specify port to listen",
			},
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Value:   "",
				Usage:   "serve alternate directory (default: current directory)",
			},
		},
		Action: func(ctx *cli.Context) error {
			addr := ctx.String("bind")
			port := ctx.Int("port")
			b := fmt.Sprintf("%s:%d", addr, port)
			log.Info().
				Int("port", port).
				Str("address", addr).
				Str("link", "http://"+b).
				Msg("Starting the server")
			http.Handle("/", serveDirHandler(ctx.String("directory")))
			return http.ListenAndServe(b, nil)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func serveDirHandler(dir string) http.HandlerFunc {
	h := http.FileServer(http.Dir(dir))
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		code := 200
		ww := WW{ResponseWriter: w, code: &code}
		h.ServeHTTP(ww, r)
		log.Info().
			Str("method", r.Method).
			Str("uri", r.URL.RequestURI()).
			Str("duration", time.Since(start).String()).
			Int("status_code", code).
			Msg(http.StatusText(code))
	}
}

type WW struct {
	http.ResponseWriter
	code *int
}

func (ww WW) WriteHeader(code int) {
	*ww.code = code
	ww.ResponseWriter.WriteHeader(code)
}
