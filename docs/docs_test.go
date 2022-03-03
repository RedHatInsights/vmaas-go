package docs

import (
	"io/ioutil"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

const openAPIPath = "openapi.json"

func TestValidateOpenAPI3DocStr(t *testing.T) {
	doc, err := ioutil.ReadFile(openAPIPath)
	assert.Nil(t, err)
	_, err = openapi3.NewSwaggerLoader().LoadSwaggerFromData(doc)
	assert.Nil(t, err)
}
