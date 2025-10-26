import { JWK } from "jose";

export type CacheableJWKSKey = {
  timestamp: number;
  publicKey: JWK;
  privateKey: JWK;
};
