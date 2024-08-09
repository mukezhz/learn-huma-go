package hello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func getDataSchemas(api huma.API, eventTypeMap map[string]any) []*huma.Schema {
	typeToEvent := make(map[reflect.Type]string, len(eventTypeMap))
	dataSchemas := make([]*huma.Schema, 0, len(eventTypeMap))
	for k, v := range eventTypeMap {
		vt := deref(reflect.TypeOf(v))
		typeToEvent[vt] = k
		required := []string{"data"}
		if k != "" && k != "message" {
			required = append(required, "event")
		}
		s := &huma.Schema{
			Title: "Event " + k,
			Type:  huma.TypeObject,
			Properties: map[string]*huma.Schema{
				"id": {
					Type:        huma.TypeInteger,
					Description: "The event ID.",
				},
				"event": {
					Type:        huma.TypeString,
					Description: "The event name.",
					Extensions: map[string]interface{}{
						"const": k,
					},
				},
				"data": api.OpenAPI().Components.Schemas.Schema(vt, true, k),
				"retry": {
					Type:        huma.TypeInteger,
					Description: "The retry time in milliseconds.",
				},
			},
			Required: required,
		}

		dataSchemas = append(dataSchemas, s)
	}
	return dataSchemas
}

func send(ctx huma.Context, data SSEvent) error {
	humaCtx := ctx.(huma.Context)
	j, err := json.Marshal(data.Data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(humaCtx.BodyWriter(), "event: %s\n", data.Event); err != nil {
		fmt.Printf("Error writing event type: %v\n", err)
		return err
	}
	message := fmt.Sprintf("data: %s\n\n", j)
	if _, err := fmt.Fprintf(humaCtx.BodyWriter(), message); err != nil {
		return err
	}

	humaCtx.BodyWriter().(http.Flusher).Flush()
	return err
}

func getSSEOperation(api huma.API, eventTypeMap map[string]any) huma.Operation {
	dataSchemas := getDataSchemas(api, eventTypeMap)
	schema := &huma.Schema{
		Title:       "Server Sent Events",
		Description: "Each oneOf object in the array represents one possible Server Sent Events (SSE) message, serialized as UTF-8 text according to the SSE specification.",
		Type:        huma.TypeArray,
		Items: &huma.Schema{
			Extensions: map[string]interface{}{
				"oneOf": dataSchemas,
			},
		},
	}
	op := huma.Operation{
		Method:      http.MethodGet,
		Description: "Event Stream",
		Tags:        []string{"sse"},
		Path:        "/sse",
		Summary:     "Stream events to the client",
	}
	if op.Responses == nil {
		op.Responses = map[string]*huma.Response{}
	}
	if op.Responses["200"] == nil {
		op.Responses["200"] = &huma.Response{}
	}
	if op.Responses["200"].Content == nil {
		op.Responses["200"].Content = map[string]*huma.MediaType{}
	}
	op.Responses["200"].Content["text/event-stream"] = &huma.MediaType{
		Schema: schema,
	}
	return op
}
