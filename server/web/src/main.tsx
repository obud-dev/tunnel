import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { createBrowserRouter, Outlet, RouterProvider } from "react-router-dom";
import { Toaster } from "~/components/ui/toaster";
import Tunnels from "./routes/tunnels.tsx";
import Routes from "./routes/tunnel.tsx";
import { Button } from "~/components/ui/button.tsx";
import { AdjustmentsHorizontalIcon } from "@heroicons/react/24/solid";
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "./components/ui/sheet.tsx";

const Layout = () => (
  <>
    <header className="h-14 bg-primary text-white flex justify-center px-8 italic font-bold">
      <div className="max-w-7xl w-full justify-between h-full flex items-center">
        <span className="text-lg">Tunnel</span>

        <Sheet>
          <SheetTrigger asChild>
            <Button size="icon">
              <AdjustmentsHorizontalIcon className="w-6 h-6" />
            </Button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>Server Info</SheetTitle>
              <SheetFooter>
                <Button>Save Changes</Button>
              </SheetFooter>
            </SheetHeader>
          </SheetContent>
        </Sheet>
      </div>
    </header>
    <div className="flex justify-center">
      <main className="p-4 w-full max-w-7xl ">
        <Outlet />
      </main>
    </div>
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
