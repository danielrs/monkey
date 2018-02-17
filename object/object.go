package object

import (
	"bytes"
	"fmt"

	"github.com/danielrs/monkey/ast"
)

type ObjectType string

const (
	NIL_OBJ          = "NIL"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION_OBJ"
	ERROR_OBJ        = "ERROR_OBJ"
)

// Environment.

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: outer}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Individual objects.

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("%q", s.Value) }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters ast.ParameterList
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(f.Parameters.String())
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
