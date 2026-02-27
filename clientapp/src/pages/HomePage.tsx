import { useCallback, useState } from "react";
import { isAxiosError } from "axios";
import { redirectToLogout, useWristbandSession } from '@wristband/react-client-auth';

import goLogo from "../assets/gologo.png";
import reactLogo from "../assets/react.svg";
import wristbandLogo from "../assets/wristband.png";
import { backendApiClient } from "../api/backend-api-client";
import { MySessionData } from "../types";

const HomePage = () => {
  // React State
  const [protectedResult, setProtectedResult] = useState<number>(0);
  const [unprotectedResult, setUnprotectedResult] = useState<number>(0);

  /* WRISTBAND_TOUCHPOINT - AUTHENTICATION */
  const { metadata, userId, tenantId, updateMetadata } = useWristbandSession<MySessionData>();

  const { hasOwnerRole } = metadata;
  const mySessionData = JSON.stringify({ userId, tenantId, metadata });

  const generateRandomNameForSession = () => {
    const newName = Math.random().toString(36).substring(2);
    updateMetadata({ fullName: newName });
    alert(`New fullName for React Context:\n\n"${newName}"`)
  };

  const callProtectedEndpoint = useCallback(async () => {
    try {
      const response = await backendApiClient.get("/api/protected");
      const result = response.data as { message: string, value: number };
      setProtectedResult((prior) => prior + result.value);
    } catch (error) {
      if (isAxiosError(error)) {
        console.error('Axios Error:', error.response?.status, error.response?.data);
      } else {
        console.error('Unexpected Error:', error);
      }
    }
  }, [setProtectedResult]);

  const callUnprotectedEndpoint = useCallback(async () => {
    try {
      const response = await backendApiClient.get("/api/unprotected");
      const result = response.data as { message: string, value: number };
      setUnprotectedResult((prior) => prior + result.value);
    } catch (error) {
      if (isAxiosError(error)) {
        console.error('Axios Error:', error.response?.status, error.response?.data);
      } else {
        console.error('Unexpected Error:', error);
      }
    }
  }, [setUnprotectedResult]);

  return (
    <>
      <div>
        <a href="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" target="_blank">
          <img src={goLogo} className="logo" alt="Go logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo" alt="React logo" />
        </a>
        <a href="https://wristband.dev" target="_blank">
          <img src={wristbandLogo} className="logo react" alt="Wristband logo" />
        </a>
      </div>
      <h1>Go + React + Wristband</h1>
      <div className="card">
        <button onClick={() => alert(`${mySessionData}`)}>
          View Session Context Data
        </button>
      </div>
      {hasOwnerRole && (
        <div className="card">
          <button onClick={generateRandomNameForSession}>
            Update Session Context Data
          </button>
        </div>
      )}
      <div className="card">
        <button onClick={callProtectedEndpoint}>
          Protected API count {protectedResult}
        </button>
      </div>
      <div className="card">
        <button onClick={callUnprotectedEndpoint}>
          Unprotected API count {unprotectedResult}
        </button>
      </div>
      <div className="card">
        <button onClick={() => redirectToLogout('/api/auth/logout')}>
          Logout
        </button>
      </div>
    </>
  );
};

export { HomePage };
