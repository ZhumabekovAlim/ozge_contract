package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, makeResponseJSON, app.logExecutionTime, gzipMiddleware)

	//dynamicMiddleware := alice.New()

	mux := pat.New()

	// USERS
	mux.Post("/too", standardMiddleware.ThenFunc(app.tooHandler.CreateTOO))
	mux.Post("/ip", standardMiddleware.ThenFunc(app.ipHandler.CreateIP))
	mux.Post("/individual", standardMiddleware.ThenFunc(app.individualHandler.CreateIndividual))

	mux.Put("/too", standardMiddleware.ThenFunc(app.tooHandler.UpdateUserContract))
	mux.Put("/ip", standardMiddleware.ThenFunc(app.ipHandler.UpdateUserContract))
	mux.Put("/individual", standardMiddleware.ThenFunc(app.individualHandler.UpdateUserContract))

	// Поиск договоров на главной странице
	mux.Get("/search/too/:iin/password/:pass/id/:id", standardMiddleware.ThenFunc(app.tooHandler.SearchTOOs))
	mux.Get("/search/ip/:iin/password/:pass/id/:id", standardMiddleware.ThenFunc(app.ipHandler.SearchIPs))
	mux.Get("/search/individual/:iin/password/:pass/id/:id", standardMiddleware.ThenFunc(app.individualHandler.SearchIndividuals))
	mux.Get("/search/all/:iin/password/:pass/id/:id", standardMiddleware.ThenFunc(app.companyDataHandler.GetAllDataByIIN))

	// Поиск договоров по QR коду
	mux.Get("/search/too/token/:token", standardMiddleware.ThenFunc(app.tooHandler.SearchTOOsByToken))
	mux.Get("/search/ip/token/:token", standardMiddleware.ThenFunc(app.ipHandler.SearchIPsByToken))
	mux.Get("/search/individual/token/:token", standardMiddleware.ThenFunc(app.individualHandler.SearchIndividualsByToken))

	mux.Get("/search/too/id/:id", standardMiddleware.ThenFunc(app.tooHandler.SearchTOOsByID))
	mux.Get("/search/ip/id/:id", standardMiddleware.ThenFunc(app.ipHandler.SearchIPsByID))
	mux.Get("/search/individual/id/:id", standardMiddleware.ThenFunc(app.individualHandler.SearchIndividualsByID))

	mux.Put("/too/:id", standardMiddleware.ThenFunc(app.tooHandler.UpdateUserContractStatus))
	mux.Put("/ip/:id", standardMiddleware.ThenFunc(app.ipHandler.UpdateUserContractStatus))
	mux.Put("/individual/:id", standardMiddleware.ThenFunc(app.individualHandler.UpdateUserContractStatus))

	// создание новой компании с паролем
	mux.Post("/companies", standardMiddleware.ThenFunc(app.companyHandler.Create))
	mux.Post("/companies/id/:id/pass/:pass", standardMiddleware.ThenFunc(app.companyHandler.CheckPassword))

	// Полный цикл для расторжения договора
	mux.Post("/discard", standardMiddleware.ThenFunc(app.discardHandler.CreateDiscard))
	mux.Put("/discard", standardMiddleware.ThenFunc(app.discardHandler.UpdateContractPath))

	return standardMiddleware.Then(mux)
}
