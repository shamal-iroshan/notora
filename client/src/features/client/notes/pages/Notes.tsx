/* eslint-disable @typescript-eslint/no-unused-vars */
import { useEffect, useState, useCallback } from "react";
import { useNavigate, Link } from "react-router-dom";

import { authAPI, notesAPI, type Note } from "@/lib/api";
import {
  encryptPassword,
  decryptContent,
  encryptContent,
  verifyPassword,
} from "@/lib/encryption";

import { Button } from "@/components/ui/button";

import { LogOut, Settings, Menu, X } from "lucide-react";

// import { PWAInstallButton } from "@/components/pwa-install-button";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useTheme } from "@/hooks/useTheme";
import NoteEditor from "../componenets/NoteEditor";
import NoteSidebar from "../componenets/NoteSidebar";
import ProtectedNoteDialog from "../componenets/ProtectedNoteDialog";
import ProtectedModeUnlock from "../componenets/ProtectedModeUnlock";

interface UnlockedNote {
  id: string;
  isUnlocked: boolean;
  decryptedContent?: string;
}

export function NotesPage() {
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();

  const [notes, setNotes] = useState<Note[]>([]);
  const [selectedNoteId, setSelectedNoteId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState<{
    id: string;
    email: string;
    full_name: string | null;
  } | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const [unlockedNotes, setUnlockedNotes] = useState<Map<string, UnlockedNote>>(
    new Map(),
  );

  const [protectedDialog, setProtectedDialog] = useState({
    isOpen: false,
    mode: "unlock" as "unlock" | "set",
  });

  const [isProtectedMode, setIsProtectedMode] = useState(false);
  const [showUnlockDialog, setShowUnlockDialog] = useState(false);
  const [showLockConfirmation, setShowLockConfirmation] = useState(false);

  const [isProtectedLocked, setIsProtectedLocked] = useState(true);
  const [protectedPassword, setProtectedPassword] = useState("");

  // ✅ Load user + notes
  useEffect(() => {
    const loadUserAndNotes = async () => {
      try {
        const currentUser = await authAPI.getCurrentUser();

        if (!currentUser) {
          navigate("/auth/login", { replace: true });
          return;
        }

        setUser(currentUser);

        const notesData = await notesAPI.getNotes(currentUser.id);
        setNotes(notesData);
      } catch (error) {
        console.error("Error loading notes:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadUserAndNotes();
  }, [navigate]);

  // ✅ Filter notes
  const filteredNotes = notes.filter((n) => {
    const noteType = n.note_type || "normal";
    return isProtectedMode ? noteType === "protected" : noteType === "normal";
  });

  // ✅ Select first note
  useEffect(() => {
    if (filteredNotes.length > 0 && !selectedNoteId) {
      setSelectedNoteId(filteredNotes[0].id);
    } else if (filteredNotes.length === 0) {
      setSelectedNoteId(null);
    }
  }, [filteredNotes, selectedNoteId]);

  // ✅ Create note
  const handleCreateNote = useCallback(
    async (type: "normal" | "protected" = "normal") => {
      if (!user) return;

      try {
        const newNote = await notesAPI.createNote(
          user.id,
          "Untitled Note",
          type,
        );

        setNotes((prev) => [newNote, ...prev]);
        setSelectedNoteId(newNote.id);

        if (type === "protected") {
          setProtectedDialog({ isOpen: true, mode: "set" });
        }
      } catch (error) {
        console.error("Error creating note:", error);
      }
    },
    [user],
  );

  // ✅ Delete note
  const handleDeleteNote = useCallback(
    async (id: string) => {
      try {
        await notesAPI.deleteNote(id);

        setNotes((prev) => prev.filter((note) => note.id !== id));

        setUnlockedNotes((prev) => {
          const updated = new Map(prev);
          updated.delete(id);
          return updated;
        });

        if (selectedNoteId === id) {
          const remaining = filteredNotes.filter((note) => note.id !== id);
          setSelectedNoteId(remaining.length > 0 ? remaining[0].id : null);
        }
      } catch (error) {
        console.error("Error deleting note:", error);
      }
    },
    [selectedNoteId, filteredNotes],
  );

  // ✅ Save note
  const handleSaveNote = useCallback(
    async (title: string, content: string) => {
      if (!selectedNoteId || !user) return;

      try {
        const selectedNote = notes.find((n) => n.id === selectedNoteId);
        if (!selectedNote) return;

        const updateData: Partial<Note> = { title };

        if (selectedNote.note_type === "protected") {
          const unlocked = unlockedNotes.get(selectedNoteId);
          if (!unlocked?.isUnlocked) return;

          const encrypted = await encryptContent(
            content,
            unlocked.decryptedContent || "",
          );
          updateData.encrypted_content = encrypted;
        } else {
          updateData.content = content;
        }

        await notesAPI.updateNote(selectedNoteId, updateData);

        setNotes((prev) =>
          prev.map((note) =>
            note.id === selectedNoteId
              ? { ...note, title, updated_at: new Date().toISOString() }
              : note,
          ),
        );
      } catch (error) {
        console.error("Error saving note:", error);
      }
    },
    [selectedNoteId, user, notes, unlockedNotes],
  );

  // ✅ Protected mode toggle
  const handleToggleProtectedMode = () => {
    if (isProtectedMode) {
      setShowLockConfirmation(true);
    } else {
      setShowUnlockDialog(true);
    }
  };

  const handleUnlockProtectedMode = async (password: string) => {
    if (password.trim().length > 0) {
      setIsProtectedMode(true);
      setShowUnlockDialog(false);
      setSelectedNoteId(null);
    }
  };

  const handleLockProtectedMode = () => {
    setIsProtectedMode(false);
    setShowLockConfirmation(false);
    setSelectedNoteId(null);
    setUnlockedNotes(new Map());
  };

  const handleUnlockProtectedNotes = useCallback(
    async (password: string) => {
      if (filteredNotes.length === 0) return;

      try {
        const firstNote = filteredNotes[0];
        if (!firstNote || firstNote.note_type !== "protected") return;

        const isCorrect = await verifyPassword(
          password,
          firstNote.password_hash || "",
        );
        if (!isCorrect) throw new Error("Incorrect password");

        setProtectedPassword(password);
        setShowUnlockDialog(false);

        for (const note of filteredNotes) {
          if (note.note_type === "protected") {
            try {
              const decrypted = await decryptContent(
                note.encrypted_content || "",
                password,
              );

              setUnlockedNotes((prev) =>
                new Map(prev).set(note.id, {
                  id: note.id,
                  isUnlocked: true,
                  decryptedContent: decrypted,
                }),
              );
            } catch (err) {
              console.error(`Failed to decrypt note ${note.id}:`, err);
            }
          }
        }

        setIsProtectedLocked(false);
      } catch (error) {
        console.error("Error unlocking notes:", error);
      }
    },
    [filteredNotes],
  );

  const handleSetProtectedPassword = useCallback(
    async (password: string) => {
      if (!selectedNoteId || !user) return;

      try {
        const passwordHash = await encryptPassword(password);

        await notesAPI.setProtectedPassword(selectedNoteId, passwordHash);

        setUnlockedNotes((prev) =>
          new Map(prev).set(selectedNoteId, {
            id: selectedNoteId,
            isUnlocked: true,
            decryptedContent: password,
          }),
        );

        setProtectedPassword(password);
        setProtectedDialog({ isOpen: false, mode: "unlock" });
        setIsProtectedLocked(false);
      } catch (error) {
        console.error("Error setting password:", error);
      }
    },
    [selectedNoteId, user],
  );

  const handleLogout = async () => {
    await authAPI.logout();
    navigate("/auth/login", { replace: true });
  };

  const selectedNote = notes.find((note) => note.id === selectedNoteId);

  const noteContent =
    selectedNote && unlockedNotes.get(selectedNoteId || "")?.isUnlocked
      ? unlockedNotes.get(selectedNoteId || "")?.decryptedContent ||
        selectedNote.content
      : selectedNote?.content || "";

  return (
    <div className="flex h-screen bg-background overflow-hidden">
      {/* SIDEBAR */}
      <div
        className={`${sidebarOpen ? "w-64" : "w-0"} hidden md:flex flex-col border-r transition-all`}
      >
        <NoteSidebar
          notes={filteredNotes}
          selectedNoteId={selectedNoteId}
          onSelectNote={setSelectedNoteId}
          onCreateNote={() =>
            handleCreateNote(isProtectedMode ? "protected" : "normal")
          }
          onDeleteNote={handleDeleteNote}
          isLoading={isLoading}
          isProtectedMode={isProtectedMode}
          onToggleProtectedMode={handleToggleProtectedMode}
        />
      </div>

      {/* MAIN */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* TOP BAR */}
        <div className="border-b px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setSidebarOpen(!sidebarOpen)}
            >
              {sidebarOpen ? (
                <X className="w-4 h-4" />
              ) : (
                <Menu className="w-4 h-4" />
              )}
            </Button>

            <div>
              <h1 className="text-lg font-semibold">MarkNotes</h1>
              {user && (
                <p className="text-xs text-muted-foreground">{user.email}</p>
              )}
            </div>
          </div>

          <div className="flex items-center gap-2">
            {/* <PWAInstallButton /> */}

            <Button
              variant="ghost"
              size="icon"
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            >
              {theme === "dark" ? "☀️" : "🌙"}
            </Button>

            <Link to="/protected/account">
              <Button variant="ghost" size="icon">
                <Settings className="w-4 h-4" />
              </Button>
            </Link>

            <Button variant="ghost" size="icon" onClick={handleLogout}>
              <LogOut className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* EDITOR */}
        <div className="flex-1 overflow-hidden flex flex-col">
          {selectedNote ? (
            <NoteEditor
              note={{ ...selectedNote, content: noteContent }}
              onSave={handleSaveNote}
              noteType={selectedNote.note_type}
            />
          ) : (
            <div className="flex items-center justify-center h-full">
              <Button
                onClick={() =>
                  handleCreateNote(isProtectedMode ? "protected" : "normal")
                }
              >
                Create Note
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* DIALOGS */}
      <ProtectedNoteDialog
        isOpen={protectedDialog.isOpen}
        mode={protectedDialog.mode}
        onUnlock={handleUnlockProtectedNotes}
        onSetPassword={handleSetProtectedPassword}
      />

      <ProtectedModeUnlock
        isOpen={showUnlockDialog}
        onUnlock={handleUnlockProtectedMode}
        onCancel={() => setShowUnlockDialog(false)}
      />

      <AlertDialog
        open={showLockConfirmation}
        onOpenChange={setShowLockConfirmation}
      >
        <AlertDialogContent>
          <AlertDialogTitle>Exit Protected Mode?</AlertDialogTitle>
          <AlertDialogDescription>
            You're about to exit protected mode.
          </AlertDialogDescription>

          <div className="flex gap-2 justify-end">
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleLockProtectedMode}>
              Exit Protected Mode
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
