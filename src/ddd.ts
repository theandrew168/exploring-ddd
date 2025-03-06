export type OrderLine = {
	orderId: string;
	sku: string;
	qty: number;
};

export class Batch {
	ref: string;
	sku: string;
	availableQuantity: number;
	eta?: Date;

	constructor(ref: string, sku: string, qty: number, eta?: Date) {
		this.ref = ref;
		this.sku = sku;
		this.availableQuantity = qty;
		this.eta = eta;
	}

	allocate(line: OrderLine) {
		this.availableQuantity -= line.qty;
	}

	canAllocate(line: OrderLine): boolean {
		return this.sku === line.sku && this.availableQuantity >= line.qty;
	}
}
