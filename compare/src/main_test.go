package main

import (
	"testing"
)

func TestPixelMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/white.png", "-c", "pixel"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Pixel compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Pixel compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 1.0 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 1.0", comparisons[0].Results[0].Index)
	}
}

func TestPixelDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "pixel"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Pixel compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Pixel compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 0.0 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 0.0", comparisons[0].Results[0].Index)
	}
}

func TestContrastMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/contrastA.png", "-B", "../../testAssets/contrastB.png", "-c", "contrast"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Contrast compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Contrast compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 1.0 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 1.0", comparisons[0].Results[0].Index)
	}
}

func TestContrastDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "contrast"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Contrast compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Contrast compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 0.0 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 0.0", comparisons[0].Results[0].Index)
	}
}

func TestQuadMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/quadA.png", "-B", "../../testAssets/quadB.png", "-c", "quad"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Quad compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Quad compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 1.0 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 1.0", comparisons[0].Results[0].Index)
	}
}

func TestQuadDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "quad"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("Quad compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("Quad compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 0.0 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 0.0", comparisons[0].Results[0].Index)
	}
}

func TestSSIMMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/white.png", "-c", "ssim"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("SSIM compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("SSIM compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 1.0 {
		t.Errorf("SSIM compare test failed, compare value was %v, expected value 1.0", comparisons[0].Results[0].Index)
	}
}

func TestSSIMDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "ssim"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("SSIM compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("SSIM compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index > 0.0001 {
		t.Errorf("SSIM compare test failed, compare value was %v, expected value < 0.0001", comparisons[0].Results[0].Index)
	}
}

func TestMSEMatch(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/white.png", "-c", "ssim"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("MSE compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("MSE compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index != 1.0 {
		t.Errorf("MSE compare test failed, compare value was %v, expected value 1.0", comparisons[0].Results[0].Index)
	}
}

func TestMSEDiff(t *testing.T) {
	args := []string{"-A", "../../testAssets/white.png", "-B", "../../testAssets/black.png", "-c", "ssim"}

	comparisons := run(args)

	if len(comparisons[0].Results) == 0 {
		t.Error("MSE compare test failed, 0 comparison results")
	}

	if len(comparisons[0].Results) > 1 {
		t.Error("MSE compare test failed, more than 1 comparison result")
	}

	if comparisons[0].Results[0].Index > 0.0001 {
		t.Errorf("SSIM compare test failed, compare value was %v, expected value < 0.0001", comparisons[0].Results[0].Index)
	}
}

func TestAll(t *testing.T) {
	args := []string{"-A", "../../testAssets/screenA.png", "-B", "../../testAssets/screenB.png"}

	comparisons := run(args)

	if len(comparisons[0].Results) < 5 {
		t.Error("All compare test failed, less than 5 comparison results")
	}

	if len(comparisons[0].Results) > 5 {
		t.Error("All compare test failed, more than 5 comparison result")
	}

	if comparisons[0].Results[0].Index != 0.9948143325617284 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 0.0051", comparisons[0].Results[0].Index)
	}

	if comparisons[0].Results[0].NumFailed != 10753 {
		t.Errorf("Pixel compare test failed, compare value was %v, expected value 10753", comparisons[0].Results[0].NumFailed)
	}

	if comparisons[0].Results[1].Index != 0.9978877314814815 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 0.0021", comparisons[0].Results[1].Index)
	}

	if comparisons[0].Results[1].NumFailed != 4380 {
		t.Errorf("Contrast compare test failed, compare value was %v, expected value 4380", comparisons[0].Results[1].NumFailed)
	}

	if comparisons[0].Results[2].Index != 0.9985976080246913 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 0.0014", comparisons[0].Results[2].Index)
	}

	if comparisons[0].Results[2].NumFailed != 2908 {
		t.Errorf("Quad compare test failed, compare value was %v, expected value 2908", comparisons[0].Results[2].NumFailed)
	}

	if comparisons[0].Results[3].Index != 0.982134444334117 {
		t.Errorf("SSIM compare test failed, compare value was %v, expected value 0.982134444334117", comparisons[0].Results[3].Index)
	}

	if comparisons[0].Results[3].NumFailed != -1 {
		t.Errorf("SSIM compare test failed, compare value was %v, expected value -1", comparisons[0].Results[3].NumFailed)
	}

	if comparisons[0].Results[4].Index != 0.9997884516460905 {
		t.Errorf("MSE compare test failed, compare value was %v, expected value 0.9997884516460905", comparisons[0].Results[4].Index)
	}

	if comparisons[0].Results[4].NumFailed != -1 {
		t.Errorf("MSE compare test failed, compare value was %v, expected value -1", comparisons[0].Results[4].NumFailed)
	}
}

func TestPixelDir(t *testing.T) {
	args := []string{"-A", "../../testAssets/DirA", "-B", "../../testAssets/DirB", "-c", "pixel"}

	comparisons := run(args)

	if len(comparisons) != 3 {
		t.Error("Pixel directory test failed, not enough results")
	}

	if comparisons[0].Results[0].NumFailed != 48 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 48", comparisons[0].Results[0].NumFailed)
	}

	if comparisons[1].Results[0].NumFailed != 10753 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 10753", comparisons[0].Results[1].NumFailed)
	}

	if comparisons[2].Results[0].Index != 1.0 {
		t.Errorf("Pixel directory test failed, compare value was %v, expected value 1.0", comparisons[0].Results[2].Index)
	}
}
