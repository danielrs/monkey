package evaluator

import (
	"fmt"

	"github.com/danielrs/monkey/ast"
	"github.com/danielrs/monkey/object"
)

var (
	NULL  = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(env *object.Environment, node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements.
	case *ast.Program:
		return evalProgram(env, node)

	case *ast.ExpressionStatement:
		return Eval(env, node.Expression)

	case *ast.BlockStatement:
		return evalBlockStatement(env, node)

	case *ast.ReturnStatement:
		return try(Eval(env, node.Value), func(val object.Object) object.Object {
			return &object.ReturnValue{val}
		})

	case *ast.LetStatement:
		return try(Eval(env, node.Value), func(val object.Object) object.Object {
			env.Set(node.Identifier.Value, val)
			return &object.Nil{}
		})

	// Expressions.
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.IntegerLiteral:
		return &object.Integer{node.Value}

	case *ast.StringLiteral:
		return &object.String{node.Value}

	case *ast.Identifier:
		return evalIdentifier(env, node)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.PrefixExpression:
		return try(Eval(env, node.Right), func(right object.Object) object.Object {
			return evalPrefixExpression(node.Operator, right)
		})

	case *ast.InfixExpression:
		return try(Eval(env, node.Left), func(l object.Object) object.Object {
			return try(Eval(env, node.Right), func(r object.Object) object.Object {
				return evalInfixExpression(node.Operator, l, r)
			})
		})

	case *ast.IfExpression:
		return evalIfExpression(env, node)

	case *ast.CallExpression:
		return try(Eval(env, node.Function), func(f object.Object) object.Object {
			args := evalExpressions(env, node.Arguments)
			if len(args) >= 1 && isError(args[0]) {
				return args[0]
			}
			return applyFunction(f, args)
		})
	}

	return nil
}

func evalProgram(env *object.Environment, program *ast.Program) object.Object {
	var result object.Object
	for _, s := range program.Statements {
		result = Eval(env, s)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(env *object.Environment, block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, s := range block.Statements {
		result = Eval(env, s)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ ||
				result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(obj object.Object) object.Object {
	switch number := obj.(type) {
	case *object.Integer:
		return &object.Integer{-number.Value}
	default:
		return newError("unknown operator: -%s", obj.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case operator == "==":
		return &object.Boolean{left == right}

	case operator == "!=":
		return &object.Boolean{left != right}
	}

	return newError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	l, lok := left.(*object.Integer)
	r, rok := right.(*object.Integer)

	if !lok || !rok {
		return NULL
	}

	switch operator {
	// Returns number.
	case "+":
		return &object.Integer{l.Value + r.Value}
	case "-":
		return &object.Integer{l.Value - r.Value}
	case "*":
		return &object.Integer{l.Value * r.Value}
	case "/":
		return &object.Integer{l.Value / r.Value}

	// Returns boolean.
	case "<":
		return nativeBooleanToObject(l.Value < r.Value)
	case ">":
		return nativeBooleanToObject(l.Value > r.Value)
	case "==":
		return nativeBooleanToObject(l.Value == r.Value)
	case "!=":
		return nativeBooleanToObject(l.Value != r.Value)
	}

	return newError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l, lok := left.(*object.String)
	r, rok := right.(*object.String)

	if !lok || !rok {
		return NULL
	}

	switch operator {
	// Returns concatenated string.
	case "+":
		return &object.String{fmt.Sprintf("%s%s", l.Value, r.Value)}
	}

	return newError("unknown operator: %s %s %s",
		left.Type(), operator, right.Type())
}

func evalIfExpression(env *object.Environment, expr *ast.IfExpression) object.Object {
	return try(Eval(env, expr.Condition), func(pred object.Object) object.Object {
		if isTruthy(pred) {
			return Eval(env, expr.Consequence)
		} else if expr.Alternative != nil {
			return Eval(env, expr.Alternative)
		}
		return NULL
	})
}

func evalIdentifier(env *object.Environment, node *ast.Identifier) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

func evalExpressions(env *object.Environment, exprs []ast.Expression) []object.Object {
	var result []object.Object

	for _, e := range exprs {
		evaluated := Eval(env, e)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	// Extends environment.
	newEnv := object.NewEnclosedEnvironment(function.Env)
	if len(function.Parameters) != len(args) {
		return newError("argument mismatch: got %d, want %d",
			len(args), len(function.Parameters))
	}
	for paramIdx, param := range function.Parameters {
		newEnv.Set(param.Value, args[paramIdx])
	}

	// Evaluates it.
	evaluated := Eval(newEnv, function.Body)
	return unwrapReturnValue(evaluated)
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// Helper functions.

// Checks the given object, if it's an error, returns it;
// otherwise, calls the given do function passing obj and
// returns its value.
func try(obj object.Object, do func(object.Object) object.Object) object.Object {
	if isError(obj) {
		return obj
	}
	return do(obj)
}

func nativeBooleanToObject(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}
