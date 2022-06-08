/*
 le package env regroupe toutes les méthodes lié aux variables d'environement et au systeme.
*/
package env

import (
	"fmt"
	"os"
	"strconv"
)

// DirFromEnvironment récupère le path d'un répertoire depuis une variable d'environnement passé en paramètre.
func DirFromEnvironment(envName string) (string, error) {
	dir, ok := os.LookupEnv(envName)
	if !ok || dir == "" {
		return "", fmt.Errorf("la variable `%s` est manquante", envName)
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !dirInfo.IsDir() {
		return "", fmt.Errorf("`%s` n'est pas un répertoire valide", dir)
	}
	return dir, nil
}

// RequiredStringFromEnvironment vérifie si la variable d'environement est renseignée et retourne sa valeur.
func RequiredStringFromEnvironment(envName string) (string, error) {
	if env, ok := os.LookupEnv(envName); ok && env != "" {
		return env, nil
	}
	return "", fmt.Errorf("la variable `%s` est manquante", envName)
}

// StringFromEnvironment récupère la valeur d'une variable d'environement ou si non définie la valeur par defaut
func StringFromEnvironment(envName, defaultValue string) string {
	if env, ok := os.LookupEnv(envName); ok && env != "" {
		return env
	}
	return defaultValue
}

// BoolFromEnvironment récupère la valeur d'une variable d'environement au format booleene sinon la valeur par defaut
func BoolFromEnvironment(envName string, defaultValue bool) bool {
	if env, ok := os.LookupEnv(envName); ok && env != "" {
		ret, err := strconv.ParseBool(env)
		if err != nil {
			return defaultValue
		}
		return ret
	}
	return defaultValue
}
