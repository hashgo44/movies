package validator

import (
	"regexp"
	"slices"
)

// Définition d'un regex pour validé vérifier si la saisi utilisateur est bien un email
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Définition d'un struct qui contient un map de validation d'erreurs
type Validator struct {
	Errors map[string]string
}

// Crée une nouvelle instance de Validator avec une carte d'erreurs vide.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid renvoie true si la carte des erreurs ne contient aucune entrée.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError ajoute un message d'erreur à la map (tant qu'aucune entrée n'existe déjà pour la clé donnée).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check ajoute un message d'erreur à la carte uniquement si un contrôle de validation n'est pas « ok ».
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Function générique qui renvoie vrai si une valeur spécifique se trouve dans une liste de valeurs autorisées.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches renvoie vrai si un string correspond à un regex.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Function générique qui renvoie vrai si toutes les valeurs d'une slice sont uniques.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
