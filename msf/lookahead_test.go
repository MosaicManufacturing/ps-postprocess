package msf

import (
	"math"
	"path"
	"testing"
)

func expectMatchingFloat32s(t *testing.T, result, expected float32) {
	if math.Abs(float64(result-expected)) > 1e-6 {
		t.Fatalf("expected %f, got %f", expected, result)
	}
}

func TestLookahead_ZLiftRetraction(t *testing.T) {
	preflightResults, err := preflightWithInputFiles(t, path.Join("side-transition", "1"))
	if err != nil {
		t.Fatal(err)
	}

	firstTransitionNextPosition := preflightResults.transitionNextPositions[0]
	expectMatchingFloat32s(t, firstTransitionNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Z, 4.2)
	if firstTransitionNextPosition.MovedXY != true {
		t.Fatal("expected firstTransitionNextPosition.MovedXY to be true")
	}
	if firstTransitionNextPosition.MovedZ != true {
		t.Fatal("expected firstTransitionNextPosition.MovedZ to be true")
	}

	firstScriptNextPosition := preflightResults.scriptNextPositions[0]
	expectMatchingFloat32s(t, firstScriptNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstScriptNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstScriptNextPosition.Z, 4.2)
	if firstScriptNextPosition.InitializedX != true {
		t.Fatal("expected firstScriptNextPosition.InitializedX to be true")
	}
	if firstScriptNextPosition.InitializedY != true {
		t.Fatal("expected firstScriptNextPosition.InitializedY to be true")
	}
	if firstScriptNextPosition.InitializedZ != true {
		t.Fatal("expected firstScriptNextPosition.InitializedZ to be true")
	}
}

func TestLookahead_ZLiftNoRetraction(t *testing.T) {
	preflightResults, err := preflightWithInputFiles(t, path.Join("side-transition", "2"))
	if err != nil {
		t.Fatal(err)
	}

	firstTransitionNextPosition := preflightResults.transitionNextPositions[0]
	expectMatchingFloat32s(t, firstTransitionNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Z, 4.2)
	if firstTransitionNextPosition.MovedXY != true {
		t.Fatal("expected firstTransitionNextPosition.MovedXY to be true")
	}
	if firstTransitionNextPosition.MovedZ != false {
		t.Fatal("expected firstTransitionNextPosition.MovedZ to be false")
	}

	firstScriptNextPosition := preflightResults.scriptNextPositions[0]
	expectMatchingFloat32s(t, firstScriptNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstScriptNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstScriptNextPosition.Z, 0)
	if firstScriptNextPosition.InitializedX != true {
		t.Fatal("expected firstScriptNextPosition.InitializedX to be true")
	}
	if firstScriptNextPosition.InitializedY != true {
		t.Fatal("expected firstScriptNextPosition.InitializedY to be true")
	}
	if firstScriptNextPosition.InitializedZ != false {
		t.Fatal("expected firstScriptNextPosition.InitializedZ to be false")
	}
}

func TestLookahead_NoZLiftRetraction(t *testing.T) {
	preflightResults, err := preflightWithInputFiles(t, path.Join("side-transition", "3"))
	if err != nil {
		t.Fatal(err)
	}

	firstTransitionNextPosition := preflightResults.transitionNextPositions[0]
	expectMatchingFloat32s(t, firstTransitionNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Z, 4.2)
	if firstTransitionNextPosition.MovedXY != true {
		t.Fatal("expected firstTransitionNextPosition.MovedXY to be true")
	}
	if firstTransitionNextPosition.MovedZ != false {
		t.Fatal("expected firstTransitionNextPosition.MovedZ to be false")
	}

	firstScriptNextPosition := preflightResults.scriptNextPositions[0]
	expectMatchingFloat32s(t, firstScriptNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstScriptNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstScriptNextPosition.Z, 0)
	if firstScriptNextPosition.InitializedX != true {
		t.Fatal("expected firstScriptNextPosition.InitializedX to be true")
	}
	if firstScriptNextPosition.InitializedY != true {
		t.Fatal("expected firstScriptNextPosition.InitializedY to be true")
	}
	if firstScriptNextPosition.InitializedZ != false {
		t.Fatal("expected firstScriptNextPosition.InitializedZ to be false")
	}
}

func TestLookahead_NoZLiftNoRetraction(t *testing.T) {
	preflightResults, err := preflightWithInputFiles(t, path.Join("side-transition", "4"))
	if err != nil {
		t.Fatal(err)
	}

	firstTransitionNextPosition := preflightResults.transitionNextPositions[0]
	expectMatchingFloat32s(t, firstTransitionNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstTransitionNextPosition.Z, 4.2)
	if firstTransitionNextPosition.MovedXY != true {
		t.Fatal("expected firstTransitionNextPosition.MovedXY to be true")
	}
	if firstTransitionNextPosition.MovedZ != false {
		t.Fatal("expected firstTransitionNextPosition.MovedZ to be false")
	}

	firstScriptNextPosition := preflightResults.scriptNextPositions[0]
	expectMatchingFloat32s(t, firstScriptNextPosition.X, 136.866)
	expectMatchingFloat32s(t, firstScriptNextPosition.Y, 71.767)
	expectMatchingFloat32s(t, firstScriptNextPosition.Z, 0)
	if firstScriptNextPosition.InitializedX != true {
		t.Fatal("expected firstScriptNextPosition.InitializedX to be true")
	}
	if firstScriptNextPosition.InitializedY != true {
		t.Fatal("expected firstScriptNextPosition.InitializedY to be true")
	}
	if firstScriptNextPosition.InitializedZ != false {
		t.Fatal("expected firstScriptNextPosition.InitializedZ to be false")
	}
}
