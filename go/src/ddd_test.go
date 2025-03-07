package ddd_test

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	ddd "github.com/theandrew168/exploring-ddd/src"
)

func makeBatchAndLine(sku string, batchQty, lineQty int) (*ddd.Batch, ddd.OrderLine) {
	batch := ddd.NewBatch("batch-001", sku, batchQty)
	line := ddd.NewOrderLine("order-123", sku, lineQty)
	return batch, line
}

func TestAllocatingToABatchReducesAvailableQuantity(t *testing.T) {
	batch := ddd.NewBatch("batch-001", "SMALL-TABLE", 20)
	line := ddd.NewOrderLine("order-ref", "SMALL-TABLE", 2)

	batch.Allocate(line)

	AssertEqual(t, batch.AvailableQuantity(), 18)
}

func TestCanAllocateIfAvailableGreaterThanRequired(t *testing.T) {
	largeBatch, smallLine := makeBatchAndLine("ELEGANT-LAMP", 20, 2)
	AssertEqual(t, largeBatch.CanAllocate(smallLine), true)
}

func TestCannotAllocateIfAvailableSmallerThanRequired(t *testing.T) {
	smallBatch, largeLine := makeBatchAndLine("ELEGANT-LAMP", 2, 20)
	AssertEqual(t, smallBatch.CanAllocate(largeLine), false)
}

func TestCanAllocateIfAvailableEqualsRequired(t *testing.T) {
	batch, line := makeBatchAndLine("ELEGANT-LAMP", 2, 2)
	AssertEqual(t, batch.CanAllocate(line), true)
}

func TestCannotAllocateIfSKUDoesNotMatch(t *testing.T) {
	batch := ddd.NewBatch("batch-001", "UNCOMFORTABLE-CHAIR", 100)
	differentSKULine := ddd.NewOrderLine("order-123", "EXPENSIVE-TOASTER", 10)
	AssertEqual(t, batch.CanAllocate(differentSKULine), false)
}

func TestAllocationIsIdempotent(t *testing.T) {
	batch, line := makeBatchAndLine("EXPENSIVE-FOOTSTOOL", 20, 2)
	batch.Allocate(line)
	batch.Allocate(line)
	AssertEqual(t, batch.AvailableQuantity(), 18)
}

func TestDeallocate(t *testing.T) {
	batch, line := makeBatchAndLine("EXPENSIVE-FOOTSTOOL", 20, 2)
	batch.Allocate(line)
	batch.Deallocate(line)
	AssertEqual(t, batch.AvailableQuantity(), 20)
}

func TestCanOnlyDeallocateAllocatedLines(t *testing.T) {
	batch, unallocatedLine := makeBatchAndLine("DECORATIVE-TRINKET", 20, 2)
	batch.Deallocate(unallocatedLine)
	AssertEqual(t, batch.AvailableQuantity(), 20)
}

func TestPrefersCurrentStockBatchesToShipments(t *testing.T) {
	inStockBatch := ddd.NewBatch("in-stock-batch", "RETRO-CLOCK", 100)
	shipmentBatch := ddd.NewBatchWithETA("shipment-batch", "RETRO-CLOCK", 100, time.Now().AddDate(0, 0, 1))
	line := ddd.NewOrderLine("oref", "RETRO-CLOCK", 10)

	ddd.Allocate(line, []*ddd.Batch{inStockBatch, shipmentBatch})

	AssertEqual(t, inStockBatch.AvailableQuantity(), 90)
	AssertEqual(t, shipmentBatch.AvailableQuantity(), 100)
}

func TestPrefersEarlierBatches(t *testing.T) {
	earliest := ddd.NewBatchWithETA("speedy-batch", "MINIMALIST-SPOON", 100, time.Now())
	medium := ddd.NewBatchWithETA("normal-batch", "MINIMALIST-SPOON", 100, time.Now().AddDate(0, 0, 1))
	latest := ddd.NewBatchWithETA("slow-batch", "MINIMALIST-SPOON", 100, time.Now().AddDate(0, 0, 2))
	line := ddd.NewOrderLine("order1", "MINIMALIST-SPOON", 10)

	ddd.Allocate(line, []*ddd.Batch{medium, earliest, latest})

	AssertEqual(t, earliest.AvailableQuantity(), 90)
	AssertEqual(t, medium.AvailableQuantity(), 100)
	AssertEqual(t, latest.AvailableQuantity(), 100)
}

func TestReturnsAllocatedBatchRef(t *testing.T) {
	inStockBatch := ddd.NewBatch("in-stock-batch-ref", "HIGHBROW-POSTER", 100)
	shipmentBatch := ddd.NewBatchWithETA("shipment-batch-ref", "HIGHBROW-POSTER", 100, time.Now().AddDate(0, 0, 1))
	line := ddd.NewOrderLine("oref", "HIGHBROW-POSTER", 10)
	allocation, _ := ddd.Allocate(line, []*ddd.Batch{inStockBatch, shipmentBatch})
	AssertEqual(t, allocation, inStockBatch.Ref())
}

func TestReturnsOutOfStockErrorIfCannotAllocate(t *testing.T) {
	batch := ddd.NewBatchWithETA("batch1", "SMALL-FORK", 10, time.Now())
	ddd.Allocate(ddd.NewOrderLine("order1", "SMALL-FORK", 10), []*ddd.Batch{batch})

	_, err := ddd.Allocate(ddd.NewOrderLine("order2", "SMALL-FORK", 1), []*ddd.Batch{batch})
	AssertErrorIs(t, err, ddd.ErrOutOfStock)
}

func AssertEqual(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func AssertNotEqual(t *testing.T, got, notwant any) {
	t.Helper()

	if reflect.DeepEqual(got, notwant) {
		t.Fatalf("got (but don't want) %v", got)
	}
}

func AssertAtLeast(t *testing.T, got, want int) {
	t.Helper()

	if got < want {
		t.Fatalf("got %v; want at least %v", got, want)
	}
}

func AssertStringContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("got %q; want to contain: %q", got, want)
	}
}

func AssertSliceContains[T comparable](t *testing.T, got []T, want T) {
	t.Helper()

	if !slices.Contains(got, want) {
		t.Fatalf("got %v; want to contain: %v", got, want)
	}
}

func AssertNilError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf("got: %v; want: nil", got)
	}
}

func AssertErrorIs(t *testing.T, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: %q", want)
	}

	if !errors.Is(got, want) {
		t.Fatalf("got %q; want: %q", got, want)
	}
}

func AssertErrorAs(t *testing.T, got error, want any) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: %T", want)
	}

	if !errors.As(got, want) {
		t.Fatalf("got %q; want: %T", got, want)
	}
}

func AssertErrorContains(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: error to contain: %q", want)
	}

	if !strings.Contains(got.Error(), want) {
		t.Fatalf("got %q; want to contain: %q", got.Error(), want)
	}
}
