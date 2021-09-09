package msf

import "testing"

func Test_intToHexString(t *testing.T) {
    type testCase struct {
        Value uint
        MinHexDigits int
        Expected string
    }
    trials := []testCase{
        {0, 4, "0000"},
        {0, 8, "00000000"},
        {1, 4, "0001"},
        {15, 4, "000f"},
        {16, 4, "0010"},
        {256, 4, "0100"},
    }
    for _, trial := range trials {
        result := intToHexString(trial.Value, trial.MinHexDigits)
        if result != trial.Expected {
            t.Fatalf("Expected '%s', got '%s'", trial.Expected, result)
        }
    }
}

func Test_int16ToHexString(t *testing.T) {
    type testCase struct {
        Value int16
        Expected string
    }
    trials := []testCase{
        {0, "0000"},
        {1, "0001"},
        {-1, "ffff"},
        {15, "000f"},
        {16, "0010"},
        {256, "0100"},
        {-256, "ff00"},
    }
    for _, trial := range trials {
        result := int16ToHexString(trial.Value)
        if result != trial.Expected {
            t.Fatalf("Expected '%s', got '%s'", trial.Expected, result)
        }
    }
}

func Test_floatToHexString(t *testing.T) {
    type testCase struct {
        Value float32
        Expected string
    }
    trials := []testCase{
        {0, "00000000"},
        {0.1, "3dcccccd"},
        {99.5, "42c70000"},
        {-83.2, "c2a66666"},
    }
    for _, trial := range trials {
        result := floatToHexString(trial.Value)
        if result != trial.Expected {
            t.Fatalf("Expected '%s', got '%s'", trial.Expected, result)
        }
    }
}