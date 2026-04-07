import type React from "react";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { authAPI, profileAPI } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { AlertCircle, ArrowLeft } from "lucide-react";
import {
  getPathWithRoot,
  PATH_CLIENT_LOGIN,
  ROOT_INDEX,
} from "@/app/router/routes";

interface Profile {
  id: string;
  email: string;
  full_name: string | null;
}

export function AccountPage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [fullName, setFullName] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const navigate = useNavigate();

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const user = await authAPI.getCurrentUser();

        if (!user) {
          navigate(getPathWithRoot(ROOT_INDEX, PATH_CLIENT_LOGIN), {
            replace: true,
          });
          return;
        }

        const profileData = await profileAPI.getProfile(user.id);

        setProfile(profileData);
        setFullName(profileData.full_name || "");
      } catch (error) {
        console.error("Error loading profile:", error);
        setMessage({
          type: "error",
          text: "Failed to load profile",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadProfile();
  }, [navigate]);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);

    try {
      const user = await authAPI.getCurrentUser();
      if (!user) throw new Error("Not authenticated");

      const { error } = await profileAPI.updateProfile(user.id, {
        full_name: fullName,
      });

      if (error) throw new Error(error);

      setMessage({
        type: "success",
        text: "Profile updated successfully",
      });

      setTimeout(() => setMessage(null), 3000);
    } catch (error) {
      console.error("Error updating profile:", error);
      setMessage({
        type: "error",
        text: "Failed to update profile",
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleChangePassword = async () => {
    const newPassword = prompt("Enter new password (minimum 8 characters):");

    if (!newPassword || newPassword.length < 8) {
      setMessage({
        type: "error",
        text: "Password must be at least 8 characters",
      });
      return;
    }

    try {
      const user = await authAPI.getCurrentUser();
      if (!user) throw new Error("Not authenticated");

      const { error } = await profileAPI.changePassword(user.id, newPassword);

      if (error) throw new Error(error);

      setMessage({
        type: "success",
        text: "Password updated successfully",
      });

      setTimeout(() => setMessage(null), 3000);
    } catch (error) {
      console.error("Error changing password:", error);
      setMessage({
        type: "error",
        text: "Failed to change password",
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-2xl mx-auto py-8 px-6">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <Link to="/protected/notes">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="w-4 h-4" />
            </Button>
          </Link>

          <div>
            <h1 className="text-2xl font-bold">Account Settings</h1>
            <p className="text-sm text-muted-foreground">Manage your account</p>
          </div>
        </div>

        {/* Messages */}
        {message && (
          <div
            className={`mb-6 p-4 rounded-lg flex gap-3 ${
              message.type === "success"
                ? "bg-green-50 dark:bg-green-950 text-green-900 dark:text-green-100"
                : "bg-red-50 dark:bg-red-950 text-red-900 dark:text-red-100"
            }`}
          >
            <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
            <p className="text-sm">{message.text}</p>
          </div>
        )}

        {/* Profile Card */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Profile Information</CardTitle>
            <CardDescription>Update your account details</CardDescription>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleUpdateProfile} className="space-y-6">
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={profile?.email || ""}
                  disabled
                  className="bg-muted"
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="full-name">Full Name</Label>
                <Input
                  id="full-name"
                  type="text"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="Your full name"
                />
              </div>

              <Button type="submit" disabled={isSaving}>
                {isSaving ? "Saving..." : "Save Changes"}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Security Card */}
        <Card>
          <CardHeader>
            <CardTitle>Security</CardTitle>
            <CardDescription>
              Manage your password and security settings
            </CardDescription>
          </CardHeader>

          <CardContent>
            <Button variant="outline" onClick={handleChangePassword}>
              Change Password
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
