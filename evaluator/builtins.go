package evaluator

import "revop/object"

var builtins map[string]*object.Builtin

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
				if args[0].Type() != object.ARRAY {
					return object.NewError("first argument to `push` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("first argument to `map` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION {
					return object.NewError("second argument to `map` must be FUNCTION, got %s", args[1].Type())
				}

				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)

				if len(fn.Parameters) != 1 {
					return object.NewError("function passed to `map` must have exactly 1 parameters, got %d", len(fn.Parameters))
				}

				newElements := make([]object.Object, len(arr.Elements))
				for i, elem := range arr.Elements {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set(fn.Parameters[0].Value, elem)

					evaluated := Eval(fn.Body, extendedEnv)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("first argument to `filter` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION {
					return object.NewError("second argument to `filter` must be FUNCTION, got %s", args[1].Type())
				}

				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)

				if len(fn.Parameters) != 1 {
					return object.NewError("function passed to `filter` must have exactly 1 parameters, got %d", len(fn.Parameters))
				}

				var newElements []object.Object
				for _, elem := range arr.Elements {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set(fn.Parameters[0].Value, elem)

					evaluated := Eval(fn.Body, extendedEnv)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("first argument to `reduce` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION {
					return object.NewError("second argument to `reduce` must be FUNCTION, got %s", args[1].Type())
				}

				arr := args[0].(*object.Array)
				fn := args[1].(*object.Function)
				accumulator := args[2]

				if len(fn.Parameters) != 2 {
					return object.NewError("function passed to `reduce` must have exactly 2 parameters, got %d", len(fn.Parameters))
				}

				for _, elem := range arr.Elements {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set(fn.Parameters[0].Value, accumulator)
					extendedEnv.Set(fn.Parameters[1].Value, elem)

					evaluated := Eval(fn.Body, extendedEnv)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("argument to `pop` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("argument to `sum` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
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
				if args[0].Type() != object.ARRAY {
					return object.NewError("argument to `reverse` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
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
