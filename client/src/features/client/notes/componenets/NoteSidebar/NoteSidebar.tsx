import type React from "react";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Plus, Trash2, FileText, Lock, LockOpen } from "lucide-react";
import { cn } from "@/lib/utils";
// import { useTheme } from "@/hooks/useTheme";

interface Note {
  id: string;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  note_type?: "normal" | "protected" | "self_destructing";
  self_destruct_at?: string | null;
}

interface NoteSidebarProps {
  notes: Note[];
  selectedNoteId: string | null;
  onSelectNote: (id: string) => void;
  onCreateNote: () => Promise<void>;
  onDeleteNote: (id: string) => Promise<void>;
  isLoading: boolean;
  isProtectedMode: boolean;
  onToggleProtectedMode: () => void;
}

export default function NoteSidebar({
  notes,
  selectedNoteId,
  onSelectNote,
  onCreateNote,
  onDeleteNote,
  isLoading,
  isProtectedMode,
  onToggleProtectedMode,
}: Readonly<NoteSidebarProps>) {
  const [searchTerm, setSearchTerm] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);
  // const { theme } = useTheme(); // Declare useTheme variable

  const filteredNotes = notes.filter((note) =>
    note.title.toLowerCase().includes(searchTerm.toLowerCase()),
  );

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setDeletingId(id);
    try {
      await onDeleteNote(id);
    } finally {
      setDeletingId(null);
    }
  };

  return (
    <div className="flex flex-col h-full bg-sidebar border-r border-sidebar-border">
      {/* Header */}
      <div className="border-b border-sidebar-border p-4">
        <Button
          onClick={onCreateNote}
          className="w-full gap-2"
          disabled={isLoading}
          size="sm"
        >
          <Plus className="w-4 h-4" />
          New Note
        </Button>
      </div>

      {/* Search */}
      <div className="p-4 border-b border-sidebar-border">
        <Input
          placeholder="Search notes..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="h-9 text-sm"
        />
      </div>

      {/* Notes List */}
      <ScrollArea className="flex-1">
        <div className="p-2 space-y-1">
          {filteredNotes.length === 0 ? (
            <div className="text-center py-8 px-4">
              <FileText className="w-8 h-8 mx-auto text-sidebar-foreground/40 mb-2" />
              <p className="text-xs text-sidebar-foreground/60 font-medium">
                No notes
              </p>
              <p className="text-xs text-sidebar-foreground/40 mt-1">
                Create one to get started
              </p>
            </div>
          ) : (
            filteredNotes.map((note) => {
              const isExpired =
                note.self_destruct_at &&
                new Date(note.self_destruct_at) < new Date();

              return (
                <button
                  key={note.id}
                  onClick={() => onSelectNote(note.id)}
                  disabled={isExpired as boolean}
                  className={cn(
                    "w-full text-left p-3 rounded-lg text-sm transition-all group hover:bg-sidebar-accent disabled:opacity-50",
                    selectedNoteId === note.id
                      ? "bg-sidebar-primary text-sidebar-primary-foreground shadow-sm"
                      : "text-sidebar-foreground hover:bg-sidebar-accent/50",
                  )}
                >
                  <div className="flex items-center justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">
                        {note.title || "Untitled"}
                      </p>
                      <p className="text-xs opacity-60 truncate line-clamp-1">
                        {note.content?.substring(0, 40).replace(/\n/g, " ") ||
                          "No content"}
                      </p>
                    </div>
                    <button
                      onClick={(e) => handleDelete(note.id, e)}
                      className="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 hover:bg-destructive/20 rounded flex-shrink-0"
                      disabled={deletingId === note.id}
                    >
                      <Trash2 className="w-3 h-3 text-destructive" />
                    </button>
                  </div>
                </button>
              );
            })
          )}
        </div>
      </ScrollArea>

      {/* Protected Mode Toggle */}
      <div className="border-t border-sidebar-border p-3 space-y-2">
        <Button
          onClick={onToggleProtectedMode}
          variant="outline"
          size="sm"
          className={cn(
            "w-full justify-start gap-2",
            isProtectedMode &&
              "bg-destructive/10 border-destructive/30 text-destructive hover:bg-destructive/15",
          )}
        >
          {isProtectedMode ? (
            <>
              <Lock className="w-4 h-4" />
              Protected Mode Active
            </>
          ) : (
            <>
              <LockOpen className="w-4 h-4" />
              Unlock Protected
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
