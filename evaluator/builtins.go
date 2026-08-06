package evaluator

import "revop/object"

var builtins map[string]*object.Builtin

// expectArray asserts that arg is an Array, returning an error object described
// by label (e.g. "argument to `sum`") when it is not.
func expectArray(arg object.Object, label string) (*object.Array, object.Object) {
	if arg.Type() != object.ARRAY {
		return nil, object.NewError("%s must be ARRAY, got %s", label, arg.Type())
	}
	return arg.(*object.Array), nil
}

// expectFunction asserts that arg is a Function, returning an error object
// described by label when it is not.
func expectFunction(arg object.Object, label string) (*object.Function, object.Object) {
	if arg.Type() != object.FUNCTION {
		return nil, object.NewError("%s must be FUNCTION, got %s", label, arg.Type())
	}
	return arg.(*object.Function), nil
}

// expectParams asserts that the callback passed to builtin name has exactly
// count parameters, returning an error object when it does not.
func expectParams(fn *object.Function, name string, count int) object.Object {
	if len(fn.Parameters) != count {
		return object.NewError("function passed to `%s` must have exactly %d parameters, got %d", name, count, len(fn.Parameters))
	}
	return nil
}

// evalCallback applies a user function to args in a fresh enclosed environment,
// binding each parameter positionally.
func evalCallback(fn *object.Function, args ...object.Object) object.Object {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, arg := range args {
		env.Set(fn.Parameters[i].Value, arg)
	}
	return Eval(fn.Body, env)
}

func init() {
	builtins = map[string]*object.Builtin{
		"len": {
			Args: 1,
			Fn: func(args ...object.Object) object.Object {
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return object.NewError("argument to `len` not supported, got %s", arg.Type())
				}
			},
		},
		"push": {
			Args: 2,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "first argument to `push`")
				if err != nil {
					return err
				}

				length := len(arr.Elements)

				newElements := make([]object.Object, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			},
		},
		"map": {
			Args: 2,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "first argument to `map`")
				if err != nil {
					return err
				}
				fn, err := expectFunction(args[1], "second argument to `map`")
				if err != nil {
					return err
				}
				if err := expectParams(fn, "map", 1); err != nil {
					return err
				}

				newElements := make([]object.Object, len(arr.Elements))
				for i, elem := range arr.Elements {
					evaluated := evalCallback(fn, elem)
					if isError(evaluated) {
						return evaluated
					}
					newElements[i] = evaluated
				}

				return &object.Array{Elements: newElements}
			},
		},
		"filter": {
			Args: 2,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "first argument to `filter`")
				if err != nil {
					return err
				}
				fn, err := expectFunction(args[1], "second argument to `filter`")
				if err != nil {
					return err
				}
				if err := expectParams(fn, "filter", 1); err != nil {
					return err
				}

				var newElements []object.Object
				for _, elem := range arr.Elements {
					evaluated := evalCallback(fn, elem)
					if isError(evaluated) {
						return evaluated
					}

					result, ok := evaluated.(*object.Boolean)
					if !ok {
						return object.NewError("function passed to `filter` must return BOOLEAN, got %s", evaluated.Type())
					}
					if result == TRUE {
						newElements = append(newElements, elem)
					}
				}

				return &object.Array{Elements: newElements}
			},
		},
		"reduce": {
			Args: 3,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "first argument to `reduce`")
				if err != nil {
					return err
				}
				fn, err := expectFunction(args[1], "second argument to `reduce`")
				if err != nil {
					return err
				}
				accumulator := args[2]

				if err := expectParams(fn, "reduce", 2); err != nil {
					return err
				}

				for _, elem := range arr.Elements {
					evaluated := evalCallback(fn, accumulator, elem)
					if isError(evaluated) {
						return evaluated
					}
					accumulator = evaluated
				}

				return accumulator
			},
		},
		"pop": {
			Args: 1,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "argument to `pop`")
				if err != nil {
					return err
				}

				length := len(arr.Elements)

				if length == 0 {
					return &object.Array{Elements: []object.Object{}}
				}

				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[:length-1])

				return &object.Array{Elements: newElements}
			},
		},
		"sum": {
			Args: 1,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "argument to `sum`")
				if err != nil {
					return err
				}

				var sum int64 = 0

				for _, elem := range arr.Elements {
					switch e := elem.(type) {
					case *object.Integer:
						sum += e.Value
					default:
						return object.NewError("all elements in array must be INTEGER for `sum`, got %s", elem.Type())
					}
				}

				return &object.Integer{Value: sum}
			},
		},
		"reverse": {
			Args: 1,
			Fn: func(args ...object.Object) object.Object {
				arr, err := expectArray(args[0], "argument to `reverse`")
				if err != nil {
					return err
				}

				length := len(arr.Elements)

				newElements := make([]object.Object, length)
				for i, elem := range arr.Elements {
					newElements[length-1-i] = elem
				}

				return &object.Array{Elements: newElements}
			},
		},
	}
}
