import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

import dotenv from 'dotenv'
dotenv.config();

const connectionString = process.env.DATABASE_URL;

if (!connectionString) {
  throw new Error("DATABASE_URL is not set in environment");
}

const client = postgres(connectionString, {
  prepare: false,
});

export const db = drizzle(client);

const verifyDbConnection = async () => {
  await client`select 1`;
};

export { client, verifyDbConnection };