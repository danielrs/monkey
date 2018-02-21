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

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Integers.
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
		{"5%10", 5},
		{"10%10", 0},
		{"3%2", 1},
		{"2*5 + 20*1", 30},
		{"-50 + 30 - 10 + 30", 0},
		{"3 * (3 * 3) + 10", 37},
		{"7 * 5 + 10 / 2 + 1000/500", 42},
		// Strings.
		{`"foo"`, "foo"},
		{`"foo" + "bar"`, "foobar"},
		{`"foo" + "bar" + "baz"`, "foobarbaz"},
		{`"foo" + ("bar" + "baz")`, "foobarbaz"},
		// Logical.
		{"false && false", false},
		{"false && true", false},
		{"true && false", false},
		{"true && true", true},
		{"false || false", false},
		{"false || true", true},
		{"true || false", true},
		{"true || true", true},
		{"true && 5", 5},
		{"false && 5", false},
		{"false || 5", 5},
		{"true || 5", true},
		{`(false || false || false || true) && (true && true && "foo")`, "foo"},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestBooleanExpression(t *testing.T) {
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

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{"[1, 2, 3]", []interface{}{1, 2, 3}},
		{"[]", []interface{}{}},
		{"[true, false]", []interface{}{true, false}},
		{`["foo", "bar", "baz"]`, []interface{}{"foo", "bar", "baz"}},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testArrayObject(t, obj, tt.expected)
	}
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1]", 3},
		{"let arr = [1, 2, 3]; arr[0]", 1},
		{"let arr = [1, 2, 3]; arr[1]", 2},
		{"let arr = [1, 2, 3]; arr[2]", 3},
		{"let i = 1; let arr = [1, 2, 3]; arr[0]; arr[i];", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestHashLiteralExpression(t *testing.T) {
	input := `
	let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	expected := map[object.HashKey]interface{}{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	obj := testEval(input)
	hash, ok := obj.(*object.Hash)
	if !ok {
		castError(t, obj, "*object.Hash")
		t.FailNow()
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("len(hash.Pairs) got %d, want %d",
			len(hash.Pairs), len(expected))
		t.FailNow()
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in map")
			break
		}
		testObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 1}["foo"]`, 1},
		{`{"foo": 1}["bar"]`, nil},
		{`let key = "foo"; {"foo": "bar"}[key]`, "bar"},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: true}[true]`, true},
		{`{false: false}[false]`, false},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		// Arrays.
		{`len([1, 2, 3])`, 3},
		{`len([1, 2])`, 2},
		{`len([1])`, 1},
		{`len([])`, 0},
		{`head([1, 2])`, 1},
		{`head([2])`, 2},
		{`head([])`, nil},
		{`last([1, 2])`, 2},
		{`last([1])`, 1},
		{`last([])`, nil},
		{`tail([1, 2, 3])`, []interface{}{2, 3}},
		{`tail([2, 3])`, []interface{}{3}},
		{`tail([3])`, []interface{}{}},
		{`tail([])`, []interface{}{}},
		{`init([1, 2, 3])`, []interface{}{1, 2}},
		{`init([2, 3])`, []interface{}{2}},
		{`init([3])`, []interface{}{}},
		{`init([])`, []interface{}{}},
		{`push([], 1)`, []interface{}{1}},
		{`push([1], 2)`, []interface{}{1, 2}},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testObject(t, obj, tt.expected)
	}
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
		// Arrays and Hashes.
		{
			`{"foo": "bar"}[fn(x) { x }]`,
			"unusable as hash key: FUNCTION_OBJ",
		},
		// Builtin.
		{
			`len(1)`,
			"argument to `len` not supported, got INTEGER",
		},
		{
			`len("one", "two")`,
			"wrong number of arguments. want 1, got 2",
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

func testArrayObject(t *testing.T, obj object.Object, expected []interface{}) bool {
	arr, ok := obj.(*object.Array)
	if !ok {
		castError(t, obj, "*object.Array")
		return false
	}
	if len(arr.Elements) != len(expected) {
		t.Errorf("len(arr.Elements) got %d, want %d",
			len(arr.Elements), len(expected))
		return false
	}
	for i := range arr.Elements {
		if !testObject(t, arr.Elements[i], expected[i]) {
			return false
		}
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
	case []interface{}:
		return testArrayObject(t, obj, v)
	}

	t.Errorf("Object %T not handled", obj)
	return false
}

// Helper functions for errors.

func castError(t *testing.T, got interface{}, want string) {
	t.Errorf("cast error: got %T, want %s", got, want)
}

func expectedError(t *testing.T, name string, got interface{}, want interface{}) {
	t.Errorf("%s is %v, want %v", name, got, want)
}
