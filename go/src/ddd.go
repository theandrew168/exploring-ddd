package ddd

import (
	"errors"
	"sort"
	"time"
)

var ErrOutOfStock = errors.New("out of stock")

type OrderLine struct {
	OrderID string
	SKU     string
	Qty     int
}

func NewOrderLine(orderID, sku string, qty int) OrderLine {
	return OrderLine{
		OrderID: orderID,
		SKU:     sku,
		Qty:     qty,
	}
}

// Ayy, we can use a set store OrderLines since structs (values, not pointers) are compared by value!
type OrderLineSet map[OrderLine]struct{}

type Batch struct {
	ref string
	sku string
	eta time.Time

	purchasedQty int
	allocations  OrderLineSet
}

func NewBatch(ref, sku string, qty int) *Batch {
	b := Batch{
		ref: ref,
		sku: sku,

		purchasedQty: qty,
		allocations:  make(OrderLineSet),
	}
	return &b
}

func NewBatchWithETA(ref, sku string, qty int, eta time.Time) *Batch {
	b := Batch{
		ref: ref,
		sku: sku,
		eta: eta,

		purchasedQty: qty,
		allocations:  make(OrderLineSet),
	}
	return &b
}

func (b *Batch) Ref() string {
	return b.ref
}

func (b *Batch) SKU() string {
	return b.sku
}

func (b *Batch) ETA() time.Time {
	return b.eta
}

func (b *Batch) AllocatedQuantity() int {
	total := 0
	for orderLine := range b.allocations {
		total += orderLine.Qty
	}
	return total
}

func (b *Batch) AvailableQuantity() int {
	return b.purchasedQty - b.AllocatedQuantity()
}

func (b *Batch) CanAllocate(line OrderLine) bool {
	return b.sku == line.SKU && b.AvailableQuantity() >= line.Qty
}

func (b *Batch) Allocate(line OrderLine) {
	if !b.CanAllocate(line) {
		return
	}

	// Much nicer than the linear lookup + lodash version in JS.
	b.allocations[line] = struct{}{}
}

func (b *Batch) Deallocate(line OrderLine) {
	// Much nicer than the linear lookup + lodash version in JS.
	delete(b.allocations, line)
}

func Allocate(line OrderLine, batches []*Batch) (string, error) {
	var validBatches []*Batch
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			validBatches = append(validBatches, batch)
		}
	}

	if len(validBatches) == 0 {
		return "", ErrOutOfStock
	}

	sort.Slice(validBatches, func(i, j int) bool {
		if validBatches[i].ETA().IsZero() {
			return true
		}
		if validBatches[j].ETA().IsZero() {
			return false
		}
		return validBatches[i].ETA().Before(validBatches[j].ETA())
	})

	chosenBatch := validBatches[0]
	chosenBatch.Allocate(line)
	return chosenBatch.Ref(), nil
}
