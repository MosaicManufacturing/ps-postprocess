package interpreter

import "testing"

func Test_getArity(t *testing.T) {
	// TODO: implement tests
}

func Test_variadicMin(t *testing.T) {
	// TODO: implement tests
}

func Test_variadicMax(t *testing.T) {
	// TODO: implement tests
}

func Test_evaluateFunction(t *testing.T) {
	// TODO: implement tests for defined functions
	if _, err := evaluateFunction("foo", []float64{}); err == nil {
		t.Error("expected error for undefined function")
	}
}

func Test_evaluateUnaryOp(t *testing.T) {
	// TODO: implement tests for defined unary operators
	if _, err := evaluateUnaryOp("foo", 0); err == nil {
		t.Error("expected error for undefined unary operator")
	}
}

func Test_evaluateBinaryOp(t *testing.T) {
	// TODO: implement tests for defined binary operators
	if _, err := evaluateBinaryOp("foo", 0, 0); err == nil {
		t.Error("expected error for undefined binary operator")
	}
}
