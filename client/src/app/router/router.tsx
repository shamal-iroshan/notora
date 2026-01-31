import { createBrowserRouter, Navigate } from "react-router-dom";
import { ClientLayout } from "../layouts/ClientLayout";
import { AdminLayout } from "../layouts/AdminLayout";
import { PATH_ADMIN_DASHBOARD, PATH_ADMIN_LOGIN, ROOT_ADMIN } from "./routes";
import { AdminLoginPage } from "@/features/admin/login/pages/Login";
import { AdminDashboard } from "@/features/admin/dashboard/pages/Dashboard";

export const router = createBrowserRouter([
  {
    element: <ClientLayout />,
    errorElement: <p>not found</p>,
    children: [
      { index: true, element: <p>home</p> },
      { path: "profile", element: <p>profile</p> },
      { path: "products", element: <p>products</p> },
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
