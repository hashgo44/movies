package main

import (
	"fmt"
	"net/http"
)

//	La méthode logError() est une aide générique pour enregistrer un message d'erreur.
//
// avec la méthode de requête actuelle et l'URL en tant qu'attributs dans les logs.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// La méthode errorResponse() est une aide générique pour envoyer des messages d'erreur au format JSON au client avec un code d'état donné.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// La méthode serverErrorResponse() sera utilisée lorsque notre application rencontre un problème inattendu au moment de l'exécution.
// Il enregistre le message d'erreur détaillé, puis utilise l'assistant errorResponse() pour envoyer un code d'état d'erreur interne du serveur 500 et une réponse JSON (contenant un message d'erreur générique) au client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and JSON response to the client. La méthode notFoundResponse() sera utilisée pour envoyer un code d'état 404 Not Found et Réponse JSON au client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// La méthode methodNotAllowedResponse() sera utilisée pour envoyer une méthode 405 non autorisée code d'état et réponse JSON au client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Méthode qui permet de valider les données lors d'une création
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// Méthode qui est renvoyée s'il y a une erreur de conflit
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
