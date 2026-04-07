// Simple encryption/decryption utility using Web Crypto API
// Note: This is client-side encryption. For production, consider using a proper encryption library

export async function encryptPassword(password: string): Promise<string> {
  // Hash the password using SHA-256 for storage
  const encoder = new TextEncoder();
  const data = encoder.encode(password);
  const hashBuffer = await crypto.subtle.digest("SHA-256", data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
  return hashHex;
}

export async function verifyPassword(
  password: string,
  hash: string,
): Promise<boolean> {
  const newHash = await encryptPassword(password);
  return newHash === hash;
}

export async function encryptContent(
  content: string,
  password: string,
): Promise<string> {
  // Generate a random salt
  const salt = crypto.getRandomValues(new Uint8Array(16));

  // Derive key from password using PBKDF2
  const encoder = new TextEncoder();
  const baseKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(password),
    "PBKDF2",
    false,
    ["deriveBits"],
  );

  const keyBits = await crypto.subtle.deriveBits(
    {
      name: "PBKDF2",
      salt: salt,
      iterations: 100000,
      hash: "SHA-256",
    },
    baseKey,
    256,
  );

  // Encrypt the content using AES-GCM
  const key = await crypto.subtle.importKey("raw", keyBits, "AES-GCM", false, [
    "encrypt",
  ]);
  const iv = crypto.getRandomValues(new Uint8Array(12));

  const encryptedData = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: iv },
    key,
    encoder.encode(content),
  );

  // Combine salt + iv + encryptedData and encode to base64
  const combined = new Uint8Array(
    salt.length + iv.length + encryptedData.byteLength,
  );
  combined.set(salt, 0);
  combined.set(iv, salt.length);
  combined.set(new Uint8Array(encryptedData), salt.length + iv.length);

  return btoa(String.fromCharCode(...Array.from(combined)));
}

export async function decryptContent(
  encryptedContent: string,
  password: string,
): Promise<string> {
  const encoder = new TextEncoder();

  // Decode from base64
  const combined = Uint8Array.from(atob(encryptedContent), (c) =>
    c.charCodeAt(0),
  );

  // Extract salt, iv, and encrypted data
  const salt = combined.slice(0, 16);
  const iv = combined.slice(16, 28);
  const encryptedData = combined.slice(28);

  // Derive key from password
  const baseKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(password),
    "PBKDF2",
    false,
    ["deriveBits"],
  );

  const keyBits = await crypto.subtle.deriveBits(
    {
      name: "PBKDF2",
      salt: salt,
      iterations: 100000,
      hash: "SHA-256",
    },
    baseKey,
    256,
  );

  // Decrypt using AES-GCM
  const key = await crypto.subtle.importKey("raw", keyBits, "AES-GCM", false, [
    "decrypt",
  ]);

  try {
    const decryptedData = await crypto.subtle.decrypt(
      { name: "AES-GCM", iv: iv },
      key,
      encryptedData,
    );

    return new TextDecoder().decode(decryptedData);
  } catch (error) {
    console.error(error);
    throw new Error("Incorrect password or corrupted data");
  }
}
