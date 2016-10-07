package router

import (
	"github.com/bmizerany/pat"
)

func New() *pat.PatternServeMux {
	router := pat.New()
	return router
}
