package function

func BuildFunctions(functionsName map[string]string) map[string]IFunction {
	funcs := map[string]IFunction{}
	funcs[functionsName["print"]] = Print{}
	return funcs
}
