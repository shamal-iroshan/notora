/* eslint-disable @typescript-eslint/no-unused-vars */
// Mock API Service - Replace with real API calls later

interface User {
  id: string;
  email: string;
  full_name: string | null;
}

interface Note {
  id: string;
  title: string;
  content: string;
  note_type?: "normal" | "protected" | "self_destructing";
  encrypted_content?: string | null;
  password_hash?: string | null;
  self_destruct_at?: string | null;
  created_at: string;
  updated_at: string;
}

interface Profile {
  id: string;
  email: string;
  full_name: string | null;
}

interface UserProfile {
  id: string;
  email: string;
  full_name: string | null;
  status: "pending" | "approved" | "rejected";
  created_at: string;
}

interface Admin {
  id: string;
  email: string;
  full_name: string | null;
}

// Mock user storage
let mockUser: User | null = {
  id: "mock-user-123",
  email: "user@example.com",
  full_name: "John Doe",
};

let mockAdmin: Admin | null = null;

let mockAllUsers: UserProfile[] = [
  {
    id: "mock-user-123",
    email: "user@example.com",
    full_name: "John Doe",
    status: "approved",
    created_at: new Date().toISOString(),
  },
];

let mockNotes: Note[] = [
  {
    id: "note-1",
    title: "Welcome to MarkNotes",
    content: "# Welcome!\n\nStart typing your markdown notes here.",
    note_type: "normal",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: "note-2",
    title: "My Protected Note",
    content: "",
    note_type: "protected",
    encrypted_content: "mock-encrypted-data",
    password_hash: "mock-hash",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
];

// Auth API calls
export const authAPI = {
  async login(email: string, password: string) {
    await new Promise((resolve) => setTimeout(resolve, 500));

    // Check if user exists and is approved
    const user = mockAllUsers.find((u) => u.email === email);
    if (!user) {
      return { user: null, error: "User not found" };
    }

    if (user.status !== "approved") {
      return { user: null, error: "Account pending approval from admin" };
    }

    if (email === "user@example.com" && password === "password123") {
      mockUser = {
        id: user.id,
        email: user.email,
        full_name: user.full_name,
      };
      return { user: mockUser, error: null };
    }

    return { user: null, error: "Invalid credentials" };
  },

  async signup(email: string, password: string, fullName: string) {
    await new Promise((resolve) => setTimeout(resolve, 600));

    // Check if user already exists
    const existingUser = mockAllUsers.find((u) => u.email === email);
    if (existingUser) {
      return { user: null, error: "Email already registered" };
    }

    // Create user with pending status
    const newUser: UserProfile = {
      id: `user-${Date.now()}`,
      email,
      full_name: fullName,
      status: "pending",
      created_at: new Date().toISOString(),
    };

    mockAllUsers.push(newUser);

    return {
      user: {
        id: newUser.id,
        email: newUser.email,
        full_name: newUser.full_name,
      },
      error: null,
    };
  },

  // async signup(email: string, password: string, fullName: string) {
  //   // Simulate API delay
  //   await new Promise((resolve) => setTimeout(resolve, 800));

  //   if (password.length < 8) {
  //     return { error: "Password must be at least 8 characters" };
  //   }

  //   mockUser = {
  //     id: `user-${Date.now()}`,
  //     email,
  //     full_name: fullName,
  //   };

  //   return { user: mockUser, error: null };
  // },

  async logout() {
    await new Promise((resolve) => setTimeout(resolve, 300));
    mockUser = null;
    return { error: null };
  },

  async getCurrentUser() {
    await new Promise((resolve) => setTimeout(resolve, 100));
    return mockUser;
  },
};

// Notes API calls
export const notesAPI = {
  async getNotes(userId: string) {
    await new Promise((resolve) => setTimeout(resolve, 300));
    return mockNotes.filter((note) => true); // In real app, filter by user
  },

  async createNote(
    userId: string,
    title: string,
    noteType: "normal" | "protected" | "self_destructing" = "normal",
  ) {
    await new Promise((resolve) => setTimeout(resolve, 200));

    const newNote: Note = {
      id: `note-${Date.now()}`,
      title,
      content: "",
      note_type: noteType,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    if (noteType === "self_destructing") {
      const expirationTime = new Date();
      expirationTime.setHours(expirationTime.getHours() + 24);
      newNote.self_destruct_at = expirationTime.toISOString();
    }

    if (noteType === "protected") {
      newNote.encrypted_content = "";
      newNote.password_hash = null;
    }

    mockNotes.unshift(newNote);
    return newNote;
  },

  async updateNote(noteId: string, updates: Partial<Note>) {
    await new Promise((resolve) => setTimeout(resolve, 200));

    const index = mockNotes.findIndex((n) => n.id === noteId);
    if (index === -1) {
      throw new Error("Note not found");
    }

    mockNotes[index] = {
      ...mockNotes[index],
      ...updates,
      updated_at: new Date().toISOString(),
    };

    return mockNotes[index];
  },

  async deleteNote(noteId: string) {
    await new Promise((resolve) => setTimeout(resolve, 200));

    mockNotes = mockNotes.filter((n) => n.id !== noteId);
    return { error: null };
  },

  async setProtectedPassword(noteId: string, passwordHash: string) {
    await new Promise((resolve) => setTimeout(resolve, 150));

    const note = mockNotes.find((n) => n.id === noteId);
    if (!note) throw new Error("Note not found");

    note.password_hash = passwordHash;
    return { error: null };
  },

  async updateSelfDestruct(noteId: string, expirationTime: string) {
    await new Promise((resolve) => setTimeout(resolve, 150));

    const note = mockNotes.find((n) => n.id === noteId);
    if (!note) throw new Error("Note not found");

    note.self_destruct_at = expirationTime;
    return { error: null };
  },
};

// Profile API calls
export const profileAPI = {
  async getProfile(userId: string) {
    await new Promise((resolve) => setTimeout(resolve, 200));

    if (!mockUser) {
      throw new Error("User not found");
    }

    const profile: Profile = {
      id: mockUser.id,
      email: mockUser.email,
      full_name: mockUser.full_name,
    };

    return profile;
  },

  async updateProfile(userId: string, updates: { full_name?: string }) {
    await new Promise((resolve) => setTimeout(resolve, 250));

    if (!mockUser) {
      throw new Error("User not found");
    }

    if (updates.full_name !== undefined) {
      mockUser.full_name = updates.full_name;
    }

    return { error: null };
  },

  async changePassword(userId: string, newPassword: string) {
    await new Promise((resolve) => setTimeout(resolve, 300));

    if (!mockUser) {
      throw new Error("User not found");
    }

    return { error: null };
  },
};

// Admin API calls
export const adminAPI = {
  async adminLogin(email: string, password: string) {
    await new Promise((resolve) => setTimeout(resolve, 500));

    // Mock admin credentials
    if (email === "admin@example.com" && password === "admin123") {
      mockAdmin = {
        id: "admin-001",
        email: "admin@example.com",
        full_name: "Admin User",
      };
      return { admin: mockAdmin, error: null };
    }

    return { admin: null, error: "Invalid credentials" };
  },

  async adminLogout() {
    await new Promise((resolve) => setTimeout(resolve, 200));
    mockAdmin = null;
    return { error: null };
  },

  async getCurrentAdmin() {
    await new Promise((resolve) => setTimeout(resolve, 100));
    return mockAdmin;
  },

  async getAllUsers() {
    await new Promise((resolve) => setTimeout(resolve, 300));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    return mockAllUsers;
  },

  async approveUser(userId: string) {
    await new Promise((resolve) => setTimeout(resolve, 250));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    const user = mockAllUsers.find((u) => u.id === userId);
    if (!user) throw new Error("User not found");

    user.status = "approved";
    return { error: null };
  },

  async rejectUser(userId: string) {
    await new Promise((resolve) => setTimeout(resolve, 250));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    const user = mockAllUsers.find((u) => u.id === userId);
    if (!user) throw new Error("User not found");

    user.status = "rejected";
    return { error: null };
  },

  async createUserProfile(email: string, fullName: string, password: string) {
    await new Promise((resolve) => setTimeout(resolve, 400));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    const newUser: UserProfile = {
      id: `user-${Date.now()}`,
      email,
      full_name: fullName,
      status: "approved",
      created_at: new Date().toISOString(),
    };

    mockAllUsers.push(newUser);
    return { user: newUser, error: null };
  },

  async changeUserPassword(userId: string, newPassword: string) {
    await new Promise((resolve) => setTimeout(resolve, 300));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    const user = mockAllUsers.find((u) => u.id === userId);
    if (!user) throw new Error("User not found");

    return { error: null };
  },

  async deleteUser(userId: string) {
    await new Promise((resolve) => setTimeout(resolve, 300));

    if (!mockAdmin) {
      throw new Error("Not authenticated as admin");
    }

    mockAllUsers = mockAllUsers.filter((u) => u.id !== userId);
    return { error: null };
  },
};
