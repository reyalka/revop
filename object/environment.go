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
	obj, _, ok := e.GetWithEnv(name)
	return obj, ok
}

func (e *Environment) GetWithEnv(name string) (Mutability, *Environment, bool) {
	obj, ok := e.store[name]
	if ok {
		return obj, e, true
	}
	if e.outer != nil {
		obj, outerEnv, ok := e.outer.GetWithEnv(name)
		if ok {
			return obj, outerEnv, true
		}
	}
	return Mutability{}, nil, false
}

func (e *Environment) Set(name string, val Mutability) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Assign(name string, val Mutability) Object {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val
	}
	if e.outer != nil {
		if _, ok := e.outer.store[name]; ok {
			e.outer.store[name] = val
			return val
		}
	}
	e.store[name] = val
	return val
}
