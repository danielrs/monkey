package evaluator

import (
	"testing"

	"github.com/danielrs/monkey/lexer"
	"github.com/danielrs/monkey/object"
	"github.com/danielrs/monkey/parser"
)

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testBooleanObject(t, obj, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"--5", 5},
		{"---10", -10},
		{"5+10", 15},
		{"5-10", -5},
		{"5*10", 50},
		{"5/10", 0},
		{"10/5", 2},
		{"2*5 + 20*1", 30},
		{"-50 + 30 - 10 + 30", 0},
		{"3 * (3 * 3) + 10", 37},
		{"7 * 5 + 10 / 2 + 1000/500", 42},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foo"`, "foo"},
		{`"foo" + "bar"`, "foobar"},
		{`"foo" + "bar" + "baz"`, "foobarbaz"},
		{`"foo" + ("bar" + "baz")`, "foobarbaz"},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testStringObject(t, obj, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false},
		{"!!5", true},
		{"5 < 10", true},
		{"5 > 10", false},
		{"5 == 10", false},
		{"5 == 5", true},
		{"5 + 1 == 6", true},
		{"6*7 == 84/2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		if !testBooleanObject(t, obj, tt.expected) {
			t.Log(tt.input)
		}
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { true } else { false }", true},
		{"if (5 > 10) { 10 } else { 5 }", 5},
		{"if (5 == 5) { 5 }", 5},
		{"if (5 == 10) { 5 }", nil},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"9; return true; false;", true},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = 5; let c = a + b + 5; c", 15},
		{"let a = 5;", nil},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let id = fn(x) { x; }; id(5);", 5},
		{"let id = fn(x) { return x; }; id(5);", 5},
		{"let double = fn(x) { x * 2 }; double(5);", 10},
		{"let add = fn(x, y) { x + y }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y }; add(5+5, add(5, 5));", 20},
		{"fn(x){ x; }(5)", 5},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
	let makeAdder = fn(x) {
		fn(y) { x + y; };
	}

	let addTwo = makeAdder(2);
	addTwo(3);
	`

	obj := testEval(input)
	testIntegerObject(t, obj, 5)
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 10;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { if (10 > 1) { true + false; } return 1; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"fn(x) { x + y; }(1)",
			"identifier not found: y",
		},
		{
			`"foo" - "bar"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testErrorObject(t, obj, tt.expected)
	}
}

// Helper functions for testing.

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(env, program)
}

func testNilObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Nil)
	if !ok {
		castError(t, obj, "*object.Nil")
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolobj, ok := obj.(*object.Boolean)
	if !ok {
		castError(t, obj, "*object.Boolean")
		return false
	}
	if boolobj.Value != expected {
		expectedError(t, "boolobj.Value", boolobj.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intobj, ok := obj.(*object.Integer)
	if !ok {
		castError(t, obj, "*object.Integer")
		return false
	}
	if intobj.Value != expected {
		expectedError(t, "intobj.Value", intobj.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	strobj, ok := obj.(*object.String)
	if !ok {
		castError(t, obj, "*object.String")
		return false
	}
	if strobj.Value != expected {
		expectedError(t, "strobj.Value", strobj.Value, expected)
		return false
	}
	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	errobj, ok := obj.(*object.Error)
	if !ok {
		castError(t, obj, "*object.Error")
		return false
	}
	if errobj.Message != expected {
		expectedError(t, "errobj.Message", errobj.Message, expected)
		return false
	}
	return true
}

func testObject(t *testing.T, obj object.Object, expected interface{}) bool {
	switch v := expected.(type) {
	case nil:
		return testNilObject(t, obj)
	case bool:
		return testBooleanObject(t, obj, v)
	case int64:
		return testIntegerObject(t, obj, v)
	case int:
		return testIntegerObject(t, obj, int64(v))
	case string:
		return testStringObject(t, obj, v)
	}

	return false
}

// Helper functions for errors.

func castError(t *testing.T, got interface{}, want string) {
	t.Errorf("cast error: got %T, want %s", got, want)
}

func expectedError(t *testing.T, name string, got interface{}, want interface{}) {
	t.Errorf("%s is %v, want %v", name, got, want)
}
