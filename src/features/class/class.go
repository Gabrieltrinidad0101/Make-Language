package class

import (
	"makeLanguages/src/languageContext"
)

type ClassBase interface {
	GetClassContext() languageContext.Context
	GetClassName() string
}

type Class struct {
	Context *languageContext.Context
	Name    string
}
