import { importJWK, SignJWT } from "jose";
import JWKSManager, { ALG } from "./JWKSManager";

export async function jwtGenerator(payload: any): Promise<string> {
  const jwksManager = new JWKSManager();
  const { kid, privateKey } = await jwksManager.getSigningKey();

  const secret = await importJWK(privateKey, ALG);

  console.log("kid", kid);
  const jwt = await new SignJWT(payload)
    .setProtectedHeader({ alg: ALG, kid: kid })
    .setIssuedAt()
    .setExpirationTime("2m")
    .sign(secret);

  console.log("JWT??", jwt);
  return jwt;
}
