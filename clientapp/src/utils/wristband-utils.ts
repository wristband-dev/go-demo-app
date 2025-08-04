import { AxiosError } from 'axios';

/**
 * Checks if an error represents a specific HTTP status code error.
 *
 * @param {unknown} error - The error to check. Must be either an AxiosError or a Response object.
 * @param {number} statusCode - The HTTP status code to check for.
 * @returns {boolean} True if the error has the specified status code, false otherwise.
 * @throws {TypeError} If the error is null, undefined, or not an AxiosError or Response object.
 *
 * @example
 * // With Axios
 * try {
 *   await axios.get('/api/resource');
 * } catch (error) {
 *   if (isHttpStatusError(error, 404)) {
 *     console.log('Resource not found');
 *   }
 * }
 *
 * @example
 * // With Fetch
 * const response = await fetch('/api/resource');
 * if (isHttpStatusError(response, 401)) {
 *   console.log('Authentication required');
 * }
 */
function isHttpStatusError(error: unknown, statusCode: number): boolean {
  // Handle null/undefined case with an exception
  if (error === null || error === undefined) {
    throw new TypeError('Argument [error] cannot be null or undefined');
  }

  // Handle Axios error format
  if (error instanceof AxiosError) {
    return error.response?.status === statusCode;
  }

  // Handle fetch Response objects
  if (error instanceof Response) {
    return error.status === statusCode;
  }

  // If it's neither of the expected types, throw an error.
  throw new TypeError(
    `Invalid error type: Expected either an AxiosError or a Response object, but received type: [${typeof error}] `
  );
}

/**
 * Checks if an error represents an HTTP 401 Unauthorized error.
 *
 * @param {unknown} error - The error to check. Must be either an AxiosError or a Response object.
 * @returns {boolean} True if the error has a 401 status code, false otherwise.
 * @throws {TypeError} If the error is null, undefined, or not an AxiosError or Response object.
 *
 * @example
 * // With Axios
 * try {
 *   await axios.get('/api/resource');
 * } catch (error) {
 *   if (isUnauthorizedError(error)) {
 *     console.log('Authentication required');
 *   }
 * }
 *
 * @example
 * // With Fetch
 * const response = await fetch('/api/resource');
 * if (isUnauthorizedError(response)) {
 *   console.log('Authentication required');
 * }
 */
export const isUnauthorizedError = (error: unknown) => isHttpStatusError(error, 401);

/**
 * Checks if an error represents an HTTP 403 Forbidden error.
 *
 * @param {unknown} error - The error to check. Must be either an AxiosError or a Response object.
 * @returns {boolean} True if the error has a 403 status code, false otherwise.
 * @throws {TypeError} If the error is null, undefined, or not an AxiosError or Response object.
 *
 * @example
 * // With Axios
 * try {
 *   await axios.get('/api/resource');
 * } catch (error) {
 *   if (isForbiddenError(error)) {
 *     console.log('Forbidden');
 *   }
 * }
 *
 * @example
 * // With Fetch
 * const response = await fetch('/api/resource');
 * if (isForbiddenError(response)) {
 *   console.log('Forbidden');
 * }
 */
export const isForbiddenError = (error: unknown) => isHttpStatusError(error, 403);

export function isOwnerRole(roleName: string) {
  // Should match the Role "name" field, i.e. "app:invotasticb2b:owner"
  return /^app:.*:owner$/.test(roleName);
}
