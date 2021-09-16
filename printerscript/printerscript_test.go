package printerscript

import (
    "fmt"
    "math"
    "testing"
)

func expectOutput(t *testing.T, input, expectedOutput string) {
    result, err := EvaluateString(input)
    if err != nil {
        t.Error(err)
        return
    }
    if result.Output != expectedOutput {
        if DEBUG {
            fmt.Println("< EXPECTED >")
            fmt.Println(expectedOutput)
            fmt.Println("< /EXPECTED >")
            fmt.Println("< RESULT >")
            fmt.Println(result.Output)
            fmt.Println("< /RESULT >")
        }
        t.Error("mismatch in output")
    }
}

func expectLocal(t *testing.T, input, local string, expectedValue float64) {
    result, err := EvaluateString(input)
    if err != nil {
        t.Error(err)
        return
    }
    if result.Locals[local] != expectedValue {
        if DEBUG {
            fmt.Printf("EXPECTED: locals[%s] = %f\n", local, expectedValue)
            fmt.Printf("RESULT: locals[%s] = %f\n", local, result.Locals[local])
        }
        t.Error("mismatch in local value")
    }
}

func expectLocalClose(t *testing.T, input, local string, expectedValue float64) {
    result, err := EvaluateString(input)
    if err != nil {
        t.Error(err)
        return
    }
    if math.Abs(result.Locals[local] - expectedValue) > 10e-5 {
        if DEBUG {
            fmt.Printf("EXPECTED: locals[%s] = %f\n", local, expectedValue)
            fmt.Printf("RESULT: locals[%s] = %f\n", local, result.Locals[local])
        }
        t.Error("mismatch in local value")
    }
}

func expectSyntaxError(t *testing.T, input string) {
    _, err := EvaluateString(input)
    if err != nil {
        // todo: ensure it's actually a syntax error
        return
    }
    t.Error("expected syntax error")
}

func expectRuntimeError(t *testing.T, input string) {
    _, err := EvaluateString(input)
    if err != nil {
        // todo: ensure it's actually a runtime error
        return
    }
    t.Error("expected runtime error")
}

// 1 - G-code statements, comments
func Test_LiteralOutput(t *testing.T) {
    input := `// code-only comment!
'T0'
/* block comment */
/* multi-line
   block comment */
"M82 // not a comment" // yes a comment
"G28 ; home all axes"
"M104 S{{TEMP}} ; heat extruder"
"M140 S{{BED}} ; heat bed"
"M109 S{{TEMP}} ; wait for extruder to heat"
"M190 S{{BED}} ; wait for bed to heat"
"G92 E0 ; reset extruder"
`
    expected := `T0
M82 // not a comment
G28 ; home all axes
M104 S0 ; heat extruder
M140 S0 ; heat bed
M109 S0 ; wait for extruder to heat
M190 S0 ; wait for bed to heat
G92 E0 ; reset extruder
`
    expectOutput(t, input, expected)
}
func Test_Escapes(t *testing.T) {
    expectOutput(t, "\"\\\"\"", "\"\n")
    expectOutput(t, "'\\''", "'\n")
    expectOutput(t, "\"\\{{\"", "{{\n")
    expectOutput(t, "\"\\{{foo}}\"", "{{foo}}\n")
    expectOutput(t, "\"\\{{\"", "{{\n")
    expectOutput(t, "\"{foo}\"", "{foo}\n")
    expectOutput(t, "\"}\"", "}\n")
}
func Test_LineComments(t *testing.T) {
    expectOutput(t, "//comment\n'G1'\n//comment", "G1\n")
    expectOutput(t, "'G1 ;A' //B", "G1 ;A\n")
    expectOutput(t, "'G1 //A ;B'", "G1 //A ;B\n")
}
func Test_BlockComments(t *testing.T) {
    expectOutput(t, "/*one-liner*/\n'G1'/*multi-\nliner*/", "G1\n")
    expectOutput(t, "/*unclosed \n'G1' /*continuation", "")
    expectOutput(t, "'G1 ;A' /*B*/", "G1 ;A\n")
    expectOutput(t, "'G1 /*A*/ ;B'", "G1 /*A*/ ;B\n")
}

// 2 - Identifiers, assignment statements, expressions and precedence
func Test_IdentifiersAssignmentAndExpressions(t *testing.T) {
    input := "x = a + b * -(c - d) / (e % f)\n'; x = {{x}}'\n"
    locals := map[string]float64{
        "a": 1.0,
        "b": 5.0,
        "c": 10.0,
        "d": 2.0,
        "e": 11.0,
        "f": 3.0,
    }
    result, err := EvaluateStringAndLocals(input, locals)
    if err != nil {
        t.Error(err)
        return
    }
    if result.Locals["x"] != -19 {
        t.Errorf("expected locals[%s] = %f, got %f", "x", -19.0, result.Locals["x"])
    }
    if result.Locals["a"] != locals["a"] {
        t.Errorf("expected locals[%s] = %f, got %f", "a", locals["a"], result.Locals["a"])
    }
    if result.Locals["b"] != locals["b"] {
        t.Errorf("expected locals[%s] = %f, got %f", "b", locals["b"], result.Locals["b"])
    }
    if result.Locals["c"] != locals["c"] {
        t.Errorf("expected locals[%s] = %f, got %f", "c", locals["c"], result.Locals["c"])
    }
    if result.Locals["d"] != locals["d"] {
        t.Errorf("expected locals[%s] = %f, got %f", "d", locals["d"], result.Locals["d"])
    }
    if result.Locals["e"] != locals["e"] {
        t.Errorf("expected locals[%s] = %f, got %f", "e", locals["e"], result.Locals["e"])
    }
    if result.Locals["f"] != locals["f"] {
        t.Errorf("expected locals[%s] = %f, got %f", "f", locals["f"], result.Locals["f"])
    }
    if result.Output != "; x = -19\n" {
        t.Error("mismatch in output")
    }
}

// 3 - Operators
func Test_Parens(t *testing.T) {
    expectLocal(t, "a = (((((5)))))", "a", 5.0)
    expectSyntaxError(t, "a = (((((5))))")
    expectSyntaxError(t, "a = ((((5)))))")
}
func Test_Not(t *testing.T) {
    // todo
}
func Test_UnaryMinus(t *testing.T) {
    // todo
}
func Test_UndefinedUnaryOp(t *testing.T) {
    // todo
}
func Test_Mult(t *testing.T) {
    // todo
}
func Test_Div(t *testing.T) {
    // todo
}
func Test_Mod(t *testing.T) {
    // todo
}
func Test_Plus(t *testing.T) {
    // todo
}
func Test_Minus(t *testing.T) {
    // todo
}
func Test_Equal(t *testing.T) {
    // todo
}
func Test_NotEqual(t *testing.T) {
    // todo
}
func Test_LessThan(t *testing.T) {
    // todo
}
func Test_LessThanEqual(t *testing.T) {
    // todo
}
func Test_GreaterThan(t *testing.T) {
    // todo
}
func Test_GreaterThanEqual(t *testing.T) {
    // todo
}
func Test_Or(t *testing.T) {
    // todo
}
func Test_And(t *testing.T) {
    // todo
}
func Test_UndefinedBinaryOp(t *testing.T) {
    // todo
}

// 4 - Function calls
func Test_Min(t *testing.T) {
    // todo
}
func Test_Max(t *testing.T) {
    // todo
}
func Test_Abs(t *testing.T) {
    // todo
}
func Test_Floor(t *testing.T) {
    // todo
}
func Test_Ceil(t *testing.T) {
    // todo
}
func Test_Round(t *testing.T) {
    // todo
}
func Test_Trunc(t *testing.T) {
    // todo
}
func Test_Pow(t *testing.T) {
    // todo
}
func Test_Sqrt(t *testing.T) {
    // todo
}
func Test_Sin(t *testing.T) {
    // todo
}
func Test_Cos(t *testing.T) {
    // todo
}
func Test_Tan(t *testing.T) {
    // todo
}
func Test_Asin(t *testing.T) {
    // todo
}
func Test_Acos(t *testing.T) {
    // todo
}
func Test_Atan(t *testing.T) {
    // todo
}
func Test_UndefinedFunction(t *testing.T) {
    // todo
}

// 5 - Branching
func Test_If(t *testing.T) {
    // todo
}
func Test_Else(t *testing.T) {
    // todo
}
func Test_ElseIf(t *testing.T) {
    // todo
}
func Test_MultiElseIf(t *testing.T) {
    // todo
}
func Test_NestedIf(t *testing.T) {
    // todo
}

// 6 - Looping
func Test_FiniteLoop(t *testing.T) {
    // todo
}
func Test_NestedLoop(t *testing.T) {
    // todo
}
func Test_InfiniteLoop(t *testing.T) {
    // todo
}
