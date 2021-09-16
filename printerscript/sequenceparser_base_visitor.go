// Code generated from SequenceParser.g4 by ANTLR 4.9. DO NOT EDIT.

package printerscript // SequenceParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseSequenceParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSequenceParserVisitor) VisitSequence(ctx *SequenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitStatements(ctx *StatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitIfBlock(ctx *IfBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitOptionalElseBlock(ctx *OptionalElseBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitWhileBlock(ctx *WhileBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitAssignment(ctx *AssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitGCode(ctx *GCodeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitGCodeText(ctx *GCodeTextContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitGCodeEscapedText(ctx *GCodeEscapedTextContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitGCodeSubExpression(ctx *GCodeSubExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitIdentExpr(ctx *IdentExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitFloatExpr(ctx *FloatExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitUnaryOpExpr(ctx *UnaryOpExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitIntExpr(ctx *IntExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitBinaryOpExpr(ctx *BinaryOpExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitBoolExpr(ctx *BoolExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSequenceParserVisitor) VisitParenExpr(ctx *ParenExprContext) interface{} {
	return v.VisitChildren(ctx)
}
