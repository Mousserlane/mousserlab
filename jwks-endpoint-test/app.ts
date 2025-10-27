import express from "express";
import JWKSManager from "./JWKSManager";
import CacheManager from "./CacheManager";
import { jwtGenerator } from "./jwtGenerator";

const PORT = 5888;
const app = express();

const cacheManager = new CacheManager();

cacheManager
  .setup()
  .then(() => {
    console.log("Cache setup successful");
  })
  .catch((error) => console.error("Cache setup error: ", error));

app.get("/", (req, res) => {
  res.send("<h1>Test page</h1>");
});

app.get("/jwks", async (_, res) => {
  const jwksManager = new JWKSManager();
  const jwksPub = await jwksManager.getPublicJWKSKey();
  res.json(jwksPub);
});

app.get("/generateJWT", async (_, res) => {
  const payload = {
    name: "YOUR NAME HERE",
    address: "SOMEWHERE",
  };

  const jwt = await jwtGenerator(payload);
  res.json({ jwt });
});

app.listen(PORT, () => {
  console.log(`App is listeneing on port ${PORT}`);
});
