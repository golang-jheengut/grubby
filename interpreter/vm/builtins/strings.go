package builtins

import (
	"errors"
	"fmt"
	"strconv"
)

type StringClass struct {
	valueStub
	classStub

	provider ClassProvider
}

func NewStringClass(classProvider ClassProvider, singletonProvider SingletonProvider) Class {
	s := &StringClass{}
	s.initialize()
	s.setStringer(s.String)

	s.provider = classProvider
	s.class = classProvider.ClassWithName("Class")
	s.superClass = classProvider.ClassWithName("Object")

	s.AddMethod(NewNativeMethod("+", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		arg := args[0].(*StringValue)
		selfAsStr := self.(*StringValue)
		return NewString(selfAsStr.value+arg.value, classProvider, singletonProvider), nil
	}))
	s.AddMethod(NewNativeMethod("==", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		asStr, ok := args[0].(*StringValue)
		if !ok {
			return singletonProvider.SingletonWithName("false"), nil
		}

		selfAsStr := self.(*StringValue)
		if selfAsStr.value == asStr.value {
			return singletonProvider.SingletonWithName("true"), nil
		} else {
			return singletonProvider.SingletonWithName("false"), nil
		}
	}))
	s.AddMethod(NewNativeMethod("<<", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		arg := args[0].(*StringValue)
		selfAsStr := self.(*StringValue)
		if selfAsStr.frozen {
			return nil, errors.New("RuntimeError: can't modify frozen String")
		}

		selfAsStr.value += arg.value
		return selfAsStr, nil
	}))
	s.AddMethod(NewNativeMethod("to_i", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		selfAsStr := self.(*StringValue)

		intValue, _ := strconv.ParseInt(selfAsStr.value, 0, 64)
		return NewFixnum(intValue, classProvider, singletonProvider), nil
	}))
	s.AddMethod(NewNativeMethod("freeze", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		selfAsStr := self.(*StringValue)
		selfAsStr.frozen = true
		return selfAsStr, nil
	}))
	s.AddMethod(NewNativeMethod("intern", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		selfAsStr := self.(*StringValue)
		maybeSymbol := singletonProvider.SymbolWithName(selfAsStr.value)
		if maybeSymbol != nil {
			return maybeSymbol, nil
		}

		symbolFromString := NewSymbol(selfAsStr.value, classProvider)
		singletonProvider.AddSymbol(symbolFromString)
		return symbolFromString, nil
	}))

	return s
}

func (c *StringClass) String() string {
	return "String"
}

func (c *StringClass) Name() string {
	return "String"
}

func (class *StringClass) New(classProvider ClassProvider, singletonProvider SingletonProvider, args ...Value) (Value, error) {
	str := &StringValue{}
	str.initialize()
	str.setStringer(str.String)
	str.setPrettyPrinter(str.PrettyPrint)
	str.class = class

	return str, nil
}

type StringValue struct {
	value string
	valueStub
	frozen bool
}

func (s *StringValue) String() string {
	return fmt.Sprintf("%s", s.value)
}

func (s *StringValue) PrettyPrint() string {
	return fmt.Sprintf("\"%s\"", s.value)
}

func (s *StringValue) RawString() string {
	return s.value
}

func NewString(str string, classProvider ClassProvider, singletonProvider SingletonProvider) Value {
	s, _ := classProvider.ClassWithName("String").New(classProvider, singletonProvider)
	s.(*StringValue).value = str
	return s
}
