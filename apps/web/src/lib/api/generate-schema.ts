import { env } from "@/env/client";
import { exec } from "node:child_process";
import { promisify } from "node:util";

const execPromise = promisify(exec);

async function generateSchema(): Promise<void> {
  const url = `${env.VITE_API_URL}/${env.VITE_API_VERSION}/openapi.yml`;
  const output = "./src/lib/api/schema.d.ts";
  try {
    console.error(`Generating schema for ${url}...`);
    await execPromise(`bunx --bun openapi-typescript ${url} -o ${output}`);
    console.error(`Successfully generated schema for ${url}`);
  } catch (error) {
    console.error(`Error generating schema for ${url}`);
    throw error;
  }
}

generateSchema();
