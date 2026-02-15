import { AxiosError } from "axios";

export function handleError(error: unknown): never {
  if (error instanceof AxiosError && error.response?.data) {
    const apiError = error.response.data as APIError;
    throw apiError;
  }
  throw {
    code: "UNKNOWN_ERROR",
    message: "Something went wrong",
    status: 500,
    override: false,
  } satisfies APIError;
}
