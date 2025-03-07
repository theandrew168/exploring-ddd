import _ from "lodash";

// Does there even exist a nice way to represent value objects in JS?
// Objects are always compared by reference, so we can't utilize sets.
export type OrderLine = {
	orderId: string;
	sku: string;
	qty: number;
};

export class Batch {
	private ref: string;
	private sku: string;
	private eta?: Date;
	private purchasedQuantity: number;
	private allocations: OrderLine[];

	constructor(ref: string, sku: string, qty: number, eta?: Date) {
		this.ref = ref;
		this.sku = sku;
		this.eta = eta;

		this.purchasedQuantity = qty;

		// It's really painful that JS lacks a way to create "compare by value"
		// objects. If we could, then this could simply be a set. Instead, we let
		// it be an array and have to implement our own comparison logic (via _.isEqual).
		this.allocations = [];
	}

	get allocatedQuantity() {
		return this.allocations.reduce((acc, line) => acc + line.qty, 0);
	}

	get availableQuantity() {
		return this.purchasedQuantity - this.allocatedQuantity;
	}

	allocate(line: OrderLine) {
		if (!this.canAllocate(line)) {
			return;
		}

		// Linear lookups every time... makes me sad.
		const index = this.allocations.findIndex((l) => _.isEqual(l, line));
		if (index === -1) {
			this.allocations.push(line);
		}
	}

	deallocate(line: OrderLine) {
		// Linear lookups every time... makes me sad.
		const index = this.allocations.findIndex((l) => _.isEqual(l, line));
		if (index !== -1) {
			this.allocations.splice(index, 1);
		}
	}

	canAllocate(line: OrderLine): boolean {
		return this.sku === line.sku && this.availableQuantity >= line.qty;
	}
}
