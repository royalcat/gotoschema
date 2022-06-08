package gotoschema

import (
	"errors"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type schemaCacheEntry struct {
	Name     string
	RefCount uint
	Value    typeSchema
}

type typeSchema = map[string]any

type DocGenerator struct {
	typeSchemaCache map[reflect.Type]schemaCacheEntry

	modelsSchema map[string]typeSchema
}

func NewDocGenerator() *DocGenerator {
	return &DocGenerator{
		typeSchemaCache: make(map[reflect.Type]schemaCacheEntry),
		modelsSchema:    make(map[string]typeSchema),
	}
}

func (gen *DocGenerator) AddModel(models ...any) error {
	for i := range models {
		model := models[i]
		rt := reflect.TypeOf(model)
		if rt.Kind() != reflect.Struct {
			return errors.New("bad type")
		}

		desc, err := gen.genSchemaForType(rt, false)
		if err != nil {
			return err
		}

		gen.modelsSchema[rt.Name()] = desc
	}

	return nil
}

func (gen *DocGenerator) EncodeYaml() (string, error) {
	out := &strings.Builder{}
	enc := yaml.NewEncoder(out)
	enc.SetIndent(4)
	enc.Encode(gen.renderDocDict())

	return out.String(), nil
}

func (gen *DocGenerator) renderDocDict() map[string]typeSchema {
	schemas := map[string]typeSchema{}
	for m, s := range gen.modelsSchema {
		schemas[m] = s
	}
	for m, s := range gen.getRefedSchemas() {
		schemas[m] = s
	}
	return schemas
}

func (gen *DocGenerator) getRefedSchemas() map[string]typeSchema {
	out := map[string]typeSchema{}
	for _, v := range gen.typeSchemaCache {
		if v.RefCount > 1 {
			out[v.Name] = v.Value
		}
	}

	return out
}
