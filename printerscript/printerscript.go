package printerscript

import (
    "fmt"
    "github.com/antlr/antlr4/runtime/Go/antlr"
)

const DEBUG = false

type InterpreterOptions struct {
    MaxLoopIterations int
    MaxOutputSize int
    EOL string
    TrailingNewline bool
    Input string
    Locals map[string]float64
}

type InterpreterResult struct {
    Output string
    Locals map[string]float64
}

func normalizeInput(input string) string {
    // todo: trim, convert newlines to only be \n, and add trailing \n
    return input
}

func Validate(input string) error {
    input = normalizeInput(input)

    // lexer
    if DEBUG { fmt.Println("===== LEXER =====") }
    istream := antlr.NewInputStream(input)
    lexer := NewSequenceLexer(istream)
    lexerErrorListener := NewSyntaxErrorListener()
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(lexerErrorListener)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    tokens.Fill()
    if err := lexerErrorListener.GetError(); err != nil {
        return err
    }
    if DEBUG {
        for _, token := range tokens.GetAllTokens() {
            fmt.Print(token.GetText())
        }
        fmt.Println()
    }

    // parser
    if DEBUG { fmt.Println("===== PARSER =====") }
    parser := NewSequenceParser(tokens)
    parserErrorListener := NewSyntaxErrorListener()
    parser.RemoveErrorListeners()
    parser.AddErrorListener(parserErrorListener)
    parser.Sequence()
    if err := parserErrorListener.GetError(); err != nil {
        return err
    }
    return nil
}

func Evaluate(opts InterpreterOptions) (InterpreterResult, error) {
    input := normalizeInput(opts.Input)
    result := InterpreterResult{}

    // lexer
    if DEBUG { fmt.Println("===== LEXER =====") }
    istream := antlr.NewInputStream(input)
    lexer := NewSequenceLexer(istream)
    lexerErrorListener := NewSyntaxErrorListener()
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(lexerErrorListener)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    tokens.Fill()
    if err := lexerErrorListener.GetError(); err != nil {
        return result, err
    }
    if DEBUG {
        for _, token := range tokens.GetAllTokens() {
            fmt.Print(token.GetText())
        }
        fmt.Println()
    }

    // parser
    if DEBUG { fmt.Println("===== PARSER =====") }
    parser := NewSequenceParser(tokens)
    parserErrorListener := NewSyntaxErrorListener()
    parser.RemoveErrorListeners()
    parser.AddErrorListener(parserErrorListener)
    tree := parser.Sequence()
    if err := parserErrorListener.GetError(); err != nil {
        return result, err
    }

    // visitor
    if DEBUG { fmt.Println("===== VISITOR =====") }
    visitorOpts := VisitorOptions{
        MaxLoopIterations: opts.MaxLoopIterations,
        MaxOutputSize:     opts.MaxOutputSize,
        EOL:               opts.EOL,
        Locals:            make(map[string]float64),
    }
    if opts.Locals != nil {
        for k, v := range opts.Locals {
            visitorOpts.Locals[k] = v
            if DEBUG { fmt.Printf("%s: %f\n", k, v) }
        }
    }
    visitor := NewVisitor(visitorOpts)
    err := visitor.Visit(tree)
    if err != nil {
        if runtimeErr, ok := err.(*RuntimeError); ok {
            return result, runtimeErr
        }
    }

    // result
    if DEBUG { fmt.Println("===== RESULT =====") }
    result.Output = visitor.GetResult()
    if !opts.TrailingNewline && len(result.Output) > 0 {
        result.Output = result.Output[:len(result.Output)-1]
    }
    result.Locals = visitor.GetLocals()
    return result, nil
}

func EvaluateStringAndLocals(input string, locals map[string]float64) (InterpreterResult, error) {
    opts := InterpreterOptions{
        MaxLoopIterations: 1e6,
        MaxOutputSize:     50 * 1024 * 1024, // 50 MiB
        EOL:               "\n",
        TrailingNewline:   true,
        Input:             input,
        Locals:            locals,
    }
    return Evaluate(opts)
}

func EvaluateString(input string) (InterpreterResult, error) {
    return EvaluateStringAndLocals(input, nil)
}