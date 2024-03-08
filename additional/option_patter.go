package additional

type MyStructOption func(m *MyStruct)

type MyStructOptionErr func(m *MyStruct) error

type MyStruct struct {
	id      uint64
	name    string
	address string
}

func NewMyStruct(id uint64, name string, opts ...MyStructOption) MyStruct {
	res := MyStruct{
		id:   id,
		name: name,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}

func AddAddress(address string) MyStructOption {
	return func(m *MyStruct) {
		m.address = address
	}
}

//func NewMyStruct1(id uint64, name string, opts ...MyStructOptionErr) (*MyStruct, error) {
//	res := &MyStruct{
//		id:   id,
//		name: name,
//	}
//	for _, opt := range opts {
//		if err := opt(res); err != nil {
//			return res, err
//		}
//	}
//	return res, nil
//}
//
//func AddAddress(address string) MyStructOptionErr {
//	return func(m *MyStruct) error {
//		if address == "" {
//			return errors.New("地址不能为空字符串")
//		}
//		m.address = address
//		return nil
//	}
//}
