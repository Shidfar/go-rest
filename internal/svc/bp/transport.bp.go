package bp

import (
	"text/template"
)

var transportTemplate = `// AUTOGENERATED CODE.
// !!  DO NOT EDIT  !!

package {{ .Pkg }}http

import (
	"context"
	"encoding/json"
	"errors"
	"{{ .PkgPath }}"
	"net/http"
)

var ErrBadRequest = errors.New("bad request")

//var ErrInvalidId = errors.New("invalid id")
{{ range $func := .Funcs }}
func encode{{ .Name }}Response(_ context.Context, writer http.ResponseWriter, rawRes any) error {
	{{ range $res := .Returns }}
	var {{ .Name }} {{ .Type.Type "DynamicField" }}
	{{ .Name }} = rawRes.({{ .Type.Type "DynamicField" }})
	
	//res := response.(*getAllResponse)
	//if res.err != nil {
	//	return json.NewEncoder(w).Encode(res.err)
	//}
	//return json.NewEncoder(w).Encode(res.payload)
	header := writer.Header()
	header.Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode({{ .Name }}); err != nil {
		return err
	}
	{{ end }}
	return nil
}
func decode{{ .Name }}Request(_ context.Context, httpReq *http.Request) (any, error) {
	{{ range $args := .Arguments }}
	var {{ .Name }} {{ .Type.Type "DynamicField" }}

	if err := decodeRequest(httpReq, &{{ .Name }}); err != nil {
		return nil, ErrBadRequest
	}

	{{ end }}

	//if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
	//	//return nil, ErrBadRequest
	//	return nil, fmt.Errorf("could not decode: %w", err)
	//}
	return req, nil
}

{{ end }}
func decodeRequest(httpReq *http.Request, req any) error {
	switch httpReq.Method {
	case http.MethodPost:
		if httpReq.Body == nil {
			return nil
		}
		return json.NewDecoder(httpReq.Body).Decode(req)
	//case http.MethodGet:
	//	if err := httpReq.ParseForm(); err != nil {
	//		return err
	//	}
	//	return decodeURLValues(httpReq.Form, vv)
	default:
		return nil
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	//case ErrInvalidId:
	//	w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
`

func NewTransportTemplate(serviceName string) *template.Template {
	t := template.New(serviceName)
	return template.Must(t.Parse(transportTemplate))
}
