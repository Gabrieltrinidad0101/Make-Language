package class

import "makeLanguages/src/languageContext"

type Class struct {
	Methods interface{}
	Context languageContext.Context
}
