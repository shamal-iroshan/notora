import type React from "react";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AlertCircle, Eye, EyeOff } from "lucide-react";

interface ProtectedNoteDialogProps {
  isOpen: boolean;
  isUnlocking?: boolean;
  onUnlock: (password: string) => Promise<void>;
  onSetPassword?: (password: string) => Promise<void>;
  mode: "unlock" | "set";
}

export default function ProtectedNoteDialog({
  isOpen,
  onUnlock,
  onSetPassword,
  mode,
}: Readonly<ProtectedNoteDialogProps>) {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      if (mode === "unlock") {
        await onUnlock(password);
      } else {
        if (password !== confirmPassword) {
          setError("Passwords do not match");
          return;
        }
        if (password.length < 6) {
          setError("Password must be at least 6 characters");
          return;
        }
        await onSetPassword?.(password);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to process password",
      );
    } finally {
      setIsLoading(false);
      setPassword("");
      setConfirmPassword("");
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-background border border-border rounded-lg p-6 max-w-sm w-full mx-4 shadow-lg">
        <h2 className="text-lg font-semibold mb-4">
          {mode === "unlock" ? "Unlock Protected Note" : "Set Password"}
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Password</label>
            <div className="relative">
              <Input
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter password"
                disabled={isLoading}
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-2.5 text-muted-foreground hover:text-foreground"
              >
                {showPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>

          {mode === "set" && (
            <div>
              <label className="block text-sm font-medium mb-2">
                Confirm Password
              </label>
              <Input
                type={showPassword ? "text" : "password"}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Confirm password"
                disabled={isLoading}
              />
            </div>
          )}

          {error && (
            <div className="flex gap-2 text-sm text-destructive bg-destructive/10 p-3 rounded">
              <AlertCircle className="w-4 h-4 flex-shrink-0 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          <div className="flex gap-2">
            <Button
              type="submit"
              disabled={
                isLoading || !password || (mode === "set" && !confirmPassword)
              }
              className="flex-1"
            >
              {isLoading
                ? "Processing..."
                : mode === "unlock"
                  ? "Unlock"
                  : "Set Password"}
            </Button>
          </div>

          {mode === "set" && (
            <p className="text-xs text-muted-foreground">
              Remember this password! You'll need to enter it to view this note.
              It cannot be recovered if forgotten.
            </p>
          )}
        </form>
      </div>
    </div>
  );
}
