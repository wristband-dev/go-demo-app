import { useEffect } from "react";

import { redirectToLogin } from "@wristband/react-client-auth";

const LoginPage = () => {
  useEffect(() => redirectToLogin('/api/auth/login'), []);
  return null;
};

export { LoginPage };
