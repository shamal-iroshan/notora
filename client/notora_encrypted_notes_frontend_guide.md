# NOTORA Encrypted Notes – Frontend Implementation Guide

This guide describes how to implement **end‑to‑end encrypted notes (E2EE)** in the NOTORA frontend.  
The backend never sees plaintext — all encryption and decryption happen exclusively in the browser.

---

## 🧠 1. Overview of Encrypted Notes Architecture

NOTORA uses a **zero‑knowledge encryption model**:

- The **user enters a master password**
- Backend provides a **user_salt**
- Frontend derives a **masterKey = PBKDF2(masterPassword, userSalt)**
- Every encrypted note uses:
  - a **noteSalt** (random per note)
  - a **noteKey = HMAC(masterKey, noteSalt)**
  - AES‑GCM to encrypt title + content

Backend stores **only ciphertext**, **nonces**, and **salts**.

Backend **never knows**:
- the master password  
- masterKey  
- per‑note keys  
- decrypted text  

---

## 🔐 2. Master Key Derivation (PBKDF2)

Used after login when frontend receives `user_salt`.

```ts
export async function deriveMasterKey(password: string, userSalt: string): Promise<CryptoKey> {
  const encoder = new TextEncoder();
  const baseKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(password),
    { name: "PBKDF2" },
    false,
    ["deriveKey"]
  );

  return crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      salt: Uint8Array.from(Buffer.from(userSalt, "hex")),
      iterations: 250000,
      hash: "SHA-256",
    },
    baseKey,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}
```

---

## 🧂 3. Per‑Note Key Derivation

Each note uses its own **note_salt**.

```ts
export async function deriveNoteKey(masterKey: CryptoKey, noteSalt: string): Promise<CryptoKey> {
  const saltBytes = Uint8Array.from(Buffer.from(noteSalt, "hex"));
  return crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      salt: saltBytes,
      iterations: 100000,
      hash: "SHA-256",
    },
    masterKey,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"]
  );
}
```

---

## 🔒 4. Encrypting a Note (AES‑GCM)

```ts
export async function encryptNote(
  masterKey: CryptoKey,
  title: string,
  content: string
) {
  const encoder = new TextEncoder();

  // Generate per-note salt
  const noteSalt = crypto.getRandomValues(new Uint8Array(16));
  const noteSaltHex = Buffer.from(noteSalt).toString("hex");

  const noteKey = await deriveNoteKey(masterKey, noteSaltHex);

  const titleNonce = crypto.getRandomValues(new Uint8Array(12));
  const contentNonce = crypto.getRandomValues(new Uint8Array(12));

  const encryptedTitle = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: titleNonce },
    noteKey,
    encoder.encode(title)
  );

  const encryptedContent = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: contentNonce },
    noteKey,
    encoder.encode(content)
  );

  return {
    title_ciphertext: Buffer.from(new Uint8Array(encryptedTitle)).toString("base64"),
    content_ciphertext: Buffer.from(new Uint8Array(encryptedContent)).toString("base64"),
    title_nonce: Buffer.from(titleNonce).toString("hex"),
    content_nonce: Buffer.from(contentNonce).toString("hex"),
    note_salt: noteSaltHex,
  };
}
```

---

## 🔓 5. Decrypting a Note (AES‑GCM)

```ts
export async function decryptNote(
  masterKey: CryptoKey,
  note: {
    title: string;
    content: string;
    title_nonce: string;
    content_nonce: string;
    note_salt: string;
  }
) {
  const decoder = new TextDecoder();

  const noteKey = await deriveNoteKey(masterKey, note.note_salt);

  const decryptedTitle = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: Uint8Array.from(Buffer.from(note.title_nonce, "hex")),
    },
    noteKey,
    Uint8Array.from(Buffer.from(note.title, "base64"))
  );

  const decryptedContent = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: Uint8Array.from(Buffer.from(note.content_nonce, "hex")),
    },
    noteKey,
    Uint8Array.from(Buffer.from(note.content, "base64"))
  );

  return {
    title: decoder.decode(decryptedTitle),
    content: decoder.decode(decryptedContent),
  };
}
```

---

## 🧩 6. Unlock Flow for Users

When a user logs in:

1. Backend returns:
   - `access_token` (cookie)
   - `refresh_token` (cookie)
   - `user_salt`
2. Frontend asks user: **Enter master password**
3. Derive master key:

```ts
const masterKey = await deriveMasterKey(masterPassword, userSalt);
```

4. Store master key **in memory**, never in localStorage.

---

## 📝 7. Creating an Encrypted Note

```ts
const encrypted = await encryptNote(masterKey, title, content);

await api.post("/api/encrypted-notes", encrypted);
```

---

## 📥 8. Fetching Notes List (Metadata)

API returns only ciphertext title + salt.  
Frontend decrypts title to display.

---

## 🧱 9. Decrypting a Full Note for Viewing

```ts
const note = await api.get(`/api/encrypted-notes/${id}`);
const decrypted = await decryptNote(masterKey, note.data);
```

---

## 🧼 10. Clearing Encryption Keys on Logout

```ts
masterKey = null;
window.location.reload();
```

---

## 🛡️ 11. Security Best Practices

- Never store the master password
- Never store the master key on disk
- Use PBKDF2 with at least 250k iterations
- Each note must have a unique **note_salt**
- Never send plaintext to backend
- Decrypt only when required
- Use HTTPS always

