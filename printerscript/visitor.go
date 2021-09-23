package printerscript

import (
    "../gcode"
    "fmt"
    "github.com/antlr/antlr4/runtime/Go/antlr"
    "strconv"
)

type VisitorOptions struct {
    MaxLoopIterations int
    EOL string
    Locals map[string]float64
}

type Visitor struct {
    antlr.ParseTreeVisitor
    maxLoopIters int
    loopIters int
    EOL string
    locals map[string]float64
    result string
}

func (v *Visitor) GetLocal(key string) float64 {
    if value, ok := v.locals[key]; ok {
        return value
    }
    // treat undefined locals as if they are defined as 0
    return 0
}

func (v *Visitor) SetLocal(key string, value float64) {
    v.locals[key] = value
}

func (v *Visitor) GetResult() string {
    return v.result
}

func (v *Visitor) GetLocals() map[string]float64 {
    return v.locals
}

// https://github.com/nikunjy/antlr-calc-go/blob/master/parser/calculator_visitor_impl.go

func NewVisitor(opts VisitorOptions) Visitor {
    return Visitor{
        ParseTreeVisitor: &BaseSequenceParserVisitor{},
        maxLoopIters:     opts.MaxLoopIterations,
        EOL:              opts.EOL,
        locals:           opts.Locals,
    }
}

func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
    return tree.Accept(v)
}

func (v *Visitor) VisitSequence(ctx *SequenceContext) interface{} {
    if DEBUG { fmt.Println("VisitSequence") }
    return ctx.Statements().Accept(v)
}

func (v *Visitor) VisitStatements(ctx *StatementsContext) interface{} {
    if DEBUG { fmt.Println("VisitStatements") }
    if ifBlock := ctx.IfBlock(); ifBlock != nil {
        err := v.Visit(ifBlock)
        if runtimeError, ok := err.(*RuntimeError); ok {
            return runtimeError
        }
    } else if whileBlock := ctx.WhileBlock(); whileBlock != nil {
        err := v.Visit(whileBlock)
        if runtimeError, ok := err.(*RuntimeError); ok {
            return runtimeError
        }
    } else if statement := ctx.Statement(); statement != nil {
        err := v.Visit(statement)
        if runtimeError, ok := err.(*RuntimeError); ok {
           return runtimeError
        }
    }
    if statements := ctx.Statements(); statements != nil {
        err := ctx.Statements().Accept(v)
        if runtimeError, ok := err.(*RuntimeError); ok {
            return runtimeError
        }
    }
    return nil
}

func (v *Visitor) VisitStatement(ctx *StatementContext) interface{} {
    if DEBUG { fmt.Println("VisitStatement") }
    assignment := ctx.Assignment()
    if assignment != nil {
        err := v.Visit(assignment)
        if runtimeError, ok := err.(*RuntimeError); ok {
            return runtimeError
        }
        return nil
    }
    err := v.Visit(ctx.GCode())
    if runtimeError, ok := err.(*RuntimeError); ok {
        return runtimeError
    }
    return nil
}

func (v *Visitor) VisitIfBlock(ctx *IfBlockContext) interface{} {
    if DEBUG { fmt.Println("VisitIfBlock") }
    condition := ctx.Expression()
    enterIf := v.Visit(condition)
    if runtimeError, ok := enterIf.(*RuntimeError); ok {
        return runtimeError
    }
    if enterIf.(float64) != 0 {
        // if condition is true
        return v.Visit(ctx.Statements())
    }
    // if condition is false -- check for else block
    elseBlock := ctx.OptionalElseBlock()
    if elseBlock != nil {
        return v.Visit(elseBlock)
    }
    return nil
}

func (v *Visitor) VisitOptionalElseBlock(ctx *OptionalElseBlockContext) interface{} {
    if DEBUG { fmt.Println("VisitOptionalElseBlock") }
    if ifBlock := ctx.IfBlock(); ifBlock != nil {
        // else-if block
        return v.Visit(ctx.IfBlock())
    }
    if statements := ctx.Statements(); statements != nil {
        // else block
        return v.Visit(ctx.Statements())
    }
    // if or else-if without a final else
    return nil
}

func (v *Visitor) VisitWhileBlock(ctx *WhileBlockContext) interface{} {
    if DEBUG { fmt.Println("VisitWhileBlock") }
    condition := ctx.Expression()
    enterWhile := v.Visit(condition)
    if runtimeError, ok := enterWhile.(*RuntimeError); ok {
        return runtimeError
    }
    for enterWhile.(float64) != 0 {
        if v.loopIters > v.maxLoopIters {
            start := ctx.GetStart()
            line := start.GetLine()
            col := start.GetColumn()
            return NewRuntimeError("maximum number of loop iterations exceeded", line, col)
        }
        v.Visit(ctx.Statements())
        enterWhile = v.Visit(condition)
        if runtimeError, ok := enterWhile.(*RuntimeError); ok {
            return runtimeError
        }
        v.loopIters++
    }
    return nil
}

func (v *Visitor) VisitAssignment(ctx *AssignmentContext) interface{} {
    if DEBUG { fmt.Println("VisitAssignment") }
    identifier := ctx.IDENTIFIER().GetText()
    value := v.Visit(ctx.Expression())
    if runtimeError, ok := value.(*RuntimeError); ok {
        return runtimeError
    }
    v.SetLocal(identifier, value.(float64))
    return nil
}

func (v *Visitor) VisitGCode(ctx *GCodeContext) interface{} {
    if DEBUG { fmt.Println("VisitGCode") }
    for _, part := range ctx.AllGCodePart() {
        err := v.Visit(part)
        if runtimeError, ok := err.(*RuntimeError); ok {
            return runtimeError
        }
    }
    v.result += v.EOL
    return nil
}

func (v *Visitor) VisitGCodeText(ctx *GCodeTextContext) interface{} {
    if DEBUG { fmt.Println("VisitGCodeText") }
    v.result += ctx.GetText()
    return nil
}

func (v *Visitor) VisitGCodeEscapedText(ctx *GCodeEscapedTextContext) interface{} {
    if DEBUG { fmt.Println("VisitGCodeEscapedText") }
    // remove leading backslash if necessary
    text := ctx.GetText()
    if len(text) > 0 {
        if text[0] == '\\' {
            v.result += text[1:]
        } else {
            v.result += text
        }
    }
    return nil
}

func (v *Visitor) VisitGCodeSubExpression(ctx *GCodeSubExpressionContext) interface{} {
    if DEBUG { fmt.Println("VisitGCodeSubExpression") }
    value := v.Visit(ctx.Expression())
    if runtimeError, ok := value.(*RuntimeError); ok {
        return runtimeError
    }
    v.result += gcode.FormatFloat(value.(float64))
    return nil
}

func (v *Visitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
    if DEBUG { fmt.Println("VisitFunctionCall") }
    // name of the function
    fn := ctx.IDENTIFIER().GetText()
    // arguments to the function
    paramCtxs := ctx.AllExpression()

    // validate arity
    argc := len(paramCtxs)
    requiredArity, err := getArity(fn)
    if err != nil {
        start := ctx.GetStart()
        line := start.GetLine()
        col := start.GetColumn()
        return NewRuntimeError(err.Error(), line, col)
    }
    argumentsMsg := "argument"
    if requiredArity != 1 {
        argumentsMsg = "arguments"
    }
    if fn == "max" || fn == "min" {
        // argc must be AT LEAST the required arity
        if argc < requiredArity {
            start := ctx.GetStart()
            line := start.GetLine()
            col := start.GetColumn()
            return NewRuntimeError(fmt.Sprintf("expected at least %d %s to '%s' (%d given)", requiredArity, argumentsMsg, fn, argc), line, col)
        }
    } else {
        // argc must be EXACTLY the required arity
        if argc != requiredArity {
            start := ctx.GetStart()
            line := start.GetLine()
            col := start.GetColumn()
            return NewRuntimeError(fmt.Sprintf("expected %d %s to '%s' (%d given)", requiredArity, argumentsMsg, fn, argc), line, col)
        }
    }

    // evaluate
    argv := make([]float64, 0, len(paramCtxs))
    for _, paramCtx := range paramCtxs {
        value := v.Visit(paramCtx)
        if runtimeError, ok := value.(*RuntimeError); ok {
            return runtimeError
        }
        argv = append(argv, value.(float64))
    }
    retval, err := evaluateFunction(fn, argv)
    if err != nil {
        start := ctx.GetStart()
        line := start.GetLine()
        col := start.GetColumn()
        return NewRuntimeError(err.Error(), line, col)
    }
    return retval
}

func (v *Visitor) VisitIdentExpr(ctx *IdentExprContext) interface{} {
    if DEBUG { fmt.Println("VisitIdentExpr") }
    identifier := ctx.IDENTIFIER().GetText()
    value := v.GetLocal(identifier)
    return value
}

func (v *Visitor) VisitIntExpr(ctx *IntExprContext) interface{} {
    if DEBUG { fmt.Println("VisitIntExpr") }
    // successful lexing + parsing guarantee no error here
    value, _ := strconv.ParseInt(ctx.INT().GetText(), 10, 64)
    return float64(value)
}

func (v *Visitor) VisitFloatExpr(ctx *FloatExprContext) interface{} {
    if DEBUG { fmt.Println("VisitFloatExpr") }
    // successful lexing + parsing guarantee no error here
    value, _ := strconv.ParseFloat(ctx.FLOAT().GetText(), 64)
    return value
}

func (v *Visitor) VisitBoolExpr(ctx *BoolExprContext) interface{} {
    if DEBUG { fmt.Println("VisitBoolExpr") }
    if ctx.TRUE() != nil {
        return float64(1)
    }
    return float64(0)
}

func (v *Visitor) VisitParenExpr(ctx *ParenExprContext) interface{} {
    if DEBUG { fmt.Println("VisitParenExpr") }
    return v.Visit(ctx.Expression())
}

func (v *Visitor) VisitUnaryOpExpr(ctx *UnaryOpExprContext) interface{} {
    if DEBUG { fmt.Println("VisitUnaryOpExpr") }
    op := ctx.GetChild(0).GetPayload().(antlr.Token).GetText()
    value := v.Visit(ctx.Expression())
    if runtimeError, ok := value.(*RuntimeError); ok {
        return runtimeError
    }
    result, err := evaluateUnaryOp(op, value.(float64))
    if err != nil {
        start := ctx.GetStart()
        line := start.GetLine()
        col := start.GetColumn()
        return NewRuntimeError(err.Error(), line, col)
    }
    return result
}

func (v *Visitor) VisitBinaryOpExpr(ctx *BinaryOpExprContext) interface{} {
    if DEBUG { fmt.Println("VisitBinaryOpExpr") }
    op := ctx.GetChild(1).GetPayload().(antlr.Token).GetText()
    lhs := v.Visit(ctx.Expression(0))
    if runtimeError, ok := lhs.(*RuntimeError); ok {
        return runtimeError
    }
    rhs := v.Visit(ctx.Expression(1))
    if runtimeError, ok := rhs.(*RuntimeError); ok {
        return runtimeError
    }
    result, err := evaluateBinaryOp(op, lhs.(float64), rhs.(float64))
    if err != nil {
        start := ctx.GetStart()
        line := start.GetLine()
        col := start.GetColumn()
        return NewRuntimeError(err.Error(), line, col)
    }
    return result
}
