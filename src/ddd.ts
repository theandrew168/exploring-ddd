import _ from "lodash";

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

		const index = this.allocations.findIndex((l) => _.isEqual(l, line));
		if (index === -1) {
			this.allocations.push(line);
		}
	}

	deallocate(line: OrderLine) {
		const index = this.allocations.findIndex((l) => _.isEqual(l, line));
		if (index !== -1) {
			this.allocations.splice(index, 1);
		}
	}

	canAllocate(line: OrderLine): boolean {
		return this.sku === line.sku && this.availableQuantity >= line.qty;
	}
}
