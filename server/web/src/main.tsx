import React from "react";
import ReactDOM from "react-dom/client";
import Home from "./routes/tunnels.tsx";
import "./index.css";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import Routes from "./routes/routes.tsx";
import { Toaster } from "@/components/ui/toaster";

const Layout = () => (
  <div>
    <Outlet />
    <Toaster />
  </div>
);

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout/>,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "routes",
        element: <Routes />,
      },
    ]
  }
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
