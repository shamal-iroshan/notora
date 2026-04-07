import { useState } from "react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Input } from "@/components/ui/input";
import { Lock } from "lucide-react";

interface ProtectedModeUnlockProps {
  isOpen: boolean;
  onUnlock: (password: string) => void;
  onCancel: () => void;
}

export default function ProtectedModeUnlock({
  isOpen,
  onUnlock,
  onCancel,
}: Readonly<ProtectedModeUnlockProps>) {
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleUnlock = () => {
    if (!password.trim()) {
      setError("Please enter a password");
      return;
    }
    setError("");
    onUnlock(password);
    setPassword("");
  };

  const handleCancel = () => {
    setPassword("");
    setError("");
    onCancel();
  };

  return (
    <AlertDialog open={isOpen} onOpenChange={(open) => !open && handleCancel()}>
      <AlertDialogContent>
        <AlertDialogTitle className="flex items-center gap-2">
          <Lock className="w-5 h-5" />
          Unlock Protected Mode
        </AlertDialogTitle>
        <AlertDialogDescription>
          Enter a password to access your protected notes. You'll need to
          remember this password.
        </AlertDialogDescription>
        <div className="space-y-3">
          <Input
            type="password"
            placeholder="Enter password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
              setError("");
            }}
            onKeyDown={(e) => e.key === "Enter" && handleUnlock()}
            className="mt-2"
            autoFocus
          />
          {error && <p className="text-sm text-destructive">{error}</p>}
        </div>
        <div className="flex gap-2 justify-end">
          <AlertDialogCancel onClick={handleCancel}>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={handleUnlock}>Unlock</AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}
