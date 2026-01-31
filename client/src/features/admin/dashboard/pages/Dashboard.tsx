import type React from "react";
import { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";

import { adminAPI } from "@/lib/api";
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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

import { LogOut, Shield, Plus, Check, X, Trash2, Key } from "lucide-react";
import {
  getPathWithRoot,
  PATH_ADMIN_LOGIN,
  ROOT_ADMIN,
} from "@/app/router/routes";

interface UserProfile {
  id: string;
  email: string;
  full_name: string | null;
  status: "pending" | "approved" | "rejected";
  created_at: string;
}

export function AdminDashboard() {
  const [users, setUsers] = useState<UserProfile[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState<
    "all" | "pending" | "approved" | "rejected"
  >("all");
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [createForm, setCreateForm] = useState({
    email: "",
    fullName: "",
    password: "",
  });
  const [newPassword, setNewPassword] = useState("");
  const [isProcessing, setIsProcessing] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const admin = await adminAPI.getCurrentAdmin();
        if (!admin) {
          navigate(getPathWithRoot(ROOT_ADMIN, PATH_ADMIN_LOGIN), {
            replace: true,
          });
          return;
        }

        const allUsers = await adminAPI.getAllUsers();
        setUsers(allUsers);
      } catch (error) {
        console.error("Error loading admin data:", error);
        navigate(getPathWithRoot(ROOT_ADMIN, PATH_ADMIN_LOGIN), {
          replace: true,
        });
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, [navigate]);

  const filteredUsers = users.filter((u) => {
    if (filter === "all") return true;
    return u.status === filter;
  });

  const handleApprove = useCallback(async (userId: string) => {
    setIsProcessing(true);
    try {
      await adminAPI.approveUser(userId);
      setUsers((prev) =>
        prev.map((u) => (u.id === userId ? { ...u, status: "approved" } : u)),
      );
    } catch (error) {
      console.error("Error approving user:", error);
    } finally {
      setIsProcessing(false);
    }
  }, []);

  const handleReject = useCallback(async (userId: string) => {
    setIsProcessing(true);
    try {
      await adminAPI.rejectUser(userId);
      setUsers((prev) =>
        prev.map((u) => (u.id === userId ? { ...u, status: "rejected" } : u)),
      );
    } catch (error) {
      console.error("Error rejecting user:", error);
    } finally {
      setIsProcessing(false);
    }
  }, []);

  const handleCreateUser = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      setIsProcessing(true);
      try {
        const { user, error } = await adminAPI.createUserProfile(
          createForm.email,
          createForm.fullName,
          createForm.password,
        );

        if (error) throw new Error(error);

        if (user) {
          setUsers((prev) => [user, ...prev]);
          setCreateForm({ email: "", fullName: "", password: "" });
          setShowCreateDialog(false);
        }
      } catch (error) {
        console.error("Error creating user:", error);
      } finally {
        setIsProcessing(false);
      }
    },
    [createForm],
  );

  const handleChangePassword = useCallback(async () => {
    if (!selectedUserId || !newPassword) return;

    setIsProcessing(true);
    try {
      await adminAPI.changeUserPassword(selectedUserId, newPassword);
      setShowPasswordDialog(false);
      setNewPassword("");
      setSelectedUserId(null);
    } catch (error) {
      console.error("Error changing password:", error);
    } finally {
      setIsProcessing(false);
    }
  }, [selectedUserId, newPassword]);

  const handleDeleteUser = useCallback(async (userId: string) => {
    setIsProcessing(true);
    try {
      await adminAPI.deleteUser(userId);
      setUsers((prev) => prev.filter((u) => u.id !== userId));
    } catch (error) {
      console.error("Error deleting user:", error);
    } finally {
      setIsProcessing(false);
    }
  }, []);

  const handleLogout = async () => {
    await adminAPI.adminLogout();
    navigate(getPathWithRoot(ROOT_ADMIN, PATH_ADMIN_LOGIN), { replace: true });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <Shield className="w-12 h-12 text-primary mx-auto mb-4 animate-pulse" />
          <p className="text-muted-foreground">Loading admin dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="max-w-7xl mx-auto px-4 md:px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Shield className="w-6 h-6 text-primary" />
            <h1 className="text-2xl font-bold">Admin Dashboard</h1>
          </div>

          <Button
            variant="outline"
            onClick={handleLogout}
            className="gap-2 bg-transparent"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 md:px-6 py-8">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Total Users
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{users.length}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Pending
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-yellow-600">
                {users.filter((u) => u.status === "pending").length}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Approved
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-600">
                {users.filter((u) => u.status === "approved").length}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                Rejected
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-600">
                {users.filter((u) => u.status === "rejected").length}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* User Management Section */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>User Management</CardTitle>
              <CardDescription>
                Manage user accounts, approvals, and permissions
              </CardDescription>
            </div>
            <Button onClick={() => setShowCreateDialog(true)} className="gap-2">
              <Plus className="w-4 h-4" />
              Create User
            </Button>
          </CardHeader>
          <CardContent>
            {/* Filter Tabs */}
            <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
              {(["all", "pending", "approved", "rejected"] as const).map(
                (f) => (
                  <Button
                    key={f}
                    variant={filter === f ? "default" : "outline"}
                    size="sm"
                    onClick={() => setFilter(f)}
                    className="capitalize"
                  >
                    {f}
                  </Button>
                ),
              )}
            </div>

            {/* Users Table */}
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left py-3 px-4 font-medium">Email</th>
                    <th className="text-left py-3 px-4 font-medium">Name</th>
                    <th className="text-left py-3 px-4 font-medium">Status</th>
                    <th className="text-left py-3 px-4 font-medium">Created</th>
                    <th className="text-right py-3 px-4 font-medium">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {filteredUsers.map((user) => (
                    <tr
                      key={user.id}
                      className="border-b border-border hover:bg-muted/50 transition-colors"
                    >
                      <td className="py-3 px-4 font-mono text-xs">
                        {user.email}
                      </td>
                      <td className="py-3 px-4">{user.full_name || "—"}</td>
                      <td className="py-3 px-4">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                            user.status === "pending"
                              ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
                              : user.status === "approved"
                                ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                                : "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                          }`}
                        >
                          {user.status}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-xs text-muted-foreground">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <div className="flex gap-2 justify-end flex-wrap">
                          {user.status === "pending" && (
                            <>
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleApprove(user.id)}
                                disabled={isProcessing}
                                className="gap-1"
                              >
                                <Check className="w-3 h-3" />
                                Approve
                              </Button>
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleReject(user.id)}
                                disabled={isProcessing}
                                className="gap-1"
                              >
                                <X className="w-3 h-3" />
                                Reject
                              </Button>
                            </>
                          )}
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => {
                              setSelectedUserId(user.id);
                              setShowPasswordDialog(true);
                            }}
                            disabled={isProcessing}
                            className="gap-1"
                          >
                            <Key className="w-3 h-3" />
                            Reset Password
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleDeleteUser(user.id)}
                            disabled={isProcessing}
                            className="gap-1 text-destructive hover:text-destructive"
                          >
                            <Trash2 className="w-3 h-3" />
                            Delete
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {filteredUsers.length === 0 && (
                <div className="text-center py-8 text-muted-foreground">
                  <p>No users found with current filter</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </main>

      {/* Create User Dialog */}
      <AlertDialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
        <AlertDialogContent>
          <AlertDialogTitle>Create New User</AlertDialogTitle>
          <AlertDialogDescription>
            Add a new user account directly without approval
          </AlertDialogDescription>
          <form onSubmit={handleCreateUser} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="create-email">Email</Label>
              <Input
                id="create-email"
                type="email"
                placeholder="user@example.com"
                value={createForm.email}
                onChange={(e) =>
                  setCreateForm({ ...createForm, email: e.target.value })
                }
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="create-name">Full Name</Label>
              <Input
                id="create-name"
                placeholder="John Doe"
                value={createForm.fullName}
                onChange={(e) =>
                  setCreateForm({ ...createForm, fullName: e.target.value })
                }
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="create-password">Password</Label>
              <Input
                id="create-password"
                type="password"
                placeholder="••••••••"
                value={createForm.password}
                onChange={(e) =>
                  setCreateForm({ ...createForm, password: e.target.value })
                }
                required
              />
            </div>
            <div className="flex gap-2 justify-end">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction type="submit" disabled={isProcessing}>
                {isProcessing ? "Creating..." : "Create"}
              </AlertDialogAction>
            </div>
          </form>
        </AlertDialogContent>
      </AlertDialog>

      {/* Change Password Dialog */}
      <AlertDialog
        open={showPasswordDialog}
        onOpenChange={setShowPasswordDialog}
      >
        <AlertDialogContent>
          <AlertDialogTitle>Reset User Password</AlertDialogTitle>
          <AlertDialogDescription>
            Enter the new password for the selected user
          </AlertDialogDescription>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="new-password">New Password</Label>
              <Input
                id="new-password"
                type="password"
                placeholder="••••••••"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
              />
            </div>
            <div className="flex gap-2 justify-end">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={handleChangePassword}
                disabled={isProcessing || !newPassword}
              >
                {isProcessing ? "Updating..." : "Update Password"}
              </AlertDialogAction>
            </div>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
