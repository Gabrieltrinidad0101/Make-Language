package function

import (
	"makeLanguages/src/languageContext"
	"makeLanguages/src/token"
)

type Function struct {
	Context languageContext.Context
	Params  *[]token.Token
	Body    interface{}
}

func (func_ Function) SetParams(params *[]interface{}) {
	i := 0
	for _, token := range *func_.Params {
		func_.Context.Set(token.Value.(string), (*params)[i])
		i++
	}
}
