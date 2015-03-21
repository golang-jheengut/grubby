package vm_test

import (
	"os"
	"path/filepath"

	. "github.com/grubby/grubby/interpreter/vm"
	. "github.com/grubby/grubby/interpreter/vm/builtins"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Strings", func() {
	var vm VM

	BeforeEach(func() {
		pathToExecutable, err := filepath.Abs(filepath.Dir(filepath.Dir(filepath.Dir(os.Args[0]))))
		if err != nil {
			panic(err)
		}

		vm = NewVM(pathToExecutable, "fake-irb-under-test")
	})

	It("can be constructed from a string literal", func() {
		val, err := vm.Run("'nonrestricted-consonantize'")

		Expect(err).ToNot(HaveOccurred())
		Expect(val).To(BeAssignableToTypeOf(NewString("", vm, vm)))
		Expect(val.String()).To(Equal("nonrestricted-consonantize"))
	})

	It("has a + method", func() {
		val, err := vm.Run("'foo' + 'bar'")

		Expect(err).ToNot(HaveOccurred())
		Expect(val.String()).To(Equal(NewString("foobar", vm, vm).String()))
	})

	Describe("interpolating values inside of a string", func() {
		It("can be done with double quoted strings", func() {
			value, err := vm.Run(`
greetings = 'hello'
adj       = 'cruel'
place     = 'world'
"#{greetings} #{adj} #{place}"
`)
			Expect(err).ToNot(HaveOccurred())
			Expect(value.(*StringValue).RawString()).To(Equal("hello cruel world"))
		})

		It("cannot be done with single quoted strings", func() {
			value, err := vm.Run(`
adj = 'cruel'
'hello #{adj} world'
`)

			Expect(err).ToNot(HaveOccurred())
			Expect(value.(*StringValue).RawString()).To(Equal("hello #{adj} world"))
		})
	})

	Describe("concatenating strings", func() {
		It("can be done with the shovel operator", func() {
			value, err := vm.Run("'hello' << ' world'")
			Expect(err).ToNot(HaveOccurred())
			Expect(value.(*StringValue).RawString()).To(Equal("hello world"))
		})
	})

	Describe("#freeze", func() {
		It("can be used to turn a string into an immutable string", func() {
			_, err := vm.Run("str = 'hello'.freeze; str << 'world'")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("RuntimeError: can't modify frozen String"))
		})
	})

	Describe("#intern", func() {
		It("can be used to create a symbol with string's value as its name", func() {
			value, err := vm.Run("'hello world'.intern")
			Expect(err).ToNot(HaveOccurred())

			Expect(value).ToNot(BeNil())
			Expect(value).To(BeAssignableToTypeOf(&SymbolValue{}))
			Expect(value).To(Equal(vm.SymbolWithName("hello world"))) // singleton instance
		})
	})
})