package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON)

	dynamicMiddleware := alice.New()

	mux := pat.New()

	// USERS
	mux.Post("/too", dynamicMiddleware.ThenFunc(app.tooHandler.CreateTOO))
	mux.Post("/ip", dynamicMiddleware.ThenFunc(app.ipHandler.CreateIP))
	mux.Post("/individual", dynamicMiddleware.ThenFunc(app.individualHandler.CreateIndividual))

	mux.Get("/search/too/:bin", dynamicMiddleware.ThenFunc(app.tooHandler.SearchTOOs))
	mux.Get("/search/ip/:iin", dynamicMiddleware.ThenFunc(app.ipHandler.SearchIPs))
	mux.Get("/search/individual/:iin", dynamicMiddleware.ThenFunc(app.individualHandler.SearchIndividuals))

	mux.Get("/health", dynamicMiddleware.ThenFunc(app.healthCheck))

	return standardMiddleware.Then(mux)
}
