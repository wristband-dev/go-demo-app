import axios from "axios";

import { isForbiddenError, isUnauthorizedError } from "../utils/wristband-utils";
import { redirectToLogin } from "@wristband/react-client-auth";

/* CSRF_TOUCHPOINT */
const backendApiClient = axios.create({
  headers: { "Content-Type": "application/json", Accept: "application/json" },
  xsrfCookieName: "XSRF-TOKEN",
  xsrfHeaderName: "X-XSRF-TOKEN",
  withXSRFToken: true,
  withCredentials: true,
});

/* WRISTBAND_TOUCHPOINT - AUTHENTICATION */
// Any HTTP 401s should trigger the user to go log in again.  This happens when their
// session cookie has expired and/or the CSRF cookie/header are missing in the request.
// You can optionally catch HTTP 403s as well.
const unauthorizedAccessInterceptor = async (error: { response: { status: number } }) => {
  if (isUnauthorizedError(error) || isForbiddenError(error)) {
    redirectToLogin('/api/auth/login');
  }

  return Promise.reject(error);
};

backendApiClient.interceptors.response.use(undefined, unauthorizedAccessInterceptor);

export { backendApiClient };
