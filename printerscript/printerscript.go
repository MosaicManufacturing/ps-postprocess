package printerscript

import (
    "fmt"
    "github.com/antlr/antlr4/runtime/Go/antlr"
)

const DEBUG = false

type InterpreterOptions struct {
    MaxLoopIterations int
    EOL string
    TrailingNewline bool
    Locals map[string]float64
}

type InterpreterResult struct {
    Output string
    Locals map[string]float64
}

func Lex(input string) (*antlr.CommonTokenStream, error) {
    input = normalizeInput(input)
    if DEBUG { fmt.Println("===== LEXER =====") }
    istream := antlr.NewInputStream(input)
    lexer := NewSequenceLexer(istream)
    lexerErrorListener := NewSyntaxErrorListener()
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(lexerErrorListener)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    tokens.Fill()
    if err := lexerErrorListener.GetError(); err != nil {
        return nil, err
    }
    if DEBUG {
        for _, token := range tokens.GetAllTokens() {
            fmt.Print(token.GetText())
        }
        fmt.Println()
    }
    return tokens, nil
}

func Parse(tokens *antlr.CommonTokenStream) (ISequenceContext, error) {
    if DEBUG { fmt.Println("===== PARSER =====") }
    parser := NewSequenceParser(tokens)
    parserErrorListener := NewSyntaxErrorListener()
    parser.RemoveErrorListeners()
    parser.AddErrorListener(parserErrorListener)
    tree := parser.Sequence()
    if err := parserErrorListener.GetError(); err != nil {
        return nil, err
    }
    return tree, nil
}

func LexAndParse(input string) (ISequenceContext, error) {
    tokens, err := Lex(input)
    if err != nil {
        return nil, err
    }
    return Parse(tokens)
}

func Validate(input string) error {
    _, err := LexAndParse(input)
    return err
}

func EvaluateTree(tree ISequenceContext, opts InterpreterOptions) (InterpreterResult, error) {
    result := InterpreterResult{}

    // visitor
    if DEBUG { fmt.Println("===== VISITOR =====") }
    visitorOpts := VisitorOptions{
        MaxLoopIterations: 1e6, // 100k total iterations, including nesting
        EOL:               "\n",
        Locals:            make(map[string]float64),
    }
    if opts.MaxLoopIterations > 0 {
        visitorOpts.MaxLoopIterations = opts.MaxLoopIterations
    }
    if opts.EOL != "" {
        visitorOpts.EOL = opts.EOL
    }
    if opts.Locals != nil {
        for k, v := range opts.Locals {
            visitorOpts.Locals[k] = v
            if DEBUG { fmt.Printf("%s: %f\n", k, v) }
        }
    }
    visitor := NewVisitor(visitorOpts)
    if err := visitor.Visit(tree); err != nil {
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

func EvaluateWithOpts(input string, opts InterpreterOptions) (InterpreterResult, error) {
    tree, err := LexAndParse(input)
    if err != nil {
        return InterpreterResult{}, err
    }
    return EvaluateTree(tree, opts)
}

func EvaluateWithLocals(input string, locals map[string]float64) (InterpreterResult, error) {
    opts := InterpreterOptions{
        TrailingNewline: true,
        Locals:          locals,
    }
    return EvaluateWithOpts(input, opts)
}

func Evaluate(input string) (InterpreterResult, error) {
    return EvaluateWithLocals(input, nil)
}