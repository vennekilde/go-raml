package codegen

import (
	"fmt"
	"strings"

	"github.com/Jumpscale/go-raml/raml"
	log "github.com/Sirupsen/logrus"
)

const (
	// Oauth2 string
	Oauth2 = "OAuth 2.0"
)

// security define a security scheme, we only support oauth2 now.
// we generate middleware that checking for oauth2 credential
type security struct {
	*raml.SecurityScheme
	Name        string
	PackageName string
	Header      *raml.Header
	QueryParams *raml.NamedParameter
	//apiDef      *raml.APIDefinition
}

// create security struct
func newSecurity(ss *raml.SecurityScheme, name, packageName string) security {
	sd := security{
		SecurityScheme: ss,
	}
	sd.Name = securitySchemeName(name)
	sd.PackageName = packageName

	// assign header, if any
	for k, v := range sd.DescribedBy.Headers {
		sd.Header = &v
		sd.Header.Name = string(k)
		break
	}

	// assign query params if any
	for k, v := range sd.DescribedBy.QueryParameters {
		sd.QueryParams = &v
		sd.QueryParams.Name = string(k)
		break
	}

	return sd
}

// generate security related code
func generateSecurity(schemes map[string]raml.SecurityScheme, dir, packageName, lang string) error {
	var err error

	// generate oauth2 middleware
	for k, ss := range schemes {
		if ss.Type != Oauth2 {
			continue
		}

		sd := newSecurity(&ss, k, packageName)

		switch lang {
		case langGo:
			gss := goSecurity{security: &sd}
			err = gss.generate(dir)

		case langPython:
			pss := pythonSecurity{security: &sd}
			err = pss.generate(dir)
		}
		if err != nil {
			log.Errorf("generateSecurity() failed to generate %v, err=%v", k, err)
			return err
		}
	}
	return nil
}

// get oauth2 middleware handler from a security scheme
func getOauth2MwrHandler(ss raml.DefinitionChoice) (string, error) {
	quotedScopes, err := getQuotedSecurityScopes(ss)
	if err != nil {
		return "", err
	}
	scopesArgs := strings.Join(quotedScopes, ", ")
	return fmt.Sprintf(`newOauth2%vMiddleware([]string{%v}).Handler`, securitySchemeName(ss.Name), scopesArgs), nil
}

// get array of security scopes in the form of quoted string
func getQuotedSecurityScopes(ss raml.DefinitionChoice) ([]string, error) {
	var quoted []string
	scopes, err := getSecurityScopes(ss)
	if err != nil {
		return quoted, err
	}
	for _, s := range scopes {
		quoted = append(quoted, fmt.Sprintf(`"%v"`, s))
	}
	return quoted, nil
}

// get scopes of a security scheme as []string
func getSecurityScopes(ss raml.DefinitionChoice) ([]string, error) {
	scopes := []string{}

	// check if there is scopes
	v, ok := ss.Parameters["scopes"]
	if !ok {
		return scopes, nil
	}

	// cast it to []string, return error if failed
	scopesArr, ok := v.([]interface{})
	if !ok {
		return scopes, fmt.Errorf("scopes must be array")
	}

	// build []string
	for _, s := range scopesArr {
		scopes = append(scopes, s.(string))
	}
	return scopes, nil
}

// return security scheme name that could be used in code
func securitySchemeName(name string) string {
	return strings.Replace(name, " ", "", -1)
}

// validate security scheme:
// - not empty
// - not 'null'
// - oauth2 -> we only support oauth2 now
func validateSecurityScheme(name string, apiDef *raml.APIDefinition) bool {
	if name == "" || name == "null" {
		return false
	}
	if ss, ok := apiDef.SecuritySchemes[name]; ok {
		return ss.Type == Oauth2
	}
	return false
}
