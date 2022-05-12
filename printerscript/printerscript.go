package printerscript

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"mosaicmfg.com/ps-postprocess/printerscript/interpreter"
)

type TokenStream antlr.CommonTokenStream

type Tree interpreter.ISequenceContext

type InterpreterOptions struct {
	MaxLoopIterations int                // raise a RuntimeError if we exceed this many iterations (default: 100,000)
	EOL               string             // line ending to use in output (default: "\n")
	TrailingNewline   bool               // include a trailing newline in output (default: true)
	Locals            map[string]float64 // initial local values (default: empty map)
}

type InterpreterResult struct {
	Output string
	Locals map[string]float64
}

// Lex takes a PrinterScript script as input and runs the PrinterScript lexer,
// returning a token stream which can be passed to Parse, along with an error
// if lexing failed. In the case of a syntax error, the error will be a
// SyntaxError including line and column information.
func Lex(input string) (*TokenStream, error) {
	input = normalizeInput(input)
	if interpreter.DEBUG {
		fmt.Println("===== LEXER =====")
	}
	istream := antlr.NewInputStream(input)
	lexer := interpreter.NewSequenceLexer(istream)
	lexerErrorListener := interpreter.NewSyntaxErrorListener()
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrorListener)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	tokens.Fill()
	if err := lexerErrorListener.GetError(); err != nil {
		return nil, err
	}
	if interpreter.DEBUG {
		for _, token := range tokens.GetAllTokens() {
			fmt.Print(token.GetText())
		}
		fmt.Println()
	}
	return (*TokenStream)(tokens), nil
}

// Parse takes a token stream produced by Lex as input and runs the PrinterScript
// parser, returning a parse tree which can be passed to EvaluateTree, along with
// an error if parsing failed. In the case of a syntax error, the error will be
// a SyntaxError including line and column information.
func Parse(tokens *TokenStream) (Tree, error) {
	if interpreter.DEBUG {
		fmt.Println("===== PARSER =====")
	}
	parser := interpreter.NewSequenceParser((*antlr.CommonTokenStream)(tokens))
	parserErrorListener := interpreter.NewSyntaxErrorListener()
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserErrorListener)
	tree := parser.Sequence()
	if err := parserErrorListener.GetError(); err != nil {
		return nil, err
	}
	return tree, nil
}

// LexAndParse is a convenience function that calls Lex and then Parse, checking
// for lexer errors before parsing. The returned parse tree can be passed to
// EvaluateTree. This function is especially useful when scripts may be
// evaluated more than once, as the parse tree can be generated just once and
// then re-used with each evaluation.
func LexAndParse(input string) (Tree, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}
	return Parse(tokens)
}

// Validate takes a PrinterScript script and checks it for syntax errors,
// running Lex and then Parse but only returning the error (if any).
func Validate(input string) error {
	_, err := LexAndParse(input)
	return err
}

// EvaluateTree takes a parse tree produced by Parse, as well as a set of
// runtime options for the interpreter, and evaluates the tree. It returns
// a result containing the output string and the final value of all locals,
// as well as a RuntimeError if evaluation failed.
func EvaluateTree(tree Tree, opts InterpreterOptions) (InterpreterResult, error) {
	result := InterpreterResult{}

	// visitor
	if interpreter.DEBUG {
		fmt.Println("===== VISITOR =====")
	}
	visitorOpts := interpreter.VisitorOptions{
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
			if interpreter.DEBUG {
				fmt.Printf("%s: %f\n", k, v)
			}
		}
	}
	visitor := interpreter.NewVisitor(visitorOpts)
	if err := visitor.Visit(tree); err != nil {
		if runtimeErr, ok := err.(*interpreter.RuntimeError); ok {
			return result, runtimeErr
		}
	}

	// result
	if interpreter.DEBUG {
		fmt.Println("===== RESULT =====")
	}
	result.Output = visitor.GetResult()
	if !opts.TrailingNewline {
		// trailing newlines are always generated -- remove now if disabled
		outputLen := len(result.Output)
		eolLen := len(opts.EOL)
		if outputLen >= eolLen {
			result.Output = result.Output[:outputLen-eolLen]
		}
	}
	result.Locals = visitor.GetLocals()
	return result, nil
}

// EvaluateWithOpts takes a PrinterScript script and a set of runtime options for
// the interpreter, and evaluates the script. It returns a result containing the
// output string and the final value of all locals, as well as a RuntimeError
// if evaluation failed.
func EvaluateWithOpts(input string, opts InterpreterOptions) (InterpreterResult, error) {
	tree, err := LexAndParse(input)
	if err != nil {
		return InterpreterResult{}, err
	}
	return EvaluateTree(tree, opts)
}

// EvaluateWithLocals takes a PrinterScript script and a set of initial local values,
// and evaluates the script. It returns a result containing the output string and the
// final value of all locals, as well as a RuntimeError if evaluation failed.
func EvaluateWithLocals(input string, locals map[string]float64) (InterpreterResult, error) {
	opts := InterpreterOptions{
		TrailingNewline: true,
		Locals:          locals,
	}
	return EvaluateWithOpts(input, opts)
}

// Evaluate takes a PrinterScript script and evaluates it. It returns a result
// containing the output string and the final value of all locals, as well as a
// RuntimeError if evaluation failed.
func Evaluate(input string) (InterpreterResult, error) {
	return EvaluateWithLocals(input, nil)
}
