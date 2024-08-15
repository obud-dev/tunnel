import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { Toaster } from "~/components/ui/toaster";
import Tunnels from "./routes/tunnels.tsx";
import Routes from "./routes/tunnel.tsx";

const Layout = () => (
  <>
    <header className="h-14 bg-slate-900 text-white flex items-center px-8 italic font-bold">
      Tunnel
    </header>
    <main className="p-4">
      <Outlet />
    </main>
    <Toaster />
  </>
);

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      {
        path: "/",
        element: <Tunnels />,
      },
      {
        path: "/tunnels/:id",
        element: <Routes />,
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
