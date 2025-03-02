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

	mux.Put("/too", dynamicMiddleware.ThenFunc(app.tooHandler.UpdateUserContract))
	mux.Put("/ip", dynamicMiddleware.ThenFunc(app.ipHandler.UpdateUserContract))
	mux.Put("/individual", dynamicMiddleware.ThenFunc(app.individualHandler.UpdateUserContract))

	// Поиск договоров на главной странице
	mux.Get("/search/too/:iin/password/:pass", dynamicMiddleware.ThenFunc(app.tooHandler.SearchTOOs))
	mux.Get("/search/ip/:iin/password/:pass", dynamicMiddleware.ThenFunc(app.ipHandler.SearchIPs))
	mux.Get("/search/individual/:iin/password/:pass", dynamicMiddleware.ThenFunc(app.individualHandler.SearchIndividuals))
	mux.Get("/search/all/:iin/password/:pass", dynamicMiddleware.ThenFunc(app.companyDataHandler.GetAllDataByIIN))

	// Поиск договоров по QR коду
	mux.Get("/search/too/token/:token", dynamicMiddleware.ThenFunc(app.tooHandler.SearchTOOsByToken))
	mux.Get("/search/ip/token/:token", dynamicMiddleware.ThenFunc(app.ipHandler.SearchIPsByToken))
	mux.Get("/search/individual/token/:token", dynamicMiddleware.ThenFunc(app.individualHandler.SearchIndividualsByToken))

	mux.Get("/search/too/id/:id", dynamicMiddleware.ThenFunc(app.tooHandler.SearchTOOsByID))
	mux.Get("/search/ip/id/:id", dynamicMiddleware.ThenFunc(app.ipHandler.SearchIPsByID))
	mux.Get("/search/individual/id/:id", dynamicMiddleware.ThenFunc(app.individualHandler.SearchIndividualsByID))

	mux.Put("/too/:id", dynamicMiddleware.ThenFunc(app.tooHandler.UpdateUserContractStatus))
	mux.Put("/ip/:id", dynamicMiddleware.ThenFunc(app.ipHandler.UpdateUserContractStatus))
	mux.Put("/individual/:id", dynamicMiddleware.ThenFunc(app.individualHandler.UpdateUserContractStatus))

	// создание новой компании с паролем
	mux.Post("/companies", dynamicMiddleware.ThenFunc(app.companyHandler.Create))
	mux.Post("/companies/id/:id/pass/:pass", dynamicMiddleware.ThenFunc(app.companyHandler.CheckPassword))

	// Полный цикл для расторжения договора
	mux.Post("/discard", dynamicMiddleware.ThenFunc(app.discardHandler.CreateDiscard))
	mux.Put("/discard", dynamicMiddleware.ThenFunc(app.discardHandler.UpdateContractPath))

	return standardMiddleware.Then(mux)
}
