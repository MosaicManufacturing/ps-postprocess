// Code generated from SequenceParser.g4 by ANTLR 4.9. DO NOT EDIT.

package interpreter // SequenceParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by SequenceParser.
type SequenceParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SequenceParser#sequence.
	VisitSequence(ctx *SequenceContext) interface{}

	// Visit a parse tree produced by SequenceParser#statements.
	VisitStatements(ctx *StatementsContext) interface{}

	// Visit a parse tree produced by SequenceParser#ifBlock.
	VisitIfBlock(ctx *IfBlockContext) interface{}

	// Visit a parse tree produced by SequenceParser#optionalElseBlock.
	VisitOptionalElseBlock(ctx *OptionalElseBlockContext) interface{}

	// Visit a parse tree produced by SequenceParser#whileBlock.
	VisitWhileBlock(ctx *WhileBlockContext) interface{}

	// Visit a parse tree produced by SequenceParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by SequenceParser#assignment.
	VisitAssignment(ctx *AssignmentContext) interface{}

	// Visit a parse tree produced by SequenceParser#gCode.
	VisitGCode(ctx *GCodeContext) interface{}

	// Visit a parse tree produced by SequenceParser#gCodeText.
	VisitGCodeText(ctx *GCodeTextContext) interface{}

	// Visit a parse tree produced by SequenceParser#gCodeEscapedText.
	VisitGCodeEscapedText(ctx *GCodeEscapedTextContext) interface{}

	// Visit a parse tree produced by SequenceParser#gCodeSubExpression.
	VisitGCodeSubExpression(ctx *GCodeSubExpressionContext) interface{}

	// Visit a parse tree produced by SequenceParser#identExpr.
	VisitIdentExpr(ctx *IdentExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#floatExpr.
	VisitFloatExpr(ctx *FloatExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#unaryOpExpr.
	VisitUnaryOpExpr(ctx *UnaryOpExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#intExpr.
	VisitIntExpr(ctx *IntExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by SequenceParser#binaryOpExpr.
	VisitBinaryOpExpr(ctx *BinaryOpExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#boolExpr.
	VisitBoolExpr(ctx *BoolExprContext) interface{}

	// Visit a parse tree produced by SequenceParser#parenExpr.
	VisitParenExpr(ctx *ParenExprContext) interface{}
}
