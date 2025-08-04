import { createBrowserRouter, RouterProvider } from "react-router";

import { ErrorPage } from "../pages/ErrorPage";
import { HomePage } from "../pages/HomePage";
import { LoginPage } from "../pages/LoginPage";
import { LogoutPage } from "../pages/LogoutPage";

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
    errorElement: <ErrorPage />,
  },
  {
    path: "login",
    element: <LoginPage />,
  },
  {
    path: "logout",
    element: <LogoutPage />,
  },
]);

const Router = () => {
  return <RouterProvider router={router} />;
};

export { Router };
