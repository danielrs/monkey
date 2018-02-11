package evaluator

import (
	"github.com/danielrs/monkey/ast"
	"github.com/danielrs/monkey/object"
)

var (
	NULL  = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements.
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.BlockStatement:
		return evalStatements(node.Statements)

	// Expressions.
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.IntegerLiteral:
		return &object.Integer{node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, s := range statements {
		result = Eval(s)
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
		return NULL
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
		return NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return &object.Boolean{left == right}
	case operator == "!=":
		return &object.Boolean{left != right}
	}

	return NULL
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

	return NULL
}

func evalIfExpression(expr *ast.IfExpression) object.Object {
	predicate := Eval(expr.Condition)
	if isTruthy(predicate) {
		return Eval(expr.Consequence)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative)
	}
	return NULL
}

// Helper functions.

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
