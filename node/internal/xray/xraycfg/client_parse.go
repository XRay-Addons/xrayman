package xraycfg

import (
	"errors"
	"strings"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/go-faster/jx"
	"github.com/tidwall/gjson"
)

func parseClientConfig(in string) ([]models.ClientConfigTemplateItem, error) {
	var out []models.ClientConfigTemplateItem
	if err := jx.DecodeStr(in).Arr(func(d *jx.Decoder) error {
		cfgItem, err := d.Raw()
		if err != nil {
			return errdefs.WrapWithStack(err)
		}
		out = append(out, cfgItem)
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}

func extractVlessEmailField(cfg string) (string, error) {
	users := getUsers(gjson.Parse(cfg))
	userIDs, err := extractFields(users, "email")
	if err != nil {
		return "", err
	}
	userID, err := getSingleValue(userIDs)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func extractVlessUUIDField(cfg string) (string, error) {
	users := getUsers(gjson.Parse(cfg))
	userIDs, err := extractFields(users, "id")
	if err != nil {
		return "", err
	}
	userID, err := getSingleValue(userIDs)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// get outbounds users
func getUsers(cfgs gjson.Result) []gjson.Result {
	users := make([]gjson.Result, 0)
	cfgs.ForEach(func(_, cfg gjson.Result) bool {
		cfg.Get(`outbounds.#(protocol=="vless")#`).ForEach(func(_, o gjson.Result) bool {
			o.Get("settings.vnext").ForEach(func(_, v gjson.Result) bool {
				v.Get("users").ForEach(func(_, u gjson.Result) bool {
					users = append(users, u)
					return true
				})
				return true
			})
			return true
		})
		return true
	})
	return users
}

// extract fields by names
func extractFields(items []gjson.Result, name string) ([]string, error) {
	uniqueFields := make(map[string]struct{}, 0)
	var errs []error
	for _, item := range items {
		val := item.Get(name)
		if !val.Exists() {
			continue
		}
		field, err := extractTemplateVar(val.String())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		uniqueFields[field] = struct{}{}
	}

	// return all errors on error
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	// get fields list
	fields := make([]string, 0, len(uniqueFields))
	for f := range uniqueFields {
		fields = append(fields, f)
	}

	return fields, nil
}

func extractTemplateVar(s string) (string, error) {
	// trim spaces
	templateVar := strings.TrimSpace(s)
	// trim "{{", "}}"
	if !strings.HasPrefix(templateVar, "{{") || !strings.HasSuffix(templateVar, "}}") {
		return "", errdefs.New("invalid template format",
			errdefs.Withf("template format: %s", s))
	}
	templateVar = templateVar[2 : len(templateVar)-2]
	// trim spaces again
	templateVar = strings.TrimSpace(templateVar)
	// trim "."
	if !strings.HasPrefix(templateVar, ".") {
		return "", errdefs.New("template variable should start with '.'")
	}
	templateVar = templateVar[1:]
	if templateVar == "" {
		return "", errdefs.New("empty variable name")
	}
	return templateVar, nil
}

func getSingleValue(values []string) (string, error) {
	if len(values) > 1 {
		return "", errdefs.New("multiple name field templates found",
			errdefs.Withf("name fields: %v", values))
	}
	for _, value := range values {
		return value, nil
	}
	return "", nil
}
