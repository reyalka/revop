package object

type Mutability struct {
	Object
	Mutable bool
}

type Environment struct {
	store map[string]Mutability
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Mutability)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Mutability, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Mutability) Object {
	e.store[name] = val
	return val
}
