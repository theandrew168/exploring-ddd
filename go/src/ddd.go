package ddd

import "time"

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

type OrderLineSet map[OrderLine]struct{}

type Batch struct {
	ref string
	sku string
	eta time.Time

	purchasedQty int
	allocations  OrderLineSet
}

func NewBatch(ref, sku string, qty int, eta time.Time) *Batch {
	b := Batch{
		ref: ref,
		sku: sku,
		eta: eta,

		purchasedQty: qty,
		allocations:  make(OrderLineSet),
	}
	return &b
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

	b.allocations[line] = struct{}{}
}

func (b *Batch) Deallocate(line OrderLine) {
	delete(b.allocations, line)
}
