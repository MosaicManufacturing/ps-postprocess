// Code generated from SequenceParser.g4 by ANTLR 4.9. DO NOT EDIT.

package interpreter // SequenceParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSequenceParserListener is a complete listener for a parse tree produced by SequenceParser.
type BaseSequenceParserListener struct{}

var _ SequenceParserListener = &BaseSequenceParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSequenceParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSequenceParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSequenceParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSequenceParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSequence is called when production sequence is entered.
func (s *BaseSequenceParserListener) EnterSequence(ctx *SequenceContext) {}

// ExitSequence is called when production sequence is exited.
func (s *BaseSequenceParserListener) ExitSequence(ctx *SequenceContext) {}

// EnterStatements is called when production statements is entered.
func (s *BaseSequenceParserListener) EnterStatements(ctx *StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *BaseSequenceParserListener) ExitStatements(ctx *StatementsContext) {}

// EnterIfBlock is called when production ifBlock is entered.
func (s *BaseSequenceParserListener) EnterIfBlock(ctx *IfBlockContext) {}

// ExitIfBlock is called when production ifBlock is exited.
func (s *BaseSequenceParserListener) ExitIfBlock(ctx *IfBlockContext) {}

// EnterOptionalElseBlock is called when production optionalElseBlock is entered.
func (s *BaseSequenceParserListener) EnterOptionalElseBlock(ctx *OptionalElseBlockContext) {}

// ExitOptionalElseBlock is called when production optionalElseBlock is exited.
func (s *BaseSequenceParserListener) ExitOptionalElseBlock(ctx *OptionalElseBlockContext) {}

// EnterWhileBlock is called when production whileBlock is entered.
func (s *BaseSequenceParserListener) EnterWhileBlock(ctx *WhileBlockContext) {}

// ExitWhileBlock is called when production whileBlock is exited.
func (s *BaseSequenceParserListener) ExitWhileBlock(ctx *WhileBlockContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseSequenceParserListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseSequenceParserListener) ExitStatement(ctx *StatementContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseSequenceParserListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseSequenceParserListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterGCode is called when production gCode is entered.
func (s *BaseSequenceParserListener) EnterGCode(ctx *GCodeContext) {}

// ExitGCode is called when production gCode is exited.
func (s *BaseSequenceParserListener) ExitGCode(ctx *GCodeContext) {}

// EnterGCodeText is called when production gCodeText is entered.
func (s *BaseSequenceParserListener) EnterGCodeText(ctx *GCodeTextContext) {}

// ExitGCodeText is called when production gCodeText is exited.
func (s *BaseSequenceParserListener) ExitGCodeText(ctx *GCodeTextContext) {}

// EnterGCodeEscapedText is called when production gCodeEscapedText is entered.
func (s *BaseSequenceParserListener) EnterGCodeEscapedText(ctx *GCodeEscapedTextContext) {}

// ExitGCodeEscapedText is called when production gCodeEscapedText is exited.
func (s *BaseSequenceParserListener) ExitGCodeEscapedText(ctx *GCodeEscapedTextContext) {}

// EnterGCodeSubExpression is called when production gCodeSubExpression is entered.
func (s *BaseSequenceParserListener) EnterGCodeSubExpression(ctx *GCodeSubExpressionContext) {}

// ExitGCodeSubExpression is called when production gCodeSubExpression is exited.
func (s *BaseSequenceParserListener) ExitGCodeSubExpression(ctx *GCodeSubExpressionContext) {}

// EnterIdentExpr is called when production identExpr is entered.
func (s *BaseSequenceParserListener) EnterIdentExpr(ctx *IdentExprContext) {}

// ExitIdentExpr is called when production identExpr is exited.
func (s *BaseSequenceParserListener) ExitIdentExpr(ctx *IdentExprContext) {}

// EnterFloatExpr is called when production floatExpr is entered.
func (s *BaseSequenceParserListener) EnterFloatExpr(ctx *FloatExprContext) {}

// ExitFloatExpr is called when production floatExpr is exited.
func (s *BaseSequenceParserListener) ExitFloatExpr(ctx *FloatExprContext) {}

// EnterUnaryOpExpr is called when production unaryOpExpr is entered.
func (s *BaseSequenceParserListener) EnterUnaryOpExpr(ctx *UnaryOpExprContext) {}

// ExitUnaryOpExpr is called when production unaryOpExpr is exited.
func (s *BaseSequenceParserListener) ExitUnaryOpExpr(ctx *UnaryOpExprContext) {}

// EnterIntExpr is called when production intExpr is entered.
func (s *BaseSequenceParserListener) EnterIntExpr(ctx *IntExprContext) {}

// ExitIntExpr is called when production intExpr is exited.
func (s *BaseSequenceParserListener) ExitIntExpr(ctx *IntExprContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseSequenceParserListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseSequenceParserListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterBinaryOpExpr is called when production binaryOpExpr is entered.
func (s *BaseSequenceParserListener) EnterBinaryOpExpr(ctx *BinaryOpExprContext) {}

// ExitBinaryOpExpr is called when production binaryOpExpr is exited.
func (s *BaseSequenceParserListener) ExitBinaryOpExpr(ctx *BinaryOpExprContext) {}

// EnterBoolExpr is called when production boolExpr is entered.
func (s *BaseSequenceParserListener) EnterBoolExpr(ctx *BoolExprContext) {}

// ExitBoolExpr is called when production boolExpr is exited.
func (s *BaseSequenceParserListener) ExitBoolExpr(ctx *BoolExprContext) {}

// EnterParenExpr is called when production parenExpr is entered.
func (s *BaseSequenceParserListener) EnterParenExpr(ctx *ParenExprContext) {}

// ExitParenExpr is called when production parenExpr is exited.
func (s *BaseSequenceParserListener) ExitParenExpr(ctx *ParenExprContext) {}
