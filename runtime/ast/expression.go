/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ast

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/turbolent/prettier"
)

const NilConstant = "nil"

type Expression interface {
	Element
	fmt.Stringer
	IfStatementTest
	isExpression()
	AcceptExp(ExpressionVisitor) Repr
	Doc() prettier.Doc
}

// BoolExpression

type BoolExpression struct {
	Value bool
	Range
}

var _ Expression = &BoolExpression{}

func (*BoolExpression) isExpression() {}

func (*BoolExpression) isIfStatementTest() {}

func (e *BoolExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*BoolExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *BoolExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitBoolExpression(e)
}

func (e *BoolExpression) String() string {
	if e.Value {
		return "true"
	}
	return "false"
}

var boolExpressionTrueDoc prettier.Doc = prettier.Text("true")
var boolExpressionFalseDoc prettier.Doc = prettier.Text("false")

func (e *BoolExpression) Doc() prettier.Doc {
	if e.Value {
		return boolExpressionTrueDoc
	} else {
		return boolExpressionFalseDoc
	}
}

func (e *BoolExpression) MarshalJSON() ([]byte, error) {
	type Alias BoolExpression
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "BoolExpression",
		Alias: (*Alias)(e),
	})
}

// NilExpression

type NilExpression struct {
	Pos Position `json:"-"`
}

var _ Expression = &NilExpression{}

func (*NilExpression) isExpression() {}

func (*NilExpression) isIfStatementTest() {}

func (e *NilExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*NilExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *NilExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitNilExpression(e)
}

func (e *NilExpression) String() string {
	return NilConstant
}

var nilExpressionDoc prettier.Doc = prettier.Text("nil")

func (*NilExpression) Doc() prettier.Doc {
	return nilExpressionDoc
}

func (e *NilExpression) StartPosition() Position {
	return e.Pos
}

func (e *NilExpression) EndPosition() Position {
	return e.Pos.Shifted(len(NilConstant) - 1)
}

func (e *NilExpression) MarshalJSON() ([]byte, error) {
	type Alias NilExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "NilExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// StringExpression

type StringExpression struct {
	Value string
	Range
}

var _ Expression = &StringExpression{}

func (*StringExpression) isExpression() {}

func (*StringExpression) isIfStatementTest() {}

func (e *StringExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*StringExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *StringExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitStringExpression(e)
}

func (e *StringExpression) String() string {
	return QuoteString(e.Value)
}

func (e *StringExpression) Doc() prettier.Doc {
	return prettier.Text(QuoteString(e.Value))
}

func (e *StringExpression) MarshalJSON() ([]byte, error) {
	type Alias StringExpression
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "StringExpression",
		Alias: (*Alias)(e),
	})
}

// IntegerExpression

type IntegerExpression struct {
	PositiveLiteral string
	Value           *big.Int `json:"-"`
	Base            int
	Range
}

var _ Expression = &IntegerExpression{}

func (*IntegerExpression) isExpression() {}

func (*IntegerExpression) isIfStatementTest() {}

func (e *IntegerExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*IntegerExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *IntegerExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitIntegerExpression(e)
}

func (e *IntegerExpression) String() string {
	literal := e.PositiveLiteral
	if e.Value.Sign() < 0 {
		literal = "-" + literal
	}
	return literal
}

func (e *IntegerExpression) Doc() prettier.Doc {
	literal := e.PositiveLiteral
	if e.Value.Sign() < 0 {
		literal = "-" + literal
	}
	return prettier.Text(literal)
}

func (e *IntegerExpression) MarshalJSON() ([]byte, error) {
	type Alias IntegerExpression
	return json.Marshal(&struct {
		Type  string
		Value string
		*Alias
	}{
		Type:  "IntegerExpression",
		Value: e.Value.String(),
		Alias: (*Alias)(e),
	})
}

// FixedPointExpression

type FixedPointExpression struct {
	PositiveLiteral string
	Negative        bool
	UnsignedInteger *big.Int `json:"-"`
	Fractional      *big.Int `json:"-"`
	Scale           uint
	Range
}

var _ Expression = &FixedPointExpression{}

func (*FixedPointExpression) isExpression() {}

func (*FixedPointExpression) isIfStatementTest() {}

func (e *FixedPointExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*FixedPointExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *FixedPointExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitFixedPointExpression(e)
}

func (e *FixedPointExpression) String() string {
	literal := e.PositiveLiteral
	if literal != "" {
		if e.Negative {
			literal = "-" + literal
		}
		return literal
	}

	var builder strings.Builder
	if e.Negative {
		builder.WriteRune('-')
	}
	builder.WriteString(e.UnsignedInteger.String())
	builder.WriteRune('.')
	fractional := e.Fractional.String()
	for i := uint(0); i < (e.Scale - uint(len(fractional))); i++ {
		builder.WriteRune('0')
	}
	builder.WriteString(fractional)
	return builder.String()
}

func (e *FixedPointExpression) Doc() prettier.Doc {
	literal := e.PositiveLiteral
	if e.Negative {
		literal = "-" + literal
	}
	return prettier.Text(literal)
}

func (e *FixedPointExpression) MarshalJSON() ([]byte, error) {
	type Alias FixedPointExpression
	return json.Marshal(&struct {
		Type            string
		UnsignedInteger string
		Fractional      string
		*Alias
	}{
		Type:            "FixedPointExpression",
		UnsignedInteger: e.UnsignedInteger.String(),
		Fractional:      e.Fractional.String(),
		Alias:           (*Alias)(e),
	})
}

// ArrayExpression

type ArrayExpression struct {
	Values []Expression
	Range
}

var _ Expression = &ArrayExpression{}

func (*ArrayExpression) isExpression() {}

func (*ArrayExpression) isIfStatementTest() {}

func (e *ArrayExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *ArrayExpression) Walk(walkChild func(Element)) {
	walkExpressions(walkChild, e.Values)
}

func (e *ArrayExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitArrayExpression(e)
}

func (e *ArrayExpression) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, value := range e.Values {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(value.String())
	}
	builder.WriteString("]")
	return builder.String()
}

var arrayExpressionSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Text(","),
	prettier.Line{},
}

func (e *ArrayExpression) Doc() prettier.Doc {
	if len(e.Values) == 0 {
		return prettier.Text("[]")
	}

	elementDocs := make([]prettier.Doc, len(e.Values))
	for i, value := range e.Values {
		elementDocs[i] = value.Doc()
	}
	return prettier.WrapBrackets(
		prettier.Join(arrayExpressionSeparatorDoc, elementDocs...),
		prettier.SoftLine{},
	)
}

func (e *ArrayExpression) MarshalJSON() ([]byte, error) {
	type Alias ArrayExpression
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "ArrayExpression",
		Alias: (*Alias)(e),
	})
}

// DictionaryExpression

type DictionaryExpression struct {
	Entries []DictionaryEntry
	Range
}

var _ Expression = &DictionaryExpression{}

func (*DictionaryExpression) isExpression() {}

func (*DictionaryExpression) isIfStatementTest() {}

func (e *DictionaryExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *DictionaryExpression) Walk(walkChild func(Element)) {
	for _, entry := range e.Entries {
		walkChild(entry.Key)
		walkChild(entry.Value)
	}
}

func (e *DictionaryExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitDictionaryExpression(e)
}

func (e *DictionaryExpression) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	for i, entry := range e.Entries {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(entry.Key.String())
		builder.WriteString(": ")
		builder.WriteString(entry.Value.String())
	}
	builder.WriteString("}")
	return builder.String()
}

var dictionaryExpressionSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Text(","),
	prettier.Line{},
}

func (e *DictionaryExpression) Doc() prettier.Doc {
	if len(e.Entries) == 0 {
		return prettier.Text("{}")
	}

	entryDocs := make([]prettier.Doc, len(e.Entries))
	for i, entry := range e.Entries {
		entryDocs[i] = entry.Doc()
	}

	return prettier.WrapBraces(
		prettier.Join(dictionaryExpressionSeparatorDoc, entryDocs...),
		prettier.SoftLine{},
	)
}

func (e *DictionaryExpression) MarshalJSON() ([]byte, error) {
	type Alias DictionaryExpression
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "DictionaryExpression",
		Alias: (*Alias)(e),
	})
}

type DictionaryEntry struct {
	Key   Expression
	Value Expression
}

func (e DictionaryEntry) MarshalJSON() ([]byte, error) {
	type Alias DictionaryEntry
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "DictionaryEntry",
		Alias: (*Alias)(&e),
	})
}

var dictionaryKeyValueSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Text(":"),
	prettier.Line{},
}

func (e DictionaryEntry) Doc() prettier.Doc {
	keyDoc := e.Key.Doc()
	valueDoc := e.Value.Doc()

	return prettier.Group{
		Doc: prettier.Concat{
			keyDoc,
			dictionaryKeyValueSeparatorDoc,
			valueDoc,
		},
	}
}

// IdentifierExpression

type IdentifierExpression struct {
	Identifier Identifier
}

var _ Expression = &IdentifierExpression{}

func (*IdentifierExpression) isExpression() {}

func (*IdentifierExpression) isIfStatementTest() {}

func (e *IdentifierExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*IdentifierExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *IdentifierExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitIdentifierExpression(e)
}

func (e *IdentifierExpression) String() string {
	return e.Identifier.Identifier
}

func (e *IdentifierExpression) Doc() prettier.Doc {
	return prettier.Text(e.Identifier.Identifier)
}

func (e *IdentifierExpression) MarshalJSON() ([]byte, error) {
	type Alias IdentifierExpression
	return json.Marshal(&struct {
		Type string
		*Alias
		Range
	}{
		Type:  "IdentifierExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

func (e *IdentifierExpression) StartPosition() Position {
	return e.Identifier.StartPosition()
}

func (e *IdentifierExpression) EndPosition() Position {
	return e.Identifier.EndPosition()
}

// Arguments

type Arguments []*Argument

func (args Arguments) String() string {
	var builder strings.Builder
	builder.WriteRune('(')
	for i, argument := range args {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(argument.String())
	}
	builder.WriteRune(')')
	return builder.String()
}

// InvocationExpression

type InvocationExpression struct {
	InvokedExpression Expression
	TypeArguments     []*TypeAnnotation
	Arguments         Arguments
	ArgumentsStartPos Position
	EndPos            Position `json:"-"`
}

var _ Expression = &InvocationExpression{}

func (*InvocationExpression) isExpression() {}

func (*InvocationExpression) isIfStatementTest() {}

func (e *InvocationExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *InvocationExpression) Walk(walkChild func(Element)) {
	walkChild(e.InvokedExpression)
	for _, argument := range e.Arguments {
		walkChild(argument.Expression)
	}
}

func (e *InvocationExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitInvocationExpression(e)
}

func (e *InvocationExpression) String() string {
	var builder strings.Builder
	builder.WriteString(e.InvokedExpression.String())
	if len(e.TypeArguments) > 0 {
		builder.WriteRune('<')
		for i, ty := range e.TypeArguments {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(ty.String())
		}
		builder.WriteRune('>')
	}
	builder.WriteString(e.Arguments.String())
	return builder.String()
}

func (e *InvocationExpression) Doc() prettier.Doc {

	result := prettier.Concat{
		// TODO: potentially parenthesize
		e.InvokedExpression.Doc(),
	}

	if len(e.TypeArguments) > 0 {
		typeArgumentDocs := make([]prettier.Doc, len(e.TypeArguments))
		for i, typeArgument := range e.TypeArguments {
			typeArgumentDocs[i] = typeArgument.Doc()
		}

		result = append(result,
			prettier.Wrap(
				prettier.Text("<"),
				prettier.Join(arrayExpressionSeparatorDoc, typeArgumentDocs...),
				prettier.Text(">"),
				prettier.SoftLine{},
			),
		)
	}

	var argumentsDoc prettier.Doc
	if len(e.Arguments) == 0 {
		argumentsDoc = prettier.Text("()")
	} else {
		argumentDocs := make([]prettier.Doc, len(e.Arguments))
		for i, argument := range e.Arguments {
			argumentDoc := argument.Expression.Doc()
			if argument.Label != "" {
				argumentDoc = prettier.Concat{
					prettier.Text(argument.Label + ": "),
					argumentDoc,
				}
			}
			argumentDocs[i] = argumentDoc
		}
		argumentsDoc = prettier.WrapParentheses(
			prettier.Join(arrayExpressionSeparatorDoc, argumentDocs...),
			prettier.SoftLine{},
		)
	}

	result = append(result, argumentsDoc)

	return result
}

func (e *InvocationExpression) StartPosition() Position {
	return e.InvokedExpression.StartPosition()
}

func (e *InvocationExpression) EndPosition() Position {
	return e.EndPos
}

func (e *InvocationExpression) MarshalJSON() ([]byte, error) {
	type Alias InvocationExpression
	return json.Marshal(&struct {
		Type string
		*Alias
		Range
	}{
		Type:  "InvocationExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// AccessExpression

type AccessExpression interface {
	Expression
	isAccessExpression()
	AccessedExpression() Expression
}

// MemberExpression

type MemberExpression struct {
	Expression Expression
	Optional   bool
	// The position of the token (`.`, `?.`) that separates the accessed expression
	// and the identifier of the member
	AccessPos  Position
	Identifier Identifier
}

var _ Expression = &MemberExpression{}

func (*MemberExpression) isExpression() {}

func (*MemberExpression) isIfStatementTest() {}

func (*MemberExpression) isAccessExpression() {}

func (e *MemberExpression) AccessedExpression() Expression {
	return e.Expression
}

func (e *MemberExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *MemberExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
}

func (e *MemberExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitMemberExpression(e)
}

func (e *MemberExpression) String() string {
	optional := ""
	if e.Optional {
		optional = "?"
	}
	return fmt.Sprintf(
		"%s%s.%s",
		e.Expression, optional, e.Identifier,
	)
}

var memberExpressionSeparatorDoc prettier.Doc = prettier.Text(".")
var memberExpressionOptionalSeparatorDoc prettier.Doc = prettier.Text("?.")

func (e *MemberExpression) Doc() prettier.Doc {
	var separatorDoc prettier.Doc
	if e.Optional {
		separatorDoc = memberExpressionOptionalSeparatorDoc
	} else {
		separatorDoc = memberExpressionSeparatorDoc
	}
	return prettier.Concat{
		// TODO: potentially parenthesize
		e.Expression.Doc(),
		prettier.Group{
			Doc: prettier.Indent{
				Doc: prettier.Concat{
					prettier.SoftLine{},
					separatorDoc,
					prettier.Text(e.Identifier.Identifier),
				},
			},
		},
	}
}

func (e *MemberExpression) StartPosition() Position {
	return e.Expression.StartPosition()
}

func (e *MemberExpression) EndPosition() Position {
	if e.Identifier.Identifier == "" {
		return e.AccessPos
	} else {
		return e.Identifier.EndPosition()
	}
}

func (e *MemberExpression) MarshalJSON() ([]byte, error) {
	type Alias MemberExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "MemberExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// IndexExpression

type IndexExpression struct {
	TargetExpression   Expression
	IndexingExpression Expression
	Range
}

var _ Expression = &IndexExpression{}

func (*IndexExpression) isExpression() {}

func (*IndexExpression) isIfStatementTest() {}

func (*IndexExpression) isAccessExpression() {}

func (e *IndexExpression) AccessedExpression() Expression {
	return e.TargetExpression
}

func (e *IndexExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *IndexExpression) Walk(walkChild func(Element)) {
	walkChild(e.TargetExpression)
	walkChild(e.IndexingExpression)
}

func (e *IndexExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitIndexExpression(e)
}
func (e *IndexExpression) String() string {
	return fmt.Sprintf(
		"%s[%s]",
		e.TargetExpression,
		e.IndexingExpression,
	)
}

func (e *IndexExpression) Doc() prettier.Doc {
	return prettier.Concat{
		// TODO: potentially parenthesize
		e.TargetExpression.Doc(),
		prettier.WrapBrackets(
			e.IndexingExpression.Doc(),
			prettier.SoftLine{},
		),
	}
}

func (e *IndexExpression) MarshalJSON() ([]byte, error) {
	type Alias IndexExpression
	return json.Marshal(&struct {
		Type string
		*Alias
	}{
		Type:  "IndexExpression",
		Alias: (*Alias)(e),
	})
}

// ConditionalExpression

type ConditionalExpression struct {
	Test Expression
	Then Expression
	Else Expression
}

var _ Expression = &ConditionalExpression{}

func (*ConditionalExpression) isExpression() {}

func (*ConditionalExpression) isIfStatementTest() {}

func (e *ConditionalExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *ConditionalExpression) Walk(walkChild func(Element)) {
	walkChild(e.Test)
	walkChild(e.Then)
	if e.Else != nil {
		walkChild(e.Else)
	}
}

func (e *ConditionalExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitConditionalExpression(e)
}
func (e *ConditionalExpression) String() string {
	return fmt.Sprintf(
		"(%s ? %s : %s)",
		e.Test, e.Then, e.Else,
	)
}

var conditionalExpressionTestSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Line{},
	prettier.Text("? "),
}
var conditionalExpressionBranchSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Line{},
	prettier.Text(": "),
}

func (e *ConditionalExpression) Doc() prettier.Doc {
	// TODO: potentially parenthesize
	testDoc := e.Test.Doc()

	// TODO: potentially parenthesize
	thenDoc := e.Then.Doc()

	// TODO: potentially parenthesize
	elseDoc := e.Else.Doc()

	return prettier.Group{
		Doc: prettier.Concat{
			testDoc,
			prettier.Indent{
				Doc: prettier.Concat{
					conditionalExpressionTestSeparatorDoc,
					prettier.Indent{
						Doc: thenDoc,
					},
					conditionalExpressionBranchSeparatorDoc,
					prettier.Indent{
						Doc: elseDoc,
					},
				},
			},
		},
	}
}

func (e *ConditionalExpression) StartPosition() Position {
	return e.Test.StartPosition()
}

func (e *ConditionalExpression) EndPosition() Position {
	return e.Else.EndPosition()
}

func (e *ConditionalExpression) MarshalJSON() ([]byte, error) {
	type Alias ConditionalExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "ConditionalExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// UnaryExpression

type UnaryExpression struct {
	Operation  Operation
	Expression Expression
	StartPos   Position `json:"-"`
}

var _ Expression = &UnaryExpression{}

func (*UnaryExpression) isExpression() {}

func (*UnaryExpression) isIfStatementTest() {}

func (e *UnaryExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *UnaryExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
}

func (e *UnaryExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitUnaryExpression(e)
}

func (e *UnaryExpression) String() string {
	return fmt.Sprintf(
		"%s%s",
		e.Operation.Symbol(),
		e.Expression,
	)
}

func (e *UnaryExpression) Doc() prettier.Doc {
	return prettier.Concat{
		prettier.Text(e.Operation.Symbol()),
		// TODO: potentially parenthesize
		e.Expression.Doc(),
	}
}

func (e *UnaryExpression) StartPosition() Position {
	return e.StartPos
}

func (e *UnaryExpression) EndPosition() Position {
	return e.Expression.EndPosition()
}

func (e *UnaryExpression) MarshalJSON() ([]byte, error) {
	type Alias UnaryExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "UnaryExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// BinaryExpression

type BinaryExpression struct {
	Operation Operation
	Left      Expression
	Right     Expression
}

var _ Expression = &BinaryExpression{}

func (*BinaryExpression) isExpression() {}

func (*BinaryExpression) isIfStatementTest() {}

func (e *BinaryExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *BinaryExpression) Walk(walkChild func(Element)) {
	walkChild(e.Left)
	walkChild(e.Right)
}

func (e *BinaryExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitBinaryExpression(e)
}

func (e *BinaryExpression) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		e.Left, e.Operation.Symbol(), e.Right,
	)
}

func (e *BinaryExpression) Doc() prettier.Doc {
	// TODO: potentially parenthesize
	leftDoc := e.Left.Doc()

	// TODO: potentially parenthesize
	rightDoc := e.Right.Doc()

	return prettier.Group{
		Doc: prettier.Concat{
			prettier.Group{
				Doc: leftDoc,
			},
			prettier.Line{},
			prettier.Text(e.Operation.Symbol()),
			prettier.Space,
			prettier.Group{
				Doc: rightDoc,
			},
		},
	}
}

func (e *BinaryExpression) StartPosition() Position {
	return e.Left.StartPosition()
}

func (e *BinaryExpression) EndPosition() Position {
	return e.Right.EndPosition()
}

func (e *BinaryExpression) MarshalJSON() ([]byte, error) {
	type Alias BinaryExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "BinaryExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// FunctionExpression

type FunctionExpression struct {
	ParameterList        *ParameterList
	ReturnTypeAnnotation *TypeAnnotation
	FunctionBlock        *FunctionBlock
	StartPos             Position `json:"-"`
}

var _ Expression = &FunctionExpression{}

func (*FunctionExpression) isExpression() {}

func (*FunctionExpression) isIfStatementTest() {}

func (e *FunctionExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *FunctionExpression) Walk(walkChild func(Element)) {
	// TODO: walk parameters
	// TODO: walk return type
	walkChild(e.FunctionBlock)
}

func (e *FunctionExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitFunctionExpression(e)
}

func (e *FunctionExpression) String() string {
	// TODO:
	return "func ..."
}

var functionExpressionFunKeywordDoc prettier.Doc = prettier.Text("fun ")
var functionExpressionParameterSeparatorDoc prettier.Doc = prettier.Concat{
	prettier.Text(","),
	prettier.Line{},
}

var typeSeparatorDoc prettier.Doc = prettier.Text(": ")
var functionExpressionEmptyBlockDoc prettier.Doc = prettier.Text(" {}")

func (e *FunctionExpression) Doc() prettier.Doc {

	signatureDoc := e.parametersDoc()

	if e.ReturnTypeAnnotation != nil &&
		!IsEmptyType(e.ReturnTypeAnnotation.Type) {

		signatureDoc = prettier.Concat{
			signatureDoc,
			typeSeparatorDoc,
			e.ReturnTypeAnnotation.Doc(),
		}
	}

	doc := prettier.Concat{
		functionExpressionFunKeywordDoc,
		prettier.Group{
			Doc: signatureDoc,
		},
	}

	if e.FunctionBlock.IsEmpty() {
		return append(doc, functionExpressionEmptyBlockDoc)
	} else {
		// TODO: pre-conditions
		// TODO: post-conditions

		blockDoc := e.FunctionBlock.Block.Doc()

		return append(
			doc,
			prettier.Space,
			blockDoc,
		)
	}
}

func (e *FunctionExpression) parametersDoc() prettier.Doc {

	if e.ParameterList == nil ||
		len(e.ParameterList.Parameters) == 0 {

		return prettier.Text("()")
	}

	parameterDocs := make([]prettier.Doc, 0, len(e.ParameterList.Parameters))

	for _, parameter := range e.ParameterList.Parameters {
		var parameterDoc prettier.Concat

		if parameter.Label != "" {
			parameterDoc = append(parameterDoc,
				prettier.Text(parameter.Label),
				prettier.Space,
			)
		}

		parameterDoc = append(
			parameterDoc,
			prettier.Text(parameter.Identifier.Identifier),
			typeSeparatorDoc,
			parameter.TypeAnnotation.Doc(),
		)

		parameterDocs = append(parameterDocs, parameterDoc)
	}

	return prettier.WrapParentheses(
		prettier.Join(
			functionExpressionParameterSeparatorDoc,
			parameterDocs...,
		),
		prettier.SoftLine{},
	)
}

func (e *FunctionExpression) StartPosition() Position {
	return e.StartPos
}

func (e *FunctionExpression) EndPosition() Position {
	return e.FunctionBlock.EndPosition()
}

func (e *FunctionExpression) MarshalJSON() ([]byte, error) {
	type Alias FunctionExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "FunctionExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// CastingExpression

type CastingExpression struct {
	Expression                Expression
	Operation                 Operation
	TypeAnnotation            *TypeAnnotation
	ParentVariableDeclaration *VariableDeclaration `json:"-"`
}

var _ Expression = &CastingExpression{}

func (*CastingExpression) isExpression() {}

func (*CastingExpression) isIfStatementTest() {}

func (e *CastingExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}
func (e *CastingExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
	// TODO: also walk type
}

func (e *CastingExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitCastingExpression(e)
}

func (e *CastingExpression) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		e.Expression, e.Operation.Symbol(), e.TypeAnnotation,
	)
}

func (e *CastingExpression) Doc() prettier.Doc {
	// TODO: potentially parenthesize
	doc := e.Expression.Doc()

	return prettier.Group{
		Doc: prettier.Concat{
			prettier.Group{
				Doc: doc,
			},
			prettier.Line{},
			prettier.Text(e.Operation.Symbol()),
			prettier.Line{},
			e.TypeAnnotation.Doc(),
		},
	}
}

func (e *CastingExpression) StartPosition() Position {
	return e.Expression.StartPosition()
}

func (e *CastingExpression) EndPosition() Position {
	return e.TypeAnnotation.EndPosition()
}

func (e *CastingExpression) MarshalJSON() ([]byte, error) {
	type Alias CastingExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "CastingExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// CreateExpression

type CreateExpression struct {
	InvocationExpression *InvocationExpression
	StartPos             Position `json:"-"`
}

var _ Expression = &CreateExpression{}

func (*CreateExpression) isExpression() {}

func (*CreateExpression) isIfStatementTest() {}

func (e *CreateExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *CreateExpression) Walk(walkChild func(Element)) {
	walkChild(e.InvocationExpression)
}

func (e *CreateExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitCreateExpression(e)
}

func (e *CreateExpression) String() string {
	return fmt.Sprintf(
		"(create %s)",
		e.InvocationExpression,
	)
}

func (e *CreateExpression) Doc() prettier.Doc {
	return prettier.Concat{
		prettier.Text("create "),
		// TODO: potentially parenthesize
		e.InvocationExpression.Doc(),
	}
}

func (e *CreateExpression) StartPosition() Position {
	return e.StartPos
}

func (e *CreateExpression) EndPosition() Position {
	return e.InvocationExpression.EndPos
}

func (e *CreateExpression) MarshalJSON() ([]byte, error) {
	type Alias CreateExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "CreateExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// DestroyExpression

type DestroyExpression struct {
	Expression Expression
	StartPos   Position `json:"-"`
}

var _ Expression = &DestroyExpression{}

func (*DestroyExpression) isExpression() {}

func (*DestroyExpression) isIfStatementTest() {}

func (e *DestroyExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *DestroyExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
}

func (e *DestroyExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitDestroyExpression(e)
}

func (e *DestroyExpression) String() string {
	return fmt.Sprintf(
		"(destroy %s)",
		e.Expression,
	)
}

const destroyExpressionKeywordDoc = prettier.Text("destroy ")

func (e *DestroyExpression) Doc() prettier.Doc {
	return prettier.Concat{
		destroyExpressionKeywordDoc,
		// TODO: potentially parenthesize
		e.Expression.Doc(),
	}
}

func (e *DestroyExpression) StartPosition() Position {
	return e.StartPos
}

func (e *DestroyExpression) EndPosition() Position {
	return e.Expression.EndPosition()
}

func (e *DestroyExpression) MarshalJSON() ([]byte, error) {
	type Alias DestroyExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "DestroyExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// ReferenceExpression

type ReferenceExpression struct {
	Expression Expression
	Type       Type     `json:"TargetType"`
	StartPos   Position `json:"-"`
}

var _ Expression = &ReferenceExpression{}

func (*ReferenceExpression) isExpression() {}

func (*ReferenceExpression) isIfStatementTest() {}

func (e *ReferenceExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *ReferenceExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
	// TODO: walk type
}

func (e *ReferenceExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitReferenceExpression(e)
}

func (e *ReferenceExpression) String() string {
	return fmt.Sprintf(
		"(&%s as %s)",
		e.Expression,
		e.Type,
	)
}

var referenceExpressionRefOperatorDoc prettier.Doc = prettier.Text("&")
var referenceExpressionAsOperatorDoc prettier.Doc = prettier.Text("as")

func (e *ReferenceExpression) Doc() prettier.Doc {
	// TODO: potentially parenthesize
	doc := e.Expression.Doc()

	return prettier.Group{
		Doc: prettier.Concat{
			referenceExpressionRefOperatorDoc,
			prettier.Group{
				Doc: doc,
			},
			prettier.Line{},
			referenceExpressionAsOperatorDoc,
			prettier.Line{},
			e.Type.Doc(),
		},
	}
}

func (e *ReferenceExpression) StartPosition() Position {
	return e.StartPos
}

func (e *ReferenceExpression) EndPosition() Position {
	return e.Type.EndPosition()
}

func (e *ReferenceExpression) MarshalJSON() ([]byte, error) {
	type Alias ReferenceExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "ReferenceExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// ForceExpression

type ForceExpression struct {
	Expression Expression
	EndPos     Position `json:"-"`
}

var _ Expression = &ForceExpression{}

func (*ForceExpression) isExpression() {}

func (*ForceExpression) isIfStatementTest() {}

func (e *ForceExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (e *ForceExpression) Walk(walkChild func(Element)) {
	walkChild(e.Expression)
}

func (e *ForceExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitForceExpression(e)
}

func (e *ForceExpression) String() string {
	return fmt.Sprintf("%s!", e.Expression)
}

const forceExpressionOperatorDoc = prettier.Text("!")

func (e *ForceExpression) Doc() prettier.Doc {
	return prettier.Concat{
		// TODO: potentially parenthesize
		e.Expression.Doc(),
		forceExpressionOperatorDoc,
	}
}

func (e *ForceExpression) StartPosition() Position {
	return e.Expression.StartPosition()
}

func (e *ForceExpression) EndPosition() Position {
	return e.EndPos
}

func (e *ForceExpression) MarshalJSON() ([]byte, error) {
	type Alias ForceExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "ForceExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}

// PathExpression

type PathExpression struct {
	StartPos   Position `json:"-"`
	Domain     Identifier
	Identifier Identifier
}

var _ Expression = &PathExpression{}

func (*PathExpression) isExpression() {}

func (*PathExpression) isIfStatementTest() {}

func (e *PathExpression) Accept(visitor Visitor) Repr {
	return e.AcceptExp(visitor)
}

func (*PathExpression) Walk(_ func(Element)) {
	// NO-OP
}

func (e *PathExpression) AcceptExp(visitor ExpressionVisitor) Repr {
	return visitor.VisitPathExpression(e)
}

func (e *PathExpression) String() string {
	return fmt.Sprintf("/%s/%s", e.Domain, e.Identifier)
}

func (e *PathExpression) Doc() prettier.Doc {
	return prettier.Text(e.String())
}

func (e *PathExpression) StartPosition() Position {
	return e.StartPos
}

func (e *PathExpression) EndPosition() Position {
	return e.Identifier.EndPosition()
}

func (e *PathExpression) MarshalJSON() ([]byte, error) {
	type Alias PathExpression
	return json.Marshal(&struct {
		Type string
		Range
		*Alias
	}{
		Type:  "PathExpression",
		Range: NewRangeFromPositioned(e),
		Alias: (*Alias)(e),
	})
}
