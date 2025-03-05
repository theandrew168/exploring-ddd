import { Hono } from "hono";
import { serve } from "@hono/node-server";

export function add(a: number, b: number): number {
	return a + b;
}

console.log(`2 + 2 = ${add(2, 2)}`);

const app = new Hono();
app.get("/", (c) => c.text("Hello Node.js!\n"));
console.log('Listening on port 3000...');
serve(app);
