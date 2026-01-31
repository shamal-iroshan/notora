import type React from "react";
import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";

import { adminAPI } from "@/lib/api";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { Shield, AlertCircle } from "lucide-react";
import {
  getPathWithRoot,
  PATH_ADMIN_DASHBOARD,
  ROOT_ADMIN,
} from "@/app/router/routes";

export function AdminLoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      const admin = await adminAPI.getCurrentAdmin();
      if (admin) {
        navigate(getPathWithRoot(ROOT_ADMIN, PATH_ADMIN_DASHBOARD), {
          replace: true,
        });
      }
    };

    checkAuth();
  }, [navigate]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const { admin, error } = await adminAPI.adminLogin(email, password);

      if (error) throw new Error(error);
      if (!admin) throw new Error("Login failed");

      navigate(getPathWithRoot(ROOT_ADMIN, PATH_ADMIN_DASHBOARD));
    } catch (error: unknown) {
      setError(error instanceof Error ? error.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted flex items-center justify-center p-4">
      <Card className="w-full max-w-md border-primary/20">
        <CardHeader className="space-y-2">
          <div className="flex items-center justify-center gap-2 mb-2">
            <Shield className="w-6 h-6 text-primary" />
            <h1 className="text-2xl font-bold">Admin Portal</h1>
          </div>
          <CardTitle>Admin Login</CardTitle>
          <CardDescription>
            Access the admin dashboard to manage users and accounts
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleLogin} className="space-y-4">
            {error && (
              <div className="flex gap-2 p-3 rounded-lg bg-destructive/10 border border-destructive/30">
                <AlertCircle className="w-4 h-4 text-destructive flex-shrink-0 mt-0.5" />
                <span className="text-sm text-destructive">{error}</span>
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="admin@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={isLoading}
              />
            </div>

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Signing in..." : "Sign In"}
            </Button>
          </form>

          <div className="mt-6 pt-6 border-t text-center text-sm text-muted-foreground">
            <p className="mb-3">Demo Credentials:</p>
            <p className="font-mono text-xs bg-muted p-2 rounded mb-2">
              Email: admin@example.com
            </p>
            <p className="font-mono text-xs bg-muted p-2 rounded">
              Password: admin123
            </p>
          </div>

          <div className="mt-6 text-center text-sm">
            <Link to="/" className="text-primary hover:underline">
              Back to App
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
