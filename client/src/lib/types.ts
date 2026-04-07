export type TabType = "normal" | "protected";

export interface Note {
  id: string;
  title: string;
  content: string;
  note_type?: "normal" | "protected";
  encrypted_content?: string | null;
  password_hash?: string | null;
  created_at: string;
  updated_at: string;
}

export interface UnlockedNote {
  id: string;
  isUnlocked: boolean;
  decryptedContent?: string;
}
