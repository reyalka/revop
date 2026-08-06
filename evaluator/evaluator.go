package evaluator

import (
	"math"

	"revop/ast"
	"revop/object"
)

// maxCallDepth bounds nested function calls so runaway recursion returns an
// error instead of exhausting the goroutine stack.
const maxCallDepth = 1000

var callDepth int

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return toBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.Block:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("identifier not found: %s", node.Value)
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.BOOLEAN {
		return object.NewError("unknown operator: !%s", right.Type())
	}

	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		// unreachable
		return NULL
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return object.NewError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evalBooleanInfixExpression(op, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(op, left, right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), op, left.Type())
	}
}

func evalBooleanInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch op {
	case "==":
		return toBooleanObject(leftVal == rightVal)
	case "!=":
		return toBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return object.NewError("division by zero")
		}
		if leftVal == math.MinInt64 && rightVal == -1 {
			return object.NewError("integer overflow: %d / %d", leftVal, rightVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return toBooleanObject(leftVal < rightVal)
	case ">":
		return toBooleanObject(leftVal > rightVal)
	case "==":
		return toBooleanObject(leftVal == rightVal)
	case "!=":
		return toBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch op {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return toBooleanObject(leftVal == rightVal)
	case "!=":
		return toBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}

	if cond.Type() != object.BOOLEAN {
		return NULL
	}

	if cond == TRUE {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalBlockStatement(block *ast.Block, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil && (result.Type() == object.RETURN_VALUE || result.Type() == object.ERROR) {
			return result
		}
	}

	return result
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		array := left.(*object.Array)
		idx := index.(*object.Integer).Value

		if idx < 0 || idx >= int64(len(array.Elements)) {
			return object.NewError("index out of range")
		}

		return array.Elements[idx]
	default:
		return object.NewError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if callDepth >= maxCallDepth {
			return object.NewError("maximum call depth exceeded: %d", maxCallDepth)
		}
		callDepth++
		defer func() { callDepth-- }()

		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if callDepth >= maxCallDepth {
			return object.NewError("maximum call depth exceeded: %d", maxCallDepth)
		}
		callDepth++
		defer func() { callDepth-- }()

		return fn.Run(args...)
	default:
		return object.NewError("not a function: %s", fn.Type())
	}

}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for index, param := range fn.Parameters {
		env.Set(param.Value, args[index])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func toBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
