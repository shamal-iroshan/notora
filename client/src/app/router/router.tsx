import { createBrowserRouter, Navigate } from "react-router-dom";
import { ClientLayout } from "../layouts/ClientLayout";
import { AdminLayout } from "../layouts/AdminLayout";
import {
  PATH_ADMIN_DASHBOARD,
  PATH_ADMIN_LOGIN,
  PATH_CLIENT_ACCOUNT,
  PATH_CLIENT_LOGIN,
  PATH_CLIENT_NOTES,
  PATH_CLIENT_SIGN_UP,
  PATH_CLIENT_SIGN_UP_SUCCESS,
  ROOT_ADMIN,
  ROOT_INDEX,
} from "./routes";
import { AdminLoginPage } from "@/features/admin/login/pages/Login";
import { AdminDashboard } from "@/features/admin/dashboard/pages/Dashboard";
import { LoginPage } from "@/features/client/login/pages/Login";
import { SignUpPage } from "@/features/client/signUp/pages/SignUp";
import { SignUpSuccessPage } from "@/features/client/signUp/pages/SignUpSuccess";
import { AccountPage } from "@/features/client/account/pages/Account";
import { NotesPage } from "@/features/client/notes/pages/Notes";

export const router = createBrowserRouter([
  {
    path: ROOT_INDEX,
    element: <ClientLayout />,
    errorElement: <p>not found</p>,
    children: [
      { index: true, element: <Navigate to={PATH_CLIENT_LOGIN} replace /> },
      { path: PATH_CLIENT_LOGIN, element: <LoginPage /> },
      { path: PATH_CLIENT_SIGN_UP, element: <SignUpPage /> },
      { path: PATH_CLIENT_SIGN_UP_SUCCESS, element: <SignUpSuccessPage /> },
      { path: PATH_CLIENT_ACCOUNT, element: <AccountPage /> },
      { path: PATH_CLIENT_NOTES, element: <NotesPage /> },
    ],
  },
  {
    path: ROOT_ADMIN,
    element: <AdminLayout />,
    children: [
      { index: true, element: <Navigate to={PATH_ADMIN_LOGIN} replace /> },
      { path: PATH_ADMIN_LOGIN, element: <AdminLoginPage /> },
      { path: PATH_ADMIN_DASHBOARD, element: <AdminDashboard /> },
    ],
  },
  {
    path: "*",
    element: <p>not found</p>,
  },
]);
