package reflect

import "reflect"

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func
		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn)
		inputValue := make([]reflect.Value, 0, numIn)
		inputValue = append(inputValue, reflect.ValueOf(entity))
		input = append(input, reflect.TypeOf(entity))
		// index为0指向接收器
		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			input = append(input, fnInType)
			inputValue = append(inputValue, reflect.Zero(fnInType))
		}
		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for m := 0; m < numOut; m++ {
			output = append(output, fn.Type().Out(m))
		}
		resValues := fn.Call(inputValue)
		result := make([]any, 0, numOut)
		for k := 0; k < numOut; k++ {
			result = append(result, resValues[k].Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}
	return res, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
