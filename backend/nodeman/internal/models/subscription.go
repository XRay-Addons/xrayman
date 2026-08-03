package models

import "github.com/go-faster/jx"

type ClientConfigItem = jx.Raw

type SubContent = []ClientConfigItem

type SubHeader struct {
	Key   string
	Value string
}

type SubHeaders = []SubHeader
