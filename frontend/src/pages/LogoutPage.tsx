import { useEffect } from "react";
import { redirectToLogout } from "@wristband/react-client-auth";

const LogoutPage = () => {
  useEffect(() => redirectToLogout('/api/auth/logout'), []);
  return null;
};

export { LogoutPage };
