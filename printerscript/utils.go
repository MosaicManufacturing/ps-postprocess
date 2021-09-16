package printerscript

import (
    "errors"
    "fmt"
    "math"
)

func getArity(fn string) (int, error) {
    // minimum arity
    if fn == "min" { return 2, nil }
    if fn == "max" { return 2, nil }
    // exact arity
    if fn == "pow" { return 2, nil }
    if fn == "abs" { return 1, nil }
    if fn == "floor" { return 1, nil }
    if fn == "ceil" { return 1, nil }
    if fn == "round" { return 1, nil }
    if fn == "trunc" { return 1, nil }
    if fn == "sqrt" { return 1, nil }
    if fn == "sin" { return 1, nil }
    if fn == "cos" { return 1, nil }
    if fn == "tan" { return 1, nil }
    if fn == "asin" { return 1, nil }
    if fn == "acos" { return 1, nil }
    if fn == "atan" { return 1, nil }
    return 0, fmt.Errorf("unknown function: '%s'", fn)
}

func variadicMin(args []float64) (float64, error) {
    if len(args) < 2 {
        return 0, errors.New("expected at least 2 arguments to 'min'")
    }
    min := args[0]
    for i := 0; i < len(args); i++ {
        min = math.Min(min, args[i])
    }
    return min, nil
}

func variadicMax(args []float64) (float64, error) {
    if len(args) < 2 {
        return 0, errors.New("expected at least 2 arguments to 'max'")
    }
    max := args[0]
    for i := 0; i < len(args); i++ {
        max = math.Max(max, args[i])
    }
    return max, nil
}

func evaluateFunction(fn string, args []float64) (float64, error) {
    if fn == "min" { return variadicMin(args) }
    if fn == "max" { return variadicMax(args) }
    if fn == "pow" {
        if args[0] == 0 && args[1] == 0 { return 1, nil }
        return math.Pow(args[0], args[1]), nil
    }
    if fn == "abs" { return math.Abs(args[0]), nil }
    if fn == "floor" { return math.Floor(args[0]), nil }
    if fn == "ceil" { return math.Ceil(args[0]), nil }
    if fn == "round" { return math.Round(args[0]), nil }
    if fn == "trunc" { return math.Trunc(args[0]), nil }
    if fn == "sqrt" {
        if args[0] < 0 {
            return 0, errors.New("expected non-negative argument to 'sqrt'")
        }
        return math.Sqrt(args[0]), nil
    }
    if fn == "sin" { return math.Sin(args[0]), nil }
    if fn == "cos" { return math.Cos(args[0]), nil }
    if fn == "tan" {
        if math.Mod(math.Abs(args[0] - (math.Pi / 2)), math.Pi) < 1e-5 {
            return 0, errors.New("undefined result of 'tan'")
        }
        return math.Tan(args[0]), nil
    }
    if fn == "asin" { return math.Asin(args[0]), nil }
    if fn == "acos" { return math.Acos(args[0]), nil }
    if fn == "atan" { return math.Atan(args[0]), nil }
    return 0, fmt.Errorf("unknown function: '%s'", fn)
}

func evaluateUnaryOp(op string, operand float64) (float64, error) {
    if op == "!" {
        if operand == 0 { return 1, nil }
        return 0, nil
    }
    if op == "-" {
        return -operand, nil
    }
    return 0, fmt.Errorf("unknown unary operator: '%s'", op)
}

func evaluateBinaryOp(op string, lhs, rhs float64) (float64, error) {
    if op == "*" { return lhs * rhs, nil }
    if op == "/" {
        if rhs == 0 { return 0, errors.New("division by zero encountered") }
        return lhs / rhs, nil
    }
    if op == "%" {
        if rhs == 0 { return 0, errors.New("mod by zero encountered") }
        return math.Mod(math.Mod(lhs, rhs) + rhs, rhs), nil
    }
    if op == "+" { return lhs + rhs, nil }
    if op == "-" { return lhs - rhs, nil }
    if op == "==" {
        if lhs == rhs { return 1, nil }
        return 0, nil
    }
    if op == "!=" {
        if lhs != rhs { return 1, nil }
        return 0, nil
    }
    if op == "<" {
        if lhs < rhs { return 1, nil }
        return 0, nil
    }
    if op == "<=" {
        if lhs <= rhs { return 1, nil }
        return 0, nil
    }
    if op == ">" {
        if lhs > rhs { return 1, nil }
        return 0, nil
    }
    if op == ">=" {
        if lhs >= rhs { return 1, nil }
        return 0, nil
    }
    if op == "||" {
        if (lhs != 0) || (rhs != 0) { return 1, nil }
        return 0, nil
    }
    if op == "&&" {
        if (lhs != 0) && (rhs != 0) { return 1, nil }
        return 0, nil
    }
    return 0, fmt.Errorf("unknown binary operator: '%s'", op)
}
