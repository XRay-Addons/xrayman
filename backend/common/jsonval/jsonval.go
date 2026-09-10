package jsonval

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/go-faster/jx"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateJsonData(data []byte) error {
	d := jx.DecodeBytes(data)
	if err := d.Skip(); err != nil {
		return xerr.WrapWithInfo(err, "json validation")
	}

	if d.Next() != jx.Invalid {
		return xerr.Newf("validation: text after end: %v", d.Next().String())
	}
	return nil
}

func ValidateJsonSchema(data, schemaData []byte) error {
	var schemaDoc any
	if err := json.Unmarshal(schemaData, &schemaDoc); err != nil {
		return xerr.WrapWithInfo(err, "parse JSON schema")
	}

	schema, err := jsonschema.CompileString("schema.json", string(schemaData))
	if err != nil {
		return xerr.WrapWithInfo(err, "compile JSON schema")
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return xerr.WrapWithInfo(err, "parse JSON")
	}

	if err := schema.Validate(value); err != nil {
		var ve *jsonschema.ValidationError
		if errors.As(err, &ve) {
			causes := []*jsonschema.ValidationError{ve}
			if len(ve.Causes) > 0 {
				causes = ve.Causes
			}
			descriptions := make([]string, len(causes))
			for i, c := range causes {
				descriptions[i] = fmt.Sprintf("%s validation rule %s violated: %s",
					c.InstanceLocation, c.KeywordLocation, c.Message)
			}
			if len(descriptions) == 1 {
				return xerr.Newf("JSON schema validation: %s", descriptions[0])
			}
			return xerr.Newf("JSON schema validation:\n\t-%s", strings.Join(descriptions, "\n\t-"))
		}

		return xerr.WrapWithInfo(err, "JSON schema validation")
	}

	return nil
}
