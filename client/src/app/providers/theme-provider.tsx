import { useEffect, useState } from "react";
import type React from "react";
import { ThemeContext, type Theme } from "../context/theme-context";

interface ThemeProviderProps {
  children: React.ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  // Initialize from localStorage synchronously
  const [theme, setTheme] = useState<Theme>(() => {
    const stored = localStorage.getItem("theme") as Theme | null;
    return stored ?? "system";
  });

  useEffect(() => {
    const root = document.documentElement;

    let effectiveTheme: "light" | "dark";

    if (theme === "system") {
      effectiveTheme = window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light";
    } else {
      effectiveTheme = theme;
    }

    root.classList.toggle("dark", effectiveTheme === "dark");

    localStorage.setItem("theme", theme);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}
