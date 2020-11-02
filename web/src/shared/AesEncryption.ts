import sjcl from "sjcl";

//aes with keysize 256, tag size 128 and 100000 iterations for password hashing
const encryptionOpts: sjcl.SjclCipherEncryptParams = {
  iter: 100000,
  ks: 256,
  ts: 128,
  v: 1,
  salt: sjcl.random.randomWords(2,0),
  iv: sjcl.random.randomWords(4,0),
};

const encrypt = (password: string) => (payload: string): string => {
  var cipher = sjcl.encrypt(password, payload, encryptionOpts);
  if (!cipher) {
    return "";
  }
  var cipherJson = JSON.parse(JSON.stringify(cipher));
  var encoded = "data:application/json;base64," + btoa(cipherJson);
  return encoded;
};

const decrypt = (password: string) => (encryptedPayload: string): string => {
  var base64String = encryptedPayload.replace(
    "data:application/json;base64,",
    ""
  );
  var decodedString = atob(base64String);
  return sjcl.decrypt(password, decodedString);
};

export interface AesEncryptor {
  encrypt: (payload: string) => string;
  decrypt: (encryptedPayload: string) => string;
  encryptObject: <T>(payload: T) => string;
  decryptObject: <T>(encryptedPayload: string) => T;
};

export const getAesEncryptor = (password: string): AesEncryptor => {
  const e = encrypt(password);
  const d = decrypt(password);
  return {
    encrypt: e,
    decrypt: d,
    encryptObject: <T>(payload: T): string => e(JSON.stringify(payload)),
    decryptObject: <T>(encryptedPayload: string): T =>
      JSON.parse(d(encryptedPayload))
  };
};
