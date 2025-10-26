import { knexClient } from "./knexClient";

type KeyType = "current-key" | "stale-key";

const TABLE_NAME = "jwks_key";
class CacheManager {
  private client = knexClient;

  async setup() {
    if (!(await this.client.schema.hasTable(TABLE_NAME))) {
      await this.client.schema.createTable(TABLE_NAME, (table) => {
        table.increments("id").primary();
        table.string("key_type").unique();
        table.text("value");
        table.timestamps(true, true);
      });
    }
  }

  async getKey(keyType: KeyType): Promise<any | null> {
    try {
      const jwksKey = await this.client(TABLE_NAME)
        .where({ key_type: keyType })
        .first();
      return jwksKey ? JSON.parse(jwksKey.value) : null;
    } catch (error) {
      throw new Error(`Error while getting key: ${error}`);
    }
  }

  async setKey(keyType: KeyType, value: any): Promise<void> {
    try {
      await this.client(TABLE_NAME)
        .insert({
          key_type: keyType,
          value: JSON.stringify(value),
        })
        .onConflict("key_type")
        .merge();
    } catch (error) {
      throw new Error(`Error while setting key: ${error}`);
    }
  }

  async deleteStaleKey(): Promise<void> {
    try {
      this.client("jwks_key").where({ key_type: "stale-key" }).del();
    } catch (error) {
      throw new Error(`Error while deleting stale key: ${error}`);
    }
  }
}

export default CacheManager;
