import { isRouteErrorResponse, useRouteError } from "react-router";

const ErrorPage = () => {
  const error = useRouteError();

  const getErrorMessage = () => {
    console.error(error);

    if (isRouteErrorResponse(error)) {
      return error.statusText || error.data?.message || "Unknown error";
    }
    if (error instanceof Error) {
      return error.message;
    }
    if (typeof error === "string") {
      return error;
    }

    return "Unknown error";
  };

  return (
    <div id="error-page">
      <h1>Oops!</h1>
      <p>Sorry, an unexpected error has occurred.</p>
      <p>
        <i>{getErrorMessage()}</i>
      </p>
    </div>
  );
};

export { ErrorPage };
