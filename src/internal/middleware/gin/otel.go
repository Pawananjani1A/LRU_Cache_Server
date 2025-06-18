/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package gin

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func OtelMiddleware(tracerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := otel.GetTracerProvider().Tracer(tracerName)
		ctx, span := tr.Start(c.Request.Context(), c.Request.URL.Path)
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.Int("http.status_code", c.Writer.Status()),
		)
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors[0])
		}
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
