package ast

import "testing"

func TestStatementNodes(_ *testing.T) {
	statements := []Statement{
		&LetStatement{},
		&ReturnStatement{},
		&ExpressionStatement{},
		&Block{},
	}

	for _, s := range statements {
		s.statementNode()
	}
}

func TestExpressionNodes(_ *testing.T) {
	expressions := []Expression{
		&Identifier{},
		&IntegerLiteral{},
		&PrefixExpression{},
		&InfixExpression{},
		&Boolean{},
		&IfExpression{},
		&Block{},
		&FunctionLiteral{},
		&CallExpression{},
		&StringLiteral{},
		&ArrayLiteral{},
		&IndexExpression{},
		&HashMapLiteral{},
	}

	for _, e := range expressions {
		e.expressionNode()
	}
}
