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
        fmt.Println("< EXPECTED >")
        fmt.Println(expectedOutput)
        fmt.Println("< /EXPECTED >")
        fmt.Println("< RESULT >")
        fmt.Println(result.Output)
        fmt.Println("< /RESULT >")
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
        fmt.Printf("EXPECTED: locals[%s] = %f\n", local, expectedValue)
        fmt.Printf("RESULT: locals[%s] = %f\n", local, result.Locals[local])
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

func expectSuccess(t *testing.T, input string) {
    if _, err := EvaluateString(input); err != nil {
        t.Error(err)
    }
}

func expectSyntaxError(t *testing.T, input string) {
    _, err := EvaluateString(input)
    if err != nil {
        if _, ok := err.(*SyntaxError); ok {
            // got what we expected
            return
        }
    }
    t.Error("expected syntax error")
}

func expectRuntimeError(t *testing.T, input string) {
    _, err := EvaluateString(input)
    if err != nil {
        if _, ok := err.(*RuntimeError); ok {
            // got what we expected
            return
        }
    }
    t.Error("expected runtime error")
}

//
// 1 - G-code statements, comments
//

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

//
// 2 - Identifiers, assignment statements, expressions and precedence
//

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

//
// 3 - Operators
//

func Test_Parens(t *testing.T) {
    expectLocal(t, "a = (((((5)))))", "a", 5.0)
    expectSyntaxError(t, "a = (((((5))))")
    expectSyntaxError(t, "a = ((((5)))))")
}

func Test_Not(t *testing.T) {
    expectSyntaxError(t, "a = !")
    expectSyntaxError(t, "a = y ! x")
    expectSuccess(t, "a = ! x")
    expectLocal(t, "x=1\na = !x", "a", 0.0)
    expectLocal(t, "a = !z", "a", 1.0)
}

func Test_UnaryMinus(t *testing.T) {
    expectLocal(t, "x=1\na = -x", "a", -1.0)
    expectLocal(t, "x=1\na = --x", "a", 1.0)
    expectLocal(t, "a = - -1", "a", 1.0)
    expectLocal(t, "a = -0", "a", 0.0)
}

func Test_UndefinedUnaryOp(t *testing.T) {
    if _, err := evaluateUnaryOp("foo", 0); err == nil {
        t.Error("expected error for undefined unary operator")
    }
}

func Test_Mult(t *testing.T) {
    expectSyntaxError(t, "a = *")
    expectSyntaxError(t, "a = * x")
    expectSuccess(t, "a = x * y")
    expectLocal(t, "x=1\ny=1\n  a = x * y", "a", 1.0)
    expectLocal(t, "x=1\ny=1\n  a = 0 * y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 * y", "a", 2.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 * 0", "a", 0.0)
}

func Test_Div(t *testing.T) {
    expectSyntaxError(t, "a = /")
    expectSyntaxError(t, "a = / x")
    expectSuccess(t, "a = x / 1")
    expectLocal(t, "x=1\ny=1\n  a = x / y", "a", 1.0)
    expectLocal(t, "x=1\ny=1\n  a = 0 / y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 / y", "a", 2.0)
    expectLocalClose(t, "x=1\ny=1\n  a = x / 2", "a", 0.5)
    expectRuntimeError(t, "a = 1 / 0")
}

func Test_Mod(t *testing.T) {
    expectSyntaxError(t, "a = %")
    expectSyntaxError(t, "a = % x")
    expectSuccess(t, "a = x % 1")
    expectLocal(t, "x=1\ny=1\n  a = x % y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = 0 % y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = 10 % 3", "a", 1.0)
    expectLocal(t, "x=1\ny=1\n  a = 11 % 3", "a", 2.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 % y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = 3 % 2", "a", 1.0)
    expectLocal(t, "x=1\ny=1\n  a = -5 % 4", "a", 3.0)
    expectLocal(t, "x=1\ny=1\n  a = -7 % 4", "a", 1.0)
    expectRuntimeError(t, "a = 1 % 0")
}

func Test_Plus(t *testing.T) {
    expectSyntaxError(t, "a = +")
    expectSyntaxError(t, "a = + x")
    expectSuccess(t, "a = x + y")
    expectLocal(t, "x=1\ny=1\n  a = x + y", "a", 2.0)
    expectLocal(t, "x=1\ny=1\n  a = x + 3", "a", 4.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 + y", "a", 3.0)
    expectLocal(t, "x=1\ny=1\n  a = -2 + 1", "a", -1.0)
    expectLocal(t, "x=1\ny=1\n  a = 0 + 0", "a", 0.0)
}

func Test_Minus(t *testing.T) {
    expectSyntaxError(t, "a = -")
    expectSuccess(t, "a = x - y")
    expectLocal(t, "x=1\ny=1\n  a = x - y", "a", 0.0)
    expectLocal(t, "x=1\ny=1\n  a = x - 3", "a", -2.0)
    expectLocal(t, "x=1\ny=1\n  a = 2 - y", "a", 1.0)
    expectLocal(t, "x=1\ny=1\n  a = -2 - 1", "a", -3.0)
    expectLocal(t, "x=1\ny=1\n  a = 0 - 0", "a", 0.0)
}

func Test_Equal(t *testing.T) {
    expectSyntaxError(t, "a = ==")
    expectSyntaxError(t, "a = == x")
    expectSuccess(t, "a = x == y")
    expectLocal(t, "a = x == y", "a", 1.0)
    expectLocal(t, "a = x == 1", "a", 0.0)
    expectLocal(t, "a = x == true", "a", 0.0)
    expectLocal(t, "a = x == false", "a", 1.0)
    expectLocal(t, "a = true == true", "a", 1.0)
    expectLocal(t, "a = true == false", "a", 0.0)
}

func Test_NotEqual(t *testing.T) {
    expectSyntaxError(t, "a = !=")
    expectSyntaxError(t, "a = != x")
    expectSuccess(t, "a = x != y")
    expectLocal(t, "a = x != y", "a", 0.0)
    expectLocal(t, "a = x != 0", "a", 0.0)
    expectLocal(t, "a = x != 1", "a", 1.0)
    expectLocal(t, "a = x != true", "a", 1.0)
    expectLocal(t, "a = x != false", "a", 0.0)
    expectLocal(t, "a = true != true", "a", 0.0)
    expectLocal(t, "a = true != false", "a", 1.0)
}

func Test_LessThan(t *testing.T) {
    expectSyntaxError(t, "a = <")
    expectSyntaxError(t, "a = < x")
    expectSuccess(t, "a = x < y")
    expectLocal(t, "a = 0 < 1", "a", 1.0)
    expectLocal(t, "a = 0 < 0", "a", 0.0)
    expectLocal(t, "a = 1 < 0", "a", 0.0)
    expectLocal(t, "a = -1 < 0", "a", 1.0)
    expectLocal(t, "a = -1 < 1", "a", 1.0)
    expectLocal(t, "a = -2 < 1", "a", 1.0)
    expectLocal(t, "a = 2 < -1", "a", 0.0)
}

func Test_LessThanEqual(t *testing.T) {
    expectSyntaxError(t, "a = <=")
    expectSyntaxError(t, "a = <= x")
    expectSuccess(t, "a = x <= y")
    expectLocal(t, "a = 0 <= 1", "a", 1.0)
    expectLocal(t, "a = 0 <= 0", "a", 1.0)
    expectLocal(t, "a = 1 <= 0", "a", 0.0)
    expectLocal(t, "a = -1 <= 0", "a", 1.0)
    expectLocal(t, "a = -1 <= 1", "a", 1.0)
    expectLocal(t, "a = -2 <= 1", "a", 1.0)
    expectLocal(t, "a = 2 <= -1", "a", 0.0)
}

func Test_GreaterThan(t *testing.T) {
    expectSyntaxError(t, "a = >")
    expectSyntaxError(t, "a = > x")
    expectSuccess(t, "a = x > y")
    expectLocal(t, "a = 1 > 0", "a", 1.0)
    expectLocal(t, "a = 0 > 0", "a", 0.0)
    expectLocal(t, "a = 0 > 1", "a", 0.0)
    expectLocal(t, "a = 0 > -1", "a", 1.0)
    expectLocal(t, "a = 1 > -1", "a", 1.0)
    expectLocal(t, "a = 1 > -2", "a", 1.0)
    expectLocal(t, "a = -1 > 2", "a", 0.0)
}

func Test_GreaterThanEqual(t *testing.T) {
    expectSyntaxError(t, "a = >=")
    expectSyntaxError(t, "a = >= x")
    expectSuccess(t, "a = x >= y")
    expectLocal(t, "a = 1 >= 0", "a", 1.0)
    expectLocal(t, "a = 0 >= 0", "a", 1.0)
    expectLocal(t, "a = 0 >= 1", "a", 0.0)
    expectLocal(t, "a = 0 >= -1", "a", 1.0)
    expectLocal(t, "a = 1 >= -1", "a", 1.0)
    expectLocal(t, "a = 1 >= -2", "a", 1.0)
    expectLocal(t, "a = -1 >= 2", "a", 0.0)
}

func Test_Or(t *testing.T) {
    expectSyntaxError(t, "a = ||")
    expectSyntaxError(t, "a = || x")
    expectSuccess(t, "a = x || y")
    expectLocal(t, "a = 0 || 0", "a", 0.0)
    expectLocal(t, "a = 1 || 0", "a", 1.0)
    expectLocal(t, "a = 0 || 1", "a", 1.0)
    expectLocal(t, "a = 1 || 1", "a", 1.0)
}

func Test_And(t *testing.T) {
    expectSyntaxError(t, "a = &&")
    expectSyntaxError(t, "a = && x")
    expectSuccess(t, "a = x && y")
    expectLocal(t, "a = 0 && 0", "a", 0.0)
    expectLocal(t, "a = 1 && 0", "a", 0.0)
    expectLocal(t, "a = 0 && 1", "a", 0.0)
    expectLocal(t, "a = 1 && 1", "a", 1.0)
}

func Test_UndefinedBinaryOp(t *testing.T) {
    if _, err := evaluateBinaryOp("foo", 0, 0); err == nil {
        t.Error("expected error for undefined binary operator")
    }
}

//
// 4 - Function calls
//

func Test_Min(t *testing.T) {
    expectRuntimeError(t, "a = min()")
    expectRuntimeError(t, "a = min(b)")
    expectSuccess(t, "a = min(b, c)")
    expectSuccess(t, "a = min(b, c, d)")
    expectSuccess(t, "a = min(b, c, d, e, f, g)")
    expectLocal(t, "a = min(1, 2)", "a", 1.0)
    expectLocal(t, "a = min(2, 1)", "a", 1.0)
    expectLocal(t, "a = min(1, -2)", "a", -2.0)
    expectLocal(t, "a = min(-2, 1)", "a", -2.0)
    expectLocal(t, "a = min(-3.5, -2)", "a", -3.5)
    expectLocal(t, "a = min(-2, -3.5)", "a", -3.5)
    expectLocal(t, "a = min(-3.5, -2.5)", "a", -3.5)
    expectLocal(t, "a = min(-2.5, -3.5)", "a", -3.5)
    expectLocal(t, "a = min(0, 1, 2)", "a", 0.0)
    expectLocal(t, "a = min(0, 1, 2, 3)", "a", 0.0)
}

func Test_Max(t *testing.T) {
    expectRuntimeError(t, "a = max()")
    expectRuntimeError(t, "a = max(b)")
    expectSuccess(t, "a = max(b, c)")
    expectSuccess(t, "a = max(b, c, d)")
    expectSuccess(t, "a = max(b, c, d, e, f, g)")
    expectLocal(t, "a = max(1, 2)", "a", 2.0)
    expectLocal(t, "a = max(2, 1)", "a", 2.0)
    expectLocal(t, "a = max(1, -2)", "a", 1.0)
    expectLocal(t, "a = max(-2, 1)", "a", 1.0)
    expectLocal(t, "a = max(-3.5, -2)", "a", -2.0)
    expectLocal(t, "a = max(-2, -3.5)", "a", -2.0)
    expectLocal(t, "a = max(-3.5, -2.5)", "a", -2.5)
    expectLocal(t, "a = max(-2.5, -3.5)", "a", -2.5)
    expectLocal(t, "a = max(0, 1, 2)", "a", 2.0)
    expectLocal(t, "a = max(0, 1, 2, 3)", "a", 3.0)
}

func Test_Abs(t *testing.T) {
    expectRuntimeError(t, "a = abs()")
    expectSuccess(t, "a = abs(b)")
    expectRuntimeError(t, "a = abs(b, c)")
    expectLocal(t, "a = abs(1)", "a", 1.0)
    expectLocal(t, "a = abs(-1)", "a", 1.0)
    expectLocal(t, "a = abs(1.5)", "a", 1.5)
    expectLocal(t, "a = abs(-1.5)", "a", 1.5)
    expectLocal(t, "a = abs(0)", "a", 0.0)
}

func Test_Floor(t *testing.T) {
    expectRuntimeError(t, "a = floor()")
    expectSuccess(t, "a = floor(b)")
    expectRuntimeError(t, "a = floor(b, c)")
    expectLocal(t, "a = floor(1)", "a", 1.0)
    expectLocal(t, "a = floor(0)", "a", 0.0)
    expectLocal(t, "a = floor(2.3)", "a", 2.0)
    expectLocal(t, "a = floor(3.8)", "a", 3.0)
    expectLocal(t, "a = floor(5.5)", "a", 5.0)
    expectLocal(t, "a = floor(-2.3)", "a", -3.0)
    expectLocal(t, "a = floor(-3.8)", "a", -4.0)
    expectLocal(t, "a = floor(-5.5)", "a", -6.0)
}

func Test_Ceil(t *testing.T) {
    expectRuntimeError(t, "a = ceil()")
    expectSuccess(t, "a = ceil(b)")
    expectRuntimeError(t, "a = ceil(b, c)")
    expectLocal(t, "a = ceil(1)", "a", 1.0)
    expectLocal(t, "a = ceil(0)", "a", 0.0)
    expectLocal(t, "a = ceil(2.3)", "a", 3.0)
    expectLocal(t, "a = ceil(3.8)", "a", 4.0)
    expectLocal(t, "a = ceil(5.5)", "a", 6.0)
    expectLocal(t, "a = ceil(-2.3)", "a", -2.0)
    expectLocal(t, "a = ceil(-3.8)", "a", -3.0)
    expectLocal(t, "a = ceil(-5.5)", "a", -5.0)
}

func Test_Round(t *testing.T) {
    expectRuntimeError(t, "a = round()")
    expectSuccess(t, "a = round(b)")
    expectRuntimeError(t, "a = round(b, c)")
    expectLocal(t, "a = round(1)", "a", 1.0)
    expectLocal(t, "a = round(0)", "a", 0.0)
    expectLocal(t, "a = round(2.3)", "a", 2.0)
    expectLocal(t, "a = round(3.8)", "a", 4.0)
    expectLocal(t, "a = round(5.5)", "a", 6.0)
    expectLocal(t, "a = round(-2.3)", "a", -2.0)
    expectLocal(t, "a = round(-3.8)", "a", -4.0)
    expectLocal(t, "a = round(-5.5)", "a", -6.0)
}

func Test_Trunc(t *testing.T) {
    expectRuntimeError(t, "a = trunc()")
    expectSuccess(t, "a = trunc(b)")
    expectRuntimeError(t, "a = trunc(b, c)")
    expectLocal(t, "a = trunc(1)", "a", 1.0)
    expectLocal(t, "a = trunc(0)", "a", 0.0)
    expectLocal(t, "a = trunc(2.3)", "a", 2.0)
    expectLocal(t, "a = trunc(3.8)", "a", 3.0)
    expectLocal(t, "a = trunc(5.5)", "a", 5.0)
    expectLocal(t, "a = trunc(-2.3)", "a", -2.0)
    expectLocal(t, "a = trunc(-3.8)", "a", -3.0)
    expectLocal(t, "a = trunc(-5.5)", "a", -5.0)
}

func Test_Pow(t *testing.T) {
    expectRuntimeError(t, "a = pow()")
    expectRuntimeError(t, "a = pow(b)")
    expectSuccess(t, "a = pow(b, c)")
    expectRuntimeError(t, "a = pow(b, c, d)")
    expectLocal(t, "a = pow(0, 0)", "a", 1.0)
    expectLocal(t, "a = pow(1, 0)", "a", 1.0)
    expectLocal(t, "a = pow(1, 1)", "a", 1.0)
    expectLocal(t, "a = pow(2, 1)", "a", 2.0)
    expectLocal(t, "a = pow(2, 2)", "a", 4.0)
    expectLocal(t, "a = pow(2, 3)", "a", 8.0)
    expectLocal(t, "a = pow(4, -1)", "a", 0.25)
    expectLocal(t, "a = pow(2, -3)", "a", 0.125)
    expectLocal(t, "a = pow(-2, 2)", "a", 4.0)
    expectLocal(t, "a = pow(-2, 3)", "a", -8.0)
    expectLocal(t, "a = pow(-2, -2)", "a", 0.25)
    expectLocal(t, "a = pow(-2, -3)", "a", -0.125)
}

func Test_Sqrt(t *testing.T) {
    expectRuntimeError(t, "a = sqrt()")
    expectSuccess(t, "a = sqrt(b)")
    expectRuntimeError(t, "a = sqrt(b, c)")
    expectLocalClose(t, "a = sqrt(4)", "a", 2.0)
    expectLocalClose(t, "a = sqrt(100)", "a", 10.0)
    expectLocalClose(t, "a = sqrt(2)", "a", math.Sqrt(2.0))
    expectSuccess(t, "a = sqrt(0)")
    expectRuntimeError(t, "a = sqrt(-1)")
    expectRuntimeError(t, "a = sqrt(-4.5)")
}

func Test_Sin(t *testing.T) {
    expectRuntimeError(t, "a = sin()")
    expectSuccess(t, "a = sin(b)")
    expectRuntimeError(t, "a = sin(b, c)")
    expectLocalClose(t, "a = sin(-3.14159)", "a", 0.0)
    expectLocalClose(t, "a = sin(-3.14159 / 2)", "a", -1.0)
    expectLocalClose(t, "a = sin(0)", "a", 0.0)
    expectLocalClose(t, "a = sin(3.14159 / 6)", "a", 0.5)
    expectLocalClose(t, "a = sin(3.14159 / 4)", "a", math.Sqrt(2.0) / 2.0)
    expectLocalClose(t, "a = sin(3.14159 / 3)", "a", math.Sqrt(3.0) / 2.0)
    expectLocalClose(t, "a = sin(3.14159 / 2)", "a", 1.0)
    expectLocalClose(t, "a = sin(3.14159)", "a", 0.0)
}

func Test_Cos(t *testing.T) {
    expectRuntimeError(t, "a = cos()")
    expectSuccess(t, "a = cos(b)")
    expectRuntimeError(t, "a = cos(b, c)")
    expectLocalClose(t, "a = cos(-3.14159)", "a", -1.0)
    expectLocalClose(t, "a = cos(-3.14159 / 2)", "a", 0.0)
    expectLocalClose(t, "a = cos(0)", "a", 1.0)
    expectLocalClose(t, "a = cos(3.14159 / 6)", "a", math.Sqrt(3.0) / 2.0)
    expectLocalClose(t, "a = cos(3.14159 / 4)", "a", math.Sqrt(2.0) / 2.0)
    expectLocalClose(t, "a = cos(3.14159 / 3)", "a", 0.5)
    expectLocalClose(t, "a = cos(3.14159 / 2)", "a", 0.0)
    expectLocalClose(t, "a = cos(3.14159)", "a", -1.0)
}

func Test_Tan(t *testing.T) {
    expectRuntimeError(t, "a = tan()")
    expectSuccess(t, "a = tan(b)")
    expectRuntimeError(t, "a = tan(b, c)")
    expectLocalClose(t, "a = tan(-3.14159)", "a", 0.0)
    expectLocalClose(t, "a = tan(0)", "a", 0.0)
    expectLocalClose(t, "a = tan(3.14159 / 6)", "a", math.Sqrt(3.0) / 3.0)
    expectLocalClose(t, "a = tan(3.14159 / 4)", "a", 1.0)
    expectLocalClose(t, "a = tan(3.14159 / 3)", "a", math.Sqrt(3.0))
    expectLocalClose(t, "a = tan(3.14159)", "a", 0.0)
    expectRuntimeError(t, fmt.Sprintf("a = tan(%f / 2)", math.Pi))
    expectRuntimeError(t, fmt.Sprintf("a = tan(-%f / 2)", math.Pi))
    expectRuntimeError(t, fmt.Sprintf("a = tan(3 * %f / 2)", math.Pi))
    expectRuntimeError(t, fmt.Sprintf("a = tan(-3 * %f / 2)", math.Pi))
}

func Test_Asin(t *testing.T) {
    expectRuntimeError(t, "a = asin()")
    expectSuccess(t, "a = asin(b)")
    expectRuntimeError(t, "a = asin(b, c)")
    expectLocalClose(t, "a = asin(0)", "a", 0.0)
    expectLocalClose(t, "a = asin(1)", "a", math.Pi / 2)
    expectLocalClose(t, "a = asin(-1)", "a", -math.Pi / 2)
}

func Test_Acos(t *testing.T) {
    expectRuntimeError(t, "a = acos()")
    expectSuccess(t, "a = acos(b)")
    expectRuntimeError(t, "a = acos(b, c)")
    expectLocalClose(t, "a = acos(0)", "a", math.Pi / 2)
    expectLocalClose(t, "a = acos(1)", "a", 0.0)
    expectLocalClose(t, "a = acos(-1)", "a", math.Pi)
}

func Test_Atan(t *testing.T) {
    expectRuntimeError(t, "a = atan()")
    expectSuccess(t, "a = atan(b)")
    expectRuntimeError(t, "a = atan(b, c)")
    expectLocalClose(t, "a = atan(0)", "a", 0.0)
    expectLocalClose(t, "a = atan(1 / sqrt(3))", "a", math.Pi / 6.0)
    expectLocalClose(t, "a = atan(1)", "a", math.Pi / 4)
    expectLocalClose(t, "a = atan(sqrt(3))", "a", math.Pi / 3.0)
}

func Test_UndefinedFunction(t *testing.T) {
    if _, err := evaluateFunction("foo", []float64{}); err == nil {
        t.Error("expected error for undefined function")
    }
}

//
// 5 - Branching
//

func Test_If(t *testing.T) {
    expectLocal(t, "x = 0\nif (x) { a = 2 }", "a", 0.0)
    expectLocal(t, "x = 1\nif (x) { a = 2 }", "a", 2.0)
    expectLocal(t, "x = 100\nif (x) { a = 2 }", "a", 2.0)
    expectLocal(t, "x = -1\nif (x) { a = 2 }", "a", 2.0)
}

func Test_Else(t *testing.T) {
    expectLocal(t, "x = 0\nif (x) { a = 2 } else { a = 3 }", "a", 3.0)
    expectLocal(t, "x = 1\nif (x) { a = 2 } else { a = 3 }", "a", 2.0)
    expectLocal(t, "x = 100\nif (x) { a = 2 } else { a = 3 }", "a", 2.0)
    expectLocal(t, "x = -1\nif (x) { a = 2 } else { a = 3 }", "a", 2.0)
}

func Test_ElseIf(t *testing.T) {
    expectLocal(t, "x=0\ny=0\nif (x) { a = 2 } else if (y) { a = 3 } else { a = 4 }", "a", 4.0)
    expectLocal(t, "x=1\ny=0\nif (x) { a = 2 } else if (y) { a = 3 } else { a = 4 }", "a", 2.0)
    expectLocal(t, "x=0\ny=1\nif (x) { a = 2 } else if (y) { a = 3 } else { a = 4 }", "a", 3.0)
    expectLocal(t, "x=1\ny=1\nif (x) { a = 2 } else if (y) { a = 3 } else { a = 4 }", "a", 2.0)
}

func Test_MultiElseIf(t *testing.T) {
    expectLocal(t, "x=0\ny=0\nz=0\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 5.0)
    expectLocal(t, "x=1\ny=0\nz=0\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 2.0)
    expectLocal(t, "x=0\ny=1\nz=0\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 3.0)
    expectLocal(t, "x=1\ny=1\nz=0\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 2.0)
    expectLocal(t, "x=0\ny=0\nz=1\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 4.0)
    expectLocal(t, "x=1\ny=0\nz=1\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 2.0)
    expectLocal(t, "x=0\ny=1\nz=1\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 3.0)
    expectLocal(t, "x=1\ny=1\nz=1\nif (x) { a = 2 } else if (y) { a = 3 } else if (z) { a = 4 } else { a = 5 }", "a", 2.0)
}

func Test_NestedIf(t *testing.T) {
    expectLocal(t, "x=0\ny=0\nif (x) { if (y) { a = 1 } else { a = 2 } } else { if (y) { a = 3 } else { a = 4 } }", "a", 4.0)
    expectLocal(t, "x=1\ny=0\nif (x) { if (y) { a = 1 } else { a = 2 } } else { if (y) { a = 3 } else { a = 4 } }", "a", 2.0)
    expectLocal(t, "x=0\ny=1\nif (x) { if (y) { a = 1 } else { a = 2 } } else { if (y) { a = 3 } else { a = 4 } }", "a", 3.0)
    expectLocal(t, "x=1\ny=1\nif (x) { if (y) { a = 1 } else { a = 2 } } else { if (y) { a = 3 } else { a = 4 } }", "a", 1.0)
}

//
// 6 - Looping
//

func Test_FiniteLoop(t *testing.T) {
    input := "while (i < 10) { a = i\ni = i + 1 }"
    expectLocal(t, input, "i", 10.0)
    expectLocal(t, input, "a", 9.0)
}

func Test_NestedLoop(t *testing.T) {
    input := "j = 1\nwhile (i < 3) { while (j < i * 30) { a = j\nj = j + 10 } i = i + 1 }"
    expectLocal(t, input, "i", 3.0)
    expectLocal(t, input, "j", 61.0)
    expectLocal(t, input, "a", 51.0)
}

func Test_InfiniteLoop(t *testing.T) {
    expectRuntimeError(t, "while (true) { }")
}
