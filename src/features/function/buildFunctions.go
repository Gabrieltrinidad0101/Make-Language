package function

import (
	interpreteStructs "makeLanguages/src/interprete/structs"
)

func BuildFunctions(functionsName map[string]string) map[string]interpreteStructs.VarType {
	funcs := map[string]interpreteStructs.VarType{}
	funcs[functionsName["print"]] = interpreteStructs.VarType{
		Value:      Print{},
		IsConstant: true,
	}
	return funcs
}
