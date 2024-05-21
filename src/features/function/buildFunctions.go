package function

import (
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
)

func BuildFunctions(ctx *languageContext.Context, functionsName map[string]string) map[string]interpreteStructs.VarType {
	funcs := map[string]interpreteStructs.VarType{}
	funcs[functionsName["print"]] = interpreteStructs.VarType{
		Value:      NewPrint(ctx),
		IsConstant: true,
	}
	return funcs
}
