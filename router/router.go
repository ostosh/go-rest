package router

import (
	"github.com/bmizerany/pat"
)

// Returns instance of HTTP request multiplexer
func New() *pat.PatternServeMux {
	router := pat.New()
	return router
}
