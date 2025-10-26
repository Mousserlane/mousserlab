import { calculateJwkThumbprint, exportJWK, generateKeyPair } from "jose";

import CacheManager from "./CacheManager";
import { CacheableJWKSKey } from "./types";

const EXPIRATION_TTL = 3600 * 1000; // expiration 1 hour for both key
const ALG = "ES256";

class JWKSManager {
  private cacheManager = new CacheManager();

  async getPublicJWKSKey() {
    let currentKey: CacheableJWKSKey =
      await this.cacheManager.getKey("current-key");

    if (!Boolean(currentKey) || this.isKeyExpired(currentKey)) {
      currentKey = await this.rotateKeys(currentKey);
    }

    const staleKey = await this.cacheManager.getKey("stale-key");
    const { timestamp, ...currentKeyRest } = currentKey; // We do not want timestamp in the JSON;

    const keys = [currentKeyRest.publicKey];

    if (staleKey && !this.isKeyExpired(staleKey)) {
      const { timestamp, ...staleKeyRest } = staleKey;
      keys.push(staleKeyRest.publicKey);
    } else if (staleKey && this.isKeyExpired(staleKey)) {
      await this.cacheManager.deleteStaleKey();
    }

    return { keys };
  }

  async rotateKeys(oldKey: CacheableJWKSKey | null): Promise<CacheableJWKSKey> {
    const newKey = await this.generateAsymmetricKey();

    // if has old key, set it to stale
    if (Boolean(oldKey)) {
      await this.cacheManager.setKey("stale-key", {
        ...oldKey,
        timestamp: Date.now(), // Keep stale key for 1 hour
      });
    }

    const serializedJWKS = await this.serializeJWKS(newKey);
    const currentKey = {
      ...serializedJWKS,
      timestamp: Date.now(),
    };

    await this.cacheManager.setKey("current-key", currentKey);
    return currentKey;
  }

  async generateAsymmetricKey() {
    // will generate the key-pair
    const { publicKey, privateKey } = await generateKeyPair(ALG, {
      extractable: true,
    });

    return { publicKey, privateKey };
  }

  private isKeyExpired(key: CacheableJWKSKey): boolean {
    return Date.now() - key.timestamp > EXPIRATION_TTL;
  }

  // To serialize the CryptoKey object to a plain JS object so it can be stored in-memory.
  // This is because CryptoKey object cannot be saved to memory.
  private async serializeJWKS(jwksKeys: {
    publicKey: CryptoKey;
    privateKey: CryptoKey;
  }) {
    const publicKey = await exportJWK(jwksKeys.publicKey);
    const privateKey = await exportJWK(jwksKeys.privateKey);

    publicKey.use = "sig";
    publicKey.alg = ALG;
    publicKey.kid = await calculateJwkThumbprint(publicKey);

    return { publicKey, privateKey };
  }
}

export default JWKSManager;
