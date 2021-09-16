// Code generated from SequenceParser.g4 by ANTLR 4.9. DO NOT EDIT.

package printerscript // SequenceParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// SequenceParserListener is a complete listener for a parse tree produced by SequenceParser.
type SequenceParserListener interface {
	antlr.ParseTreeListener

	// EnterSequence is called when entering the sequence production.
	EnterSequence(c *SequenceContext)

	// EnterStatements is called when entering the statements production.
	EnterStatements(c *StatementsContext)

	// EnterIfBlock is called when entering the ifBlock production.
	EnterIfBlock(c *IfBlockContext)

	// EnterOptionalElseBlock is called when entering the optionalElseBlock production.
	EnterOptionalElseBlock(c *OptionalElseBlockContext)

	// EnterWhileBlock is called when entering the whileBlock production.
	EnterWhileBlock(c *WhileBlockContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterGCode is called when entering the gCode production.
	EnterGCode(c *GCodeContext)

	// EnterGCodeText is called when entering the gCodeText production.
	EnterGCodeText(c *GCodeTextContext)

	// EnterGCodeEscapedText is called when entering the gCodeEscapedText production.
	EnterGCodeEscapedText(c *GCodeEscapedTextContext)

	// EnterGCodeSubExpression is called when entering the gCodeSubExpression production.
	EnterGCodeSubExpression(c *GCodeSubExpressionContext)

	// EnterIdentExpr is called when entering the identExpr production.
	EnterIdentExpr(c *IdentExprContext)

	// EnterFloatExpr is called when entering the floatExpr production.
	EnterFloatExpr(c *FloatExprContext)

	// EnterUnaryOpExpr is called when entering the unaryOpExpr production.
	EnterUnaryOpExpr(c *UnaryOpExprContext)

	// EnterIntExpr is called when entering the intExpr production.
	EnterIntExpr(c *IntExprContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterBinaryOpExpr is called when entering the binaryOpExpr production.
	EnterBinaryOpExpr(c *BinaryOpExprContext)

	// EnterBoolExpr is called when entering the boolExpr production.
	EnterBoolExpr(c *BoolExprContext)

	// EnterParenExpr is called when entering the parenExpr production.
	EnterParenExpr(c *ParenExprContext)

	// ExitSequence is called when exiting the sequence production.
	ExitSequence(c *SequenceContext)

	// ExitStatements is called when exiting the statements production.
	ExitStatements(c *StatementsContext)

	// ExitIfBlock is called when exiting the ifBlock production.
	ExitIfBlock(c *IfBlockContext)

	// ExitOptionalElseBlock is called when exiting the optionalElseBlock production.
	ExitOptionalElseBlock(c *OptionalElseBlockContext)

	// ExitWhileBlock is called when exiting the whileBlock production.
	ExitWhileBlock(c *WhileBlockContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitGCode is called when exiting the gCode production.
	ExitGCode(c *GCodeContext)

	// ExitGCodeText is called when exiting the gCodeText production.
	ExitGCodeText(c *GCodeTextContext)

	// ExitGCodeEscapedText is called when exiting the gCodeEscapedText production.
	ExitGCodeEscapedText(c *GCodeEscapedTextContext)

	// ExitGCodeSubExpression is called when exiting the gCodeSubExpression production.
	ExitGCodeSubExpression(c *GCodeSubExpressionContext)

	// ExitIdentExpr is called when exiting the identExpr production.
	ExitIdentExpr(c *IdentExprContext)

	// ExitFloatExpr is called when exiting the floatExpr production.
	ExitFloatExpr(c *FloatExprContext)

	// ExitUnaryOpExpr is called when exiting the unaryOpExpr production.
	ExitUnaryOpExpr(c *UnaryOpExprContext)

	// ExitIntExpr is called when exiting the intExpr production.
	ExitIntExpr(c *IntExprContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitBinaryOpExpr is called when exiting the binaryOpExpr production.
	ExitBinaryOpExpr(c *BinaryOpExprContext)

	// ExitBoolExpr is called when exiting the boolExpr production.
	ExitBoolExpr(c *BoolExprContext)

	// ExitParenExpr is called when exiting the parenExpr production.
	ExitParenExpr(c *ParenExprContext)
}
