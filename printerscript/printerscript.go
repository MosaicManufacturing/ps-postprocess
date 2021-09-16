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

func Validate(input string) error {
    // todo: run lexer and parser, and just make sure
    //   no syntax errors were found (i.e. no visitor)
    return nil
}

func Evaluate(opts InterpreterOptions) (InterpreterResult, error) {
    input := opts.Input // todo: trim, convert newlines to only be \n, and add trailing \n

    // lexer
    if DEBUG { fmt.Println("===== LEXER =====") }
    istream := antlr.NewInputStream(input)
    lexer := NewSequenceLexer(istream)
    lexerErrorListener := antlr.NewConsoleErrorListener() // todo: implement our own
    lexer.RemoveErrorListeners()
    lexer.AddErrorListener(lexerErrorListener)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    tokens.Fill()
    if DEBUG {
        for _, token := range tokens.GetAllTokens() {
            fmt.Print(token.GetText())
        }
        fmt.Println()
    }

    // parser
    if DEBUG { fmt.Println("===== PARSER =====") }
    parser := NewSequenceParser(tokens)
    parserErrorListener := antlr.NewConsoleErrorListener() // todo: implement our own
    // todo: no equivalent of parser.removeParseListeners() ?
    parser.RemoveErrorListeners()
    parser.AddErrorListener(parserErrorListener)
    tree := parser.Sequence()

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
    visitor.Visit(tree)

    // result
    if DEBUG { fmt.Println("===== RESULT =====") }
    result := InterpreterResult{}
    result.Output = visitor.GetResult()
    if !opts.TrailingNewline && len(result.Output) > 0 {
        result.Output = result.Output[:len(result.Output)-1]
    }
    result.Locals = visitor.GetLocals()
    return result, nil // todo: actually return errors, don't just print them
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