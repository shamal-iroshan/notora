import { Outlet } from "react-router-dom";

export function ClientLayout() {
  return (
    <>
      {/* Client Navbar */}
      <header>Client Header</header>

      <main>
        <Outlet />
      </main>

      {/* Client Footer */}
    </>
  );
}
