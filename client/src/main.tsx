import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { ThemeProvider } from "./app/providers/theme-provider";
import { RouterProvider } from "react-router-dom";
import { router } from "./app/router/router";
import { Toaster } from "sonner";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider>
      <Toaster richColors position="top-right" />
      <RouterProvider router={router} />
    </ThemeProvider>
  </StrictMode>,
);
