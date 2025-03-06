import { Hono } from "hono";
import { serve } from "@hono/node-server";

const app = new Hono();
app.get("/", (c) => c.text("Hello Node.js!\n"));
console.log('Listening on port 3000...');
serve(app);
