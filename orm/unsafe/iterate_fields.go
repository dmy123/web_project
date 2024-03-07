package unsafe

import (
	"fmt"
	"reflect"
)

func PrintFieldOffset(entity any) {
	typ := reflect.TypeOf(entity)
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		fmt.Println(fd.Offset)
	}
}
