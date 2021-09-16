package printerscript

import (
    "fmt"
    "github.com/antlr/antlr4/runtime/Go/antlr"
    "strconv"
)

type VisitorOptions struct {
    MaxLoopIterations int
    MaxOutputSize int
    EOL string
    Locals map[string]float64
}

type Visitor struct {
    antlr.ParseTreeVisitor
    maxOutputSize int
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
        maxOutputSize:    opts.MaxOutputSize,
        maxLoopIters:     opts.MaxLoopIterations,
        EOL:              opts.EOL,
        locals:           opts.Locals,
    }
}

func (v *Visitor) Visit(tree antlr.ParseTree) interface{} {
    return tree.Accept(v)
   //
   //switch val := tree.(type) {
   //case *StatementsContext:
   //    return val.Accept(v)
   //}
   //fmt.Println("couldn't figure out root of tree... :S")
   //return nil
}

func (v *Visitor) VisitSequence(ctx *SequenceContext) interface{} {
    if DEBUG { fmt.Println("VisitSequence") }
    return ctx.Statements().Accept(v)
}

func (v *Visitor) VisitStatements(ctx *StatementsContext) interface{} {
    if DEBUG { fmt.Println("VisitStatements") }
    if statement := ctx.Statement(); statement != nil {
        statement.Accept(v)
    }
    if statements := ctx.Statements(); statements != nil {
        ctx.Statements().Accept(v)
    }
    return nil
}

func (v *Visitor) VisitStatement(ctx *StatementContext) interface{} {
    if DEBUG { fmt.Println("VisitStatement") }
    assignment := ctx.Assignment()
    if assignment != nil {
        v.Visit(assignment)
        return nil
    }
    v.Visit(ctx.GCode())
    return nil
}

func (v *Visitor) VisitIfBlock(ctx *IfBlockContext) interface{} {
    if DEBUG { fmt.Println("VisitIfBlock") }
    condition := ctx.Expression()
    enterIf := v.Visit(condition)
    if enterIf != 0 {
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
    return v.VisitChildren(ctx)
}

func (v *Visitor) VisitWhileBlock(ctx *WhileBlockContext) interface{} {
    if DEBUG { fmt.Println("VisitWhileBlock") }
    condition := ctx.Expression()
    enterWhile := v.Visit(condition)
    for enterWhile != 0 {
        if v.loopIters > v.maxLoopIters {
            //start := ctx.GetStart()
            //line := start.GetLine()
            //col := start.GetColumn()
            // todo: runtime error
        }
        v.Visit(ctx.Statements())
        enterWhile = v.Visit(condition)
        v.loopIters++
    }
    return nil
}

func (v *Visitor) VisitAssignment(ctx *AssignmentContext) interface{} {
    if DEBUG { fmt.Println("VisitAssignment") }
    identifier := ctx.IDENTIFIER().GetText()
    value := v.Visit(ctx.Expression()).(float64)
    v.SetLocal(identifier, value)
    return nil
}

func (v *Visitor) VisitGCode(ctx *GCodeContext) interface{} {
    if DEBUG { fmt.Println("VisitGCode") }
    v.VisitChildren(ctx)
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
    value := v.Visit(ctx.Expression()).(float64)
    v.result += strconv.FormatFloat(value, 'f', -1, 64)
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
        // todo: runtime error
    }
    if fn == "max" || fn == "min" {
        // argc must be AT LEAST the required arity
        if argc < requiredArity {
            // todo: runtime error
        }
    } else {
        // argc must be EXACTLY the required arity
        if argc != requiredArity {
            // todo: runtime error
        }
    }

    // evaluate
    argv := make([]float64, 0, len(paramCtxs))
    for _, paramCtx := range paramCtxs {
        argv = append(argv, v.Visit(paramCtx).(float64))
    }
    retval, err := evaluateFunction(fn, argv)
    if err != nil {
        // todo: runtime error
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
    value, err := strconv.ParseInt(ctx.INT().GetText(), 10, 64)
    if err != nil {
        // todo: runtime error
    }
    return float64(value)
}

func (v *Visitor) VisitFloatExpr(ctx *FloatExprContext) interface{} {
    if DEBUG { fmt.Println("VisitFloatExpr") }
    value, err := strconv.ParseFloat(ctx.FLOAT().GetText(), 64)
    if err != nil {
        // todo: runtime error
    }
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
    operand := v.Visit(ctx.Expression()).(float64)
    result, err := evaluateUnaryOp(op, operand)
    if err != nil {
        // todo: runtime error
    }
    return result
}

func (v *Visitor) VisitBinaryOpExpr(ctx *BinaryOpExprContext) interface{} {
    if DEBUG { fmt.Println("VisitBinaryOpExpr") }
    op := ctx.GetChild(1).GetPayload().(antlr.Token).GetText()
    lhs := v.Visit(ctx.Expression(0)).(float64)
    rhs := v.Visit(ctx.Expression(1)).(float64)
    result, err := evaluateBinaryOp(op, lhs, rhs)
    if err != nil {
        // todo: runtime error
    }
    return result
}
