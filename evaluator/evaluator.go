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

var builtins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			vararg := make([]interface{}, 0, len(args))
			for _, a := range args {
				vararg = append(vararg, a.Inspect())
			}
			fmt.Println(vararg...)
			return NULL
		},
	},

	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. want %d, got %d",
					1, len(args))
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Elements))}
			}

			return newError("argument to `len` not supported, got %s",
				args[0].Type())
		},
	},

	"head": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. want %d, got %d",
					1, len(args))
			}

			switch arr := args[0].(type) {
			case *object.Array:
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return NULL
			}

			return newError("argument to `head` not supported, got %s",
				args[0].Type())
		},
	},

	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. want %d, got %d",
					1, len(args))
			}

			switch arr := args[0].(type) {
			case *object.Array:
				if len(arr.Elements) > 0 {
					return arr.Elements[len(arr.Elements)-1]
				}
				return NULL
			}

			return newError("argument to `last` not supported, got %s",
				args[0].Type())
		},
	},

	"tail": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. want %d, got %d",
					1, len(args))
			}

			switch arr := args[0].(type) {
			case *object.Array:
				if len(arr.Elements) > 0 {
					newElems := make([]object.Object, len(arr.Elements)-1)
					copy(newElems, arr.Elements[1:])
					return &object.Array{newElems}
				}
				return &object.Array{[]object.Object{}}
			}

			return newError("argument to `tail` not supported, got %s",
				args[0].Type())
		},
	},

	"init": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. want %d, got %d",
					1, len(args))
			}

			switch arr := args[0].(type) {
			case *object.Array:
				if len(arr.Elements) > 0 {
					newElems := make([]object.Object, len(arr.Elements)-1)
					copy(newElems, arr.Elements[:len(arr.Elements)-1])
					return &object.Array{newElems}
				}
				return &object.Array{[]object.Object{}}
			}

			return newError("argument to `init` not supported, got %s",
				args[0].Type())
		},
	},

	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. want %d, got %d",
					2, len(args))
			}

			switch arr := args[0].(type) {
			case *object.Array:
				length := len(arr.Elements)
				newElems := make([]object.Object, length, length+1)
				copy(newElems, arr.Elements)
				newElems = append(newElems, args[1])
				return &object.Array{newElems}
			}

			return newError("argument to `push` not supported, got %s",
				args[0].Type())
		},
	},
}

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

	case *ast.ArrayLiteral:
		elems := evalExpressions(env, node.Elements)
		if len(elems) >= 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{elems}

	case *ast.IndexExpression:
		return try(Eval(env, node.Left), func(l object.Object) object.Object {
			return try(Eval(env, node.Index), func(i object.Object) object.Object {
				return evalIndexExpression(l, i)
			})
		})

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

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arr := left.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(arr.Elements) - 1)
		if idx < 0 || idx > max {
			return NULL
		}
		return arr.Elements[idx]
	}

	return newError("index operator not supported: %s", left.Type())
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
	case "%":
		return &object.Integer{l.Value % r.Value}

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
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
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
	switch fn := fn.(type) {
	case *object.Function:
		// Extends environment.
		newEnv := object.NewEnclosedEnvironment(fn.Env)
		if len(fn.Parameters) != len(args) {
			return newError("argument mismatch: got %d, want %d",
				len(args), len(fn.Parameters))
		}
		for paramIdx, param := range fn.Parameters {
			newEnv.Set(param.Value, args[paramIdx])
		}
		// Evaluates it.
		evaluated := Eval(newEnv, fn.Body)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)
	}

	return newError("not a function: %s", fn.Type())
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
