package parser

import (
	"fmt"
	"testing"

	"github.com/danielrs/monkey/ast"
	"github.com/danielrs/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		value      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{`let foo = "bar";`, "foo", "bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.LetStatement")
			t.FailNow()
		}

		if !testLetStatement(t, stmt, tt.identifier) ||
			!testLiteralExpression(t, stmt.Value, tt.value) {
			t.FailNow()
		}
	}

}

func testLetStatement(t *testing.T, s ast.Statement, identifier string) bool {
	letstmt, ok := s.(*ast.LetStatement)
	if !ok {
		castError(t, s, "*ast.LetStatement")
		return false
	}

	if letstmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral got %q, want 'let'", s.TokenLiteral())
		return false
	}

	if letstmt.Identifier.TokenLiteral() != identifier {
		t.Errorf("letstmt.Identifier=%s, want %q",
			letstmt.Identifier.TokenLiteral(),
			identifier)
		return false
	}

	if letstmt.Identifier.Value != identifier {
		t.Errorf("letstmt.Identifier.Value=%s, want %q",
			letstmt.Identifier.Value,
			identifier)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"return 5", 5},
		{"return false", false},
		{`return "foobar"`, "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ReturnStatement")
			t.FailNow()
		}

		if !testReturnStatement(t, stmt) ||
			!testLiteralExpression(t, stmt.Value, tt.value) {
			t.FailNow()
		}
	}
}

func testReturnStatement(t *testing.T, s ast.Statement) bool {
	retstmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		castError(t, s, "*ast.LetStatement")
		return false
	}

	if retstmt.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral got %q, want 'return'", s.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		castError(t, program.Statements[0], "*ast.ExpressionStatement")
	}
	testIdentifierExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"5;", 5},
		{"05;", 5},
		{"005", 5},
		{"010", 10},
		{"0010", 10},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
		}

		testIntegerLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"false;", false},
		{"true;", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
		}
		testBooleanLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar";`, "foobar"},
		{`"foo bar";`, "foo bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
		}
		testStringLiteralExpression(t, stmt.Expression, tt.expected)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		if !testPrefixExpression(t, program.Statements[0], tt.operator, tt.value) {
			t.FailNow()
		}
	}
}

func testPrefixExpression(t *testing.T, s ast.Statement, operator string, value interface{}) bool {
	stmt, ok := s.(*ast.ExpressionStatement)
	if !ok {
		castError(t, s, "*ast.ExpressionStatement")
		return false
	}

	expr, ok := stmt.Expression.(*ast.PrefixExpression)
	if !ok {
		castError(t, stmt.Expression, "*ast.PrefixExpression")
		return false
	}
	if expr.Operator != operator {
		t.Errorf("expr.Operator is %q, want %q", expr.Operator, operator)
		return false
	}
	if !testLiteralExpression(t, expr.Right, value) {
		return false
	}

	return true
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 % 5", 5, "%", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{`"foo" + "bar"`, "foo", "+", "bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		if !testInfixExpression(t, program.Statements[0], tt.left, tt.operator, tt.right) {
			t.FailNow()
		}
	}
}

func testInfixExpression(t *testing.T, s ast.Statement, left interface{}, operator string, right interface{}) bool {
	stmt, ok := s.(*ast.ExpressionStatement)
	if !ok {
		castError(t, stmt, "*ast.ExpressionStatement")
		return false
	}

	expr, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		castError(t, stmt.Expression, "*ast.InfixExpression")
		return false
	}

	if !testLiteralExpression(t, expr.Left, left) {
		return false
	}
	if expr.Operator != operator {
		t.Errorf("expr.Operator is %q, want %q", expr.Operator, operator)
	}
	if !testLiteralExpression(t, expr.Right, right) {
		return false
	}

	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// Grouped expressions.
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		// Index.
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("got %q, want %q", actual, tt.expected)
		}
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input       string
		condition   string
		consequence string
		alternative string
	}{
		{"if (x < y) { x }", "(x < y)", "x", ""},
		{"if (x < y) { x } else { x+y }", "(x < y)", "x", "(x + y)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
			t.FailNow()
		}

		expr, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			castError(t, stmt.Expression, "*ast.IfExpression")
			t.FailNow()
		}

		if expr.Condition.String() != tt.condition {
			t.Errorf("expr.Condition is %s, want %s",
				expr.Condition.String(), tt.condition)
		}
		if expr.Consequence.String() != tt.consequence {
			t.Errorf("expr.Consequence is %s, want %s",
				expr.Consequence.String(), tt.consequence)
		}
		if len(tt.alternative) > 0 && expr.Alternative.String() != tt.alternative {
			t.Errorf("expr.Alternative is %s, want %s",
				expr.Alternative.String(), tt.alternative)
		}
	}
}

func TestFunctionLiteral(t *testing.T) {
	tests := []struct {
		input      string
		parameters string
		body       string
	}{
		{"fn() {}", "", ""},
		{"fn(a) { a+1 }", "a", "(a + 1)"},
		{"fn(a,b) {a+b}; ", "a, b", "(a + b)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
			t.FailNow()
		}

		expr, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			castError(t, stmt.Expression, "*ast.FunctionLiteral")
			t.FailNow()
		}

		if expr.Parameters.String() != tt.parameters {
			t.Errorf("expr.Parameters.String() is %s, want %s",
				expr.Parameters.String(), tt.parameters)
			t.FailNow()
		}
		if expr.Body.String() != tt.body {
			t.Errorf("expr.Body.String() is %s, want %s",
				expr.Body.String(), tt.body)
			t.FailNow()
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		castError(t, stmt.Expression, "*ast.ArrayLiteral")
		t.FailNow()
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) got %d, want %d",
			len(array.Elements), 3)
	}

	testIntegerLiteralExpression(t, array.Elements[0], 1)
	testInfixExpression(t,
		&ast.ExpressionStatement{Expression: array.Elements[1]}, 2, "*", 2)
	testInfixExpression(t,
		&ast.ExpressionStatement{Expression: array.Elements[2]}, 3, "+", 3)
}

func TestEmptyArrayLiteral(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		castError(t, stmt.Expression, "*ast.ArrayLiteral")
		t.FailNow()
	}

	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) got %d, want %d",
			len(array.Elements), 0)
	}
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	expr, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		castError(t, stmt.Expression, "*ast.IndexExpression")
		t.FailNow()
	}

	if !testIdentifierExpression(t, expr.Left, "myArray") {
		t.FailNow()
	}
	if !testInfixExpression(t, &ast.ExpressionStatement{Expression: expr.Index}, 1, "+", 1) {
		t.FailNow()
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		input     string
		function  string
		arguments []string
	}{
		{"init()", "init", []string{}},
		{"add(x, y)", "add", []string{"x", "y"}},
		{"sub(x+10, y*5)", "sub", []string{"(x + 10)", "(y * 5)"}},
		{"fn(a) { a * a }(5)", "fn(a) (a * a)", []string{"5"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement, got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
			t.FailNow()
		}

		expr, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			castError(t, stmt.Expression, "*ast.CallExpression")
			t.FailNow()
		}

		if expr.Function.String() != tt.function {
			t.Errorf("expr.Function.String() is %s, want %s",
				expr.Function.String(), tt.function)
			t.FailNow()
		}
		if len(expr.Arguments) != len(tt.arguments) {
			t.Errorf("got %d arguments, want %d",
				len(expr.Arguments), len(tt.arguments))
			t.FailNow()
		}

		for i, a := range expr.Arguments {
			if a.String() != tt.arguments[i] {
				t.Errorf("Argument #%d is %s, want %s",
					i, a.String(), tt.arguments[i])
				t.FailNow()
			}
		}
	}
}

// Helper functions for testing.

func testIdentifierExpression(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		castError(t, expr, "*ast.Identifier")
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is %q, want %q", ident.TokenLiteral(), value)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value is %q, want %q", ident.Value, value)
		return false
	}
	return true
}

func testIntegerLiteralExpression(t *testing.T, expr ast.Expression, value int64) bool {
	integer, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		castError(t, expr, "*ast.IntegerLiteral")
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value is %d, want %d", integer.Value, value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() is %s, want %d",
			integer.TokenLiteral(), value)
		return false
	}
	return true
}

func testBooleanLiteralExpression(t *testing.T, expr ast.Expression, value bool) bool {
	boolean, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		castError(t, expr, "*ast.BooleanLiteral")
		return false
	}
	if boolean.Value != value {
		t.Errorf("boolean.Value is %t, want %t", boolean.Value, value)
		return false
	}
	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() is %t, want %t",
			boolean.TokenLiteral(), value)
		return false
	}
	return true
}

func testStringLiteralExpression(t *testing.T, expr ast.Expression, value string) bool {
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		castError(t, expr, "*ast.StringLiteral")
		return false
	}
	if str.Value != value {
		t.Errorf("str.Value is %q, want %q", str.Value, value)
		return false
	}
	if str.TokenLiteral() != value {
		t.Errorf("str.TokenLiteral() is %q, want %q", str.TokenLiteral(), value)
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteralExpression(t, expr, int64(v))
	case int64:
		return testIntegerLiteralExpression(t, expr, v)
	case bool:
		return testBooleanLiteralExpression(t, expr, v)
	case string:
		return testStringLiteralExpression(t, expr, v)
	}

	t.Errorf("type %T of expr not handled, want %T", expr, expected)
	return false
}

// Helper functions for errors.

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}

	t.FailNow()
}

func castError(t *testing.T, got interface{}, want string) {
	t.Errorf("cast error: got %T, want %s", got, want)
}
