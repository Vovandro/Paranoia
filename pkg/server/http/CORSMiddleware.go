package http

import (
	"context"
	"strconv"
	"strings"

	"gitlab.com/devpro_studio/Paranoia/paranoia/interfaces"
	"gitlab.com/devpro_studio/go_utils/decode"
)

type CORSMiddleware struct {
	name   string
	config CORSMiddlewareConfig
}

type CORSMiddlewareConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	MaxAge           int      `yaml:"max_age"`
}

func NewCORSMiddleware(name string) interfaces.IMiddleware {
	return &CORSMiddleware{
		name: name,
	}
}

func (c *CORSMiddleware) Init(app interfaces.IEngine, cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &c.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	// Set defaults if not provided
	if len(c.config.AllowOrigins) == 0 {
		c.config.AllowOrigins = []string{"*"}
	}

	if len(c.config.AllowMethods) == 0 {
		c.config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	if len(c.config.AllowHeaders) == 0 {
		c.config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	if c.config.MaxAge == 0 {
		c.config.MaxAge = 86400 // 24 hours
	}

	return nil
}

func (c *CORSMiddleware) Stop() error {
	return nil
}

func (c *CORSMiddleware) Name() string {
	return c.name
}

func (c *CORSMiddleware) Type() string {
	return "middleware"
}

func (c *CORSMiddleware) Invoke(next RouteFunc) RouteFunc {
	return func(ctx context.Context, httpCtx ICtx) {
		origin := httpCtx.GetRequest().GetHeader().Get("Origin")
		if origin == "" {
			// Not a CORS request, continue
			next(ctx, httpCtx)
			return
		}

		// Check if the origin is allowed
		originAllowed := false
		for _, allowedOrigin := range c.config.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				originAllowed = true
				break
			}
		}

		if !originAllowed {
			// Origin not allowed, proceed without CORS headers
			next(ctx, httpCtx)
			return
		}

		// Set CORS headers
		header := httpCtx.GetResponse().Header()
		header.Set("Access-Control-Allow-Origin", origin)

		if c.config.AllowCredentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}

		if len(c.config.ExposeHeaders) > 0 {
			header.Set("Access-Control-Expose-Headers", strings.Join(c.config.ExposeHeaders, ", "))
		}

		// Handle preflight request
		if httpCtx.GetRequest().GetMethod() == "OPTIONS" {
			if reqMethod := httpCtx.GetRequest().GetHeader().Get("Access-Control-Request-Method"); reqMethod != "" {
				// This is a preflight request
				header.Set("Access-Control-Allow-Methods", strings.Join(c.config.AllowMethods, ", "))
				header.Set("Access-Control-Allow-Headers", strings.Join(c.config.AllowHeaders, ", "))
				header.Set("Access-Control-Max-Age", strconv.Itoa(c.config.MaxAge))

				// Return 204 No Content for preflight requests
				httpCtx.GetResponse().SetStatus(204)
				return
			}
		}

		// Process actual request
		next(ctx, httpCtx)
	}
}
