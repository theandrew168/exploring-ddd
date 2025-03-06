import { expect, test } from "vitest";

import { Batch, type OrderLine } from "./ddd";

function makeBatchAndLine(sku: string, batchQty: number, lineQty: number): [Batch, OrderLine] {
	return [
		new Batch("batch-001", sku, batchQty, new Date()),
		{ orderId: "order-123", sku, qty: lineQty },
	];
}

test("allocating to a batch reduces available quantity", () => {
	const batch = new Batch("batch-001", "SMALL-TABLE", 20, new Date());
	const line: OrderLine = { orderId: "order-ref", sku: "SMALL-TABLE", qty: 2 };

	batch.allocate(line);

	expect(batch.availableQuantity).toBe(18);
});

test("can allocate if available greater than required", () => {
	const [largeBatch, smallLine] = makeBatchAndLine("ELEGANT-LAMP", 20, 2);
	expect(largeBatch.canAllocate(smallLine)).toBe(true);
});

test("cannot allocate if available smaller than required", () => {
	const [smallBatch, largeLine] = makeBatchAndLine("ELEGANT-LAMP", 2, 20);
	expect(smallBatch.canAllocate(largeLine)).toBe(false);
});

test("can allocate if available equal to required", () => {
	const [batch, line] = makeBatchAndLine("ELEGANT-LAMP", 2, 2);
	expect(batch.canAllocate(line)).toBe(true);
});

test("cannot allocate if sku does not match", () => {
	const batch = new Batch("batch-001", "UNCOMFORTABLE-CHAIR", 100);
	const differentSkuLine: OrderLine = { orderId: "order-123", sku: "EXPENSIVE-TOASTER", qty: 10 };
	expect(batch.canAllocate(differentSkuLine)).toBe(false);
});
