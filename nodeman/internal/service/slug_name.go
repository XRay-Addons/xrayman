package service

import "github.com/gosimple/slug"

func makeSlugName(name string) string {
	return slug.Make(name)
}
