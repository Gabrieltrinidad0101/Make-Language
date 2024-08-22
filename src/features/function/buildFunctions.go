package function

import (
	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"
)

func BuildFunctions(ctx *languageContext.Context, functionsName map[string]string) map[string]interpreteStructs.VarType {
	funcs := map[string]interpreteStructs.VarType{}
	funcs[functionsName["print"]] = interpreteStructs.VarType{
		Value:      NewPrint(ctx),
		IsConstant: true,
	}
	return funcs
}
