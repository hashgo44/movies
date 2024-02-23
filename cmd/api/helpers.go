package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Aller a la ligne pour faciliter la visualisation dans le terminal.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// Cette méthode va servir à lire les requêtes entrantes dans l'application et formater les erreurs
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Utilisez http.MaxBytesReader() pour limiter la taille du corps de la requête à 1 Mo.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialisez le json.Decoder et appelez la méthode DisallowUnknownFields() dessus avant le décodage.
	// Cela signifie que si le JSON du client inclut désormais un champ qui ne peut pas être mappé à la destination cible, le décodeur renverra une erreur au lieu de simplement ignorer le champ.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Décodez le corps de la requête dans la destination cible.
	err := dec.Decode(dst)

	if err != nil {
		// S'il y a une erreur lors du décodage, récupérez-la
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarchalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Utilisez la fonction error.As() pour vérifier si l’erreur est du type *json.SyntaxError. Si c'est le cas, renvoyez un message d'erreur en anglais simple qui inclut l'emplacement du problème.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Decode() peut également renvoyer une erreur io.ErrUnexpectedEOF pour les erreurs de syntaxe dans le JSON. Nous vérifions donc cela à l’aide de error.Is() et renvoyons un message d’erreur générique.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// *json.UnmarshalTypeError. Cette erreur se produit lorsque la valeur JSON n'est pas du bon type pour la destination cible. Si l'erreur concerne un champ spécifique, nous l'incluons dans notre message d'erreur pour faciliter le débogage du client.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for filed %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Une erreur io.EOF sera renvoyée par Decode() si le corps de la requête est vide.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Si le JSON contient un champ inconnu nous renvoyons une erreur
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Utilisez la fonction error.As() pour vérifier si l’erreur est du type *http.MaxBytesError. Si tel est le cas, cela signifie que le corps de la requête a dépassé nos limites. limite de taille de 1 Mo et nous renvoyons un message d'erreur clair.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// Une erreur json.InvalidUnmarshalError sera renvoyée si nous transmettons quelque chose qui n'est pas un pointeur non nul vers Decode(). Nous captons cela utilisons panic(), plutôt que de renvoyer une erreur à notre gestionnaire.
		case errors.As(err, &invalidUnmarchalError):
			panic(err)

		// Pour toute autre chose, renvoyez le message d'erreur tel quel.
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
