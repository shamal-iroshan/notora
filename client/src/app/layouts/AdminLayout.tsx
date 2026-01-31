import { Outlet } from "react-router-dom";

export function AdminLayout() {
  return (
    <div className="admin-layout">
      <section>
        <main>
          <Outlet />
        </main>
      </section>
    </div>
  );
}
