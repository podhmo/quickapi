package define

import (
	"strings"

	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

//go:embed skeleton.json
var docSkeleton []byte

type DocModifier func() *openapi3.T

func Doc() DocModifier {
	doc, err := reflectopenapi.NewDocFromSkeleton(docSkeleton)
	if err != nil {
		panic(err) // skeleton template is always corect
	}
	return func() *openapi3.T { return doc }
}

func (m DocModifier) After(f func(doc *openapi3.T)) DocModifier {
	return func() *openapi3.T {
		doc := m()
		f(doc)
		return doc
	}
}

func (m DocModifier) Title(title string) DocModifier {
	return m.After(func(doc *openapi3.T) {
		doc.Info.Title = strings.TrimSpace(title)
	})
}
func (m DocModifier) Version(version string) DocModifier {
	return m.After(func(doc *openapi3.T) {
		doc.Info.Version = strings.TrimSpace(version)
	})
}
func (m DocModifier) Server(url string, description string) DocModifier {
	return m.After(func(doc *openapi3.T) {
		if len(doc.Servers) == 1 && doc.Servers[0].Description == "local development server" {
			doc.Servers = []*openapi3.Server{{URL: url, Description: description}}
			return
		}
		doc.Servers = append(doc.Servers, &openapi3.Server{URL: url, Description: description})
	})
}
