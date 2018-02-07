package parser

import (
	"fmt"
	"testing"

	"github.com/danielrs/monkey/ast"
	"github.com/danielrs/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let	y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
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
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		testReturnStatement(t, stmt)
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
	input := "5;"

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
	testIntegerLiteralExpression(t, stmt.Expression, 5)
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

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			castError(t, program.Statements[0], "*ast.ExpressionStatement")
		}
		testBooleanLiteralExpression(t, stmt.Expression, tt.expected)
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
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

// Helper functions for parsing.

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
		t.Errorf("integer.TokenLiteral() is %d, want %d",
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

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteralExpression(t, expr, int64(v))
	case int64:
		return testIntegerLiteralExpression(t, expr, v)
	case string:
		return testIdentifierExpression(t, expr, v)
	case bool:
		return testBooleanLiteralExpression(t, expr, v)
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
