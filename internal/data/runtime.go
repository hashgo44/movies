package data

import (
	"fmt"
	"strconv"
)

// Création d'un type Runtime int32 pour notre model de données data/movies
type Runtime int32

// Création d'un méthode MarshalJSON via le type Runtime
// Si nous utilisons un type Runtime dans un model alors Go utilisera cette méthode MarshalJSON pour formater le Runtime
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}
