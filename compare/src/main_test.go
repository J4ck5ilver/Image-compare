package main

import (
	"testing"
)

func TestPixelMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/white.png", "-c", "pixel"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Pixel compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Pixel compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 1.0 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 1.0", results[0].Fraction)
	}
}

func TestPixelDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "pixel"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Pixel compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Pixel compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 0.0 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 0.0", results[0].Fraction)
	}
}

func TestContrastMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/contrastA.png", "-B", "../../testAssets/contrastB.png", "-c", "contrast"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Contrast compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Contrast compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 1.0 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 1.0", results[0].Fraction)
	}
}

func TestContrastDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "contrast"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Contrast compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Contrast compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 0.0 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 0.0", results[0].Fraction)
	}
}

func TestQuadMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/quadA.png", "-B", "../../testAssets/quadB.png", "-c", "quad"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Quad compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Quad compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 1.0 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 1.0", results[0].Fraction)
	}
}

func TestQuadDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "quad"}

	results := run(args)

	if len(results) == 0 {
		t.Error("Quad compare test failed, 0 comparison results")
	}

	if len(results) > 1 {
		t.Error("Quad compare test failed, more than 1 comparison result")
	}

	if results[0].Fraction != 0.0 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 0.0", results[0].Fraction)
	}
}

func TestAll(t *testing.T) {
	args := []string{"-A", "../../testAssets/screenA.png", "-B", "../../testAssets/screenB.png"}

	results := run(args)

	if len(results) < 3 {
		t.Error("All compare test failed, less than 3 comparison results")
	}

	if len(results) > 3 {
		t.Error("All compare test failed, more than 3 comparison result")
	}

	if results[0].Fraction != 0.9948143325617284 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 0.0051", results[0].Fraction)
	}

	if results[0].NumFailed != 10753 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 10753", results[0].NumFailed)
	}

	if results[1].Fraction != 0.9978877314814815 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 0.0021", results[1].Fraction)
	}

	if results[1].NumFailed != 4380 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 4380", results[1].NumFailed)
	}

	if results[2].Fraction != 0.9985976080246913 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 0.0014", results[2].Fraction)
	}

	if results[2].NumFailed != 727 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 727", results[2].NumFailed)
	}
}

func TestPixelDir(t *testing.T) {
	args := []string{"-A", "../../testAssets/DirA", "-B", "../../testAssets/DirB", "-c", "pixel"}

	results := run(args)

	if len(results) != 3 {
		t.Error("Pixel directory test failed, not enough results")
	}

	if results[0].NumFailed != 48 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 48", results[0].NumFailed)
	}

	if results[1].NumFailed != 10753 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 10753", results[1].NumFailed)
	}

	if results[2].Fraction != 1.0 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 1.0", results[2].Fraction)
	}
}
