import { toast as sonnerToast } from "sonner";

type ToastVariant = "default" | "destructive" | "success";

interface ToastOptions {
  title?: string;
  description?: string;
  variant?: ToastVariant;
}

export function useToast() {
  function toast({ title, description, variant = "default" }: ToastOptions) {
    const message = title ?? "";

    switch (variant) {
      case "destructive":
        sonnerToast.error(message, {
          description,
        });
        break;

      case "success":
        sonnerToast.success(message, {
          description,
        });
        break;

      default:
        sonnerToast(message, {
          description,
        });
    }
  }

  return { toast };
}
