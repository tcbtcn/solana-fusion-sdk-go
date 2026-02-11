package math

import (
	"math/big"
	"testing"
)

func TestMulDiv_Basic(t *testing.T) {
	a := big.NewInt(100)
	b := big.NewInt(200)
	x := big.NewInt(50)

	result := MulDiv(a, b, x, RoundingFloor)
	expected := big.NewInt(400)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s, got %s", expected.String(), result.String())
	}
}

func TestMulDiv_RoundingCeil(t *testing.T) {
	a := big.NewInt(100)
	b := big.NewInt(201)
	x := big.NewInt(50)

	result := MulDiv(a, b, x, RoundingCeil)
	// (100 * 201) / 50 = 20100 / 50 = 402 exactly (no remainder)
	// Ceiling of 402.0 is 402 (TypeScript mulDiv returns 402 when remainder is 0)
	expected := big.NewInt(402)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s (ceiled), got %s", expected.String(), result.String())
	}
}

func TestMulDiv_RoundingFloor(t *testing.T) {
	a := big.NewInt(100)
	b := big.NewInt(201)
	x := big.NewInt(50)

	result := MulDiv(a, b, x, RoundingFloor)
	expected := big.NewInt(402)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s (floored), got %s", expected.String(), result.String())
	}
}

func TestMulDiv_DivisionByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for division by zero")
		}
	}()

	a := big.NewInt(100)
	b := big.NewInt(200)
	x := big.NewInt(0)

	MulDiv(a, b, x, RoundingFloor)
}

func TestMulDiv_LargeNumbers(t *testing.T) {
	a, _ := new(big.Int).SetString("1000000000000000000", 10)
	b, _ := new(big.Int).SetString("2000000000000000000", 10)
	x, _ := new(big.Int).SetString("500000000000000000", 10)

	result := MulDiv(a, b, x, RoundingFloor)
	expected, _ := new(big.Int).SetString("4000000000000000000", 10)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s, got %s", expected.String(), result.String())
	}
}

func TestCalcTakingAmount_Basic(t *testing.T) {
	swapMakerAmount := big.NewInt(500)
	orderMakerAmount := big.NewInt(1000)
	orderTakerAmount := big.NewInt(2000)

	result := CalcTakingAmount(swapMakerAmount, orderMakerAmount, orderTakerAmount)
	expected := big.NewInt(1000)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s, got %s", expected.String(), result.String())
	}
}

func TestCalcTakingAmount_ZeroMakerAmount(t *testing.T) {
	swapMakerAmount := big.NewInt(500)
	orderMakerAmount := big.NewInt(0)
	orderTakerAmount := big.NewInt(2000)

	result := CalcTakingAmount(swapMakerAmount, orderMakerAmount, orderTakerAmount)
	expected := big.NewInt(0)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s for zero maker amount, got %s", expected.String(), result.String())
	}
}

func TestCalcTakingAmount_PartialFill(t *testing.T) {
	swapMakerAmount := big.NewInt(250)
	orderMakerAmount := big.NewInt(1000)
	orderTakerAmount := big.NewInt(2000)

	result := CalcTakingAmount(swapMakerAmount, orderMakerAmount, orderTakerAmount)
	expected := big.NewInt(500)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s, got %s", expected.String(), result.String())
	}
}

func TestCalcTakingAmount_Ceiling(t *testing.T) {
	swapMakerAmount := big.NewInt(333)
	orderMakerAmount := big.NewInt(1000)
	orderTakerAmount := big.NewInt(2000)

	result := CalcTakingAmount(swapMakerAmount, orderMakerAmount, orderTakerAmount)
	expected := big.NewInt(666)

	if result.Cmp(expected) != 0 {
		t.Errorf("Expected %s (ceiled), got %s", expected.String(), result.String())
	}
}
