package define

import (
	"strings"

	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

//go:embed skeleton.json
var docSkeleton []byte

type DocModifier func() (doc *openapi3.T, loaded bool, error error)

func Doc(data []byte) DocModifier {
	if data != nil {
		return func() (*openapi3.T, bool, error) {
			doc, err := openapi3.NewLoader().LoadFromData(data)
			return doc, true, err
		}
	}
	return func() (*openapi3.T, bool, error) {
		doc, err := reflectopenapi.NewDocFromSkeleton(docSkeleton)
		return doc, false, err
	}
}

func (m DocModifier) After(f func(doc *openapi3.T)) DocModifier {
	return func() (*openapi3.T, bool, error) {
		doc, loaded, err := m()
		if loaded {
			return doc, loaded, err
		}
		if err != nil {
			return doc, loaded, err
		}
		f(doc)
		return doc, loaded, nil
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
