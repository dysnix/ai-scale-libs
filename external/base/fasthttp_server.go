package base

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
	"net/http/pprof"
	rtp "runtime/pprof"

	libs "github.com/dysnix/ai-scale-libs/external/configs"
	libsSrv "github.com/dysnix/ai-scale-libs/external/grpc/server"
)

type fastHttpLogger struct {
	zap.SugaredLogger
}

func (l *fastHttpLogger) Printf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

var (
	_ libsSrv.Server = (*FastHttpServer)(nil)
)

type FastHttpServer struct {
	conf   *libs.Base
	logger *fastHttpLogger
	server *fasthttp.Server
}

func NewFastHttpServer(conf *libs.Base, lg *zap.SugaredLogger) (out *FastHttpServer) {
	defer func() {
		out.init()
	}()

	return &FastHttpServer{
		conf: conf,
		logger: &fastHttpLogger{
			SugaredLogger: *lg,
		},
	}
}

func (s *FastHttpServer) routing() *router.Router {
	r := router.New()

	if s.conf.Monitoring.Enabled {
		// metrics
		r.ANY("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))
	}

	if s.conf.Profiling.Enabled {
		// profiling
		grPprof := r.Group("/debug/pprof")
		grPprof.ANY("/cmdline", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline))
		grPprof.ANY("/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))
		grPprof.ANY("/symbol", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol))
		grPprof.ANY("/trace", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace))
		grPprof.ANY("/{path:*}", s.indexPprofRoute)
	}

	return r
}

func (s *FastHttpServer) init() {
	s.server = &fasthttp.Server{
		Name:               s.conf.Single.Name,
		Concurrency:        int(s.conf.Single.Concurrency),
		TCPKeepalive:       s.conf.Single.TCPKeepalive.Enabled,
		TCPKeepalivePeriod: s.conf.Single.TCPKeepalive.Period,
		ReadBufferSize:     int(s.conf.Single.Buffer.ReadBufferSize),
		WriteBufferSize:    int(s.conf.Single.Buffer.WriteBufferSize),
		ReadTimeout:        s.conf.Single.HTTPTransport.ReadTimeout,
		WriteTimeout:       s.conf.Single.HTTPTransport.WriteTimeout,
		IdleTimeout:        s.conf.Single.HTTPTransport.MaxIdleConnDuration,
		Logger:             s.logger,
		Handler:            fasthttp.CompressHandler(s.PanicMiddleware(s.CorsMiddleware(s.routing().Handler))),
	}

	if s.conf.IsDebugMode {
		s.server.LogAllErrors = true
	}
}

func (s *FastHttpServer) Start() <-chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)
		if s.conf.Single != nil && s.conf.Single.Enabled {
			s.logger.Info("✔️ FataHttp server started.")
			if err := s.server.ListenAndServe(fmt.Sprintf("%s:%d", s.conf.Single.Host, s.conf.Single.Port)); err != nil {
				if s.conf.IsDebugMode {
					s.logger.Errorw(err.Error(), "serving fasthttp server with error")
				}

				errCh <- err
				return
			}
		}
	}()

	return errCh
}

func (s *FastHttpServer) Stop() error {
	defer func() {
		if s.conf.IsDebugMode {
			s.logger.Info("🛑 FataHttp server stoped.")
		}
	}()

	return s.server.Shutdown()
}

const (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func (s *FastHttpServer) CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func errorPrint(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	if err1 := json.NewEncoder(ctx).Encode(map[string]string{
		"error": err.Error(),
	}); err1 != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func (s *FastHttpServer) PanicMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					s.logger.Errorw("panic middleware detect", err)
					errorPrint(ctx, err, fasthttp.StatusInternalServerError)
				}
			}

			ctx.Done()
		}()

		next(ctx)
	}
}

func (s *FastHttpServer) indexPprofRoute(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	for _, v := range rtp.Profiles() {
		ppName := v.Name()
		if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/"+ppName) {
			namedHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler(ppName).ServeHTTP)
			namedHandler(ctx)
			return
		}
	}

	index := fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index)
	index(ctx)
}
