package clientcfg

import (
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

func parseClientConfig(in string) ([]models.ClientConfigTemplateItem, error) {
	arr := gjson.Parse(in)
	if !arr.IsArray() {
		return nil, xerr.New("client config must be an array")
	}
	out := make([]models.ClientConfigTemplateItem, 0, len(arr.Array()))
	for _, item := range arr.Array() {
		out = append(out, item.Raw)
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
		return nil, xerr.Join(errs...)
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
		return "", xerr.Newf("invalid template format: %s", s)
	}
	templateVar = templateVar[2 : len(templateVar)-2]
	// trim spaces again
	templateVar = strings.TrimSpace(templateVar)
	// trim "."
	if !strings.HasPrefix(templateVar, ".") {
		return "", xerr.New("template variable should start with '.'")
	}
	templateVar = templateVar[1:]
	if templateVar == "" {
		return "", xerr.New("empty variable name")
	}
	return templateVar, nil
}

func getSingleValue(values []string) (string, error) {
	if len(values) > 1 {
		return "", xerr.Newf("multiple name field templates found: %v", values)
	}
	for _, value := range values {
		return value, nil
	}
	return "", nil
}
