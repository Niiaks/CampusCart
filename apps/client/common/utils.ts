import { AxiosError } from "axios";
import type { APIError } from "@/common/types";

export function handleError(error: unknown): never {
  if (error instanceof AxiosError && error.response?.data) {
    const data = error.response.data;
    if (typeof data === "object" && data !== null && "code" in data) {
      throw data as APIError;
    }
  }
  throw {
    code: "UNKNOWN_ERROR",
    message: "Something went wrong",
    status: 500,
    override: false,
  } satisfies APIError;
}
