import { handleError } from "@/common/utils";
import { API_URL } from "@/constants/constants";
import axios from "axios";

export const authService = {
  login: async (req: LoginRequest) => {
    try {
      const { data } = await axios.post(`${API_URL}/auth/login`, req, {
        withCredentials: true,
      });
      return data;
    } catch (error) {
      handleError(error);
    }
  },
  register: async (req: RegisterRequest) => {
    try {
      const { data } = await axios.post(`${API_URL}/auth/register`, req, {
        withCredentials: true,
      });
      return data;
    } catch (error) {
      handleError(error);
    }
  },
  verify: async (req: VerifyEmailRequest) => {
    try {
      const { data } = await axios.post(`${API_URL}/auth/verify-email`, req, {
        withCredentials: true,
      });
      return data;
    } catch (error) {
      handleError(error);
    }
  },
  me: async () => {
    try {
      const { data } = await axios.get(`${API_URL}/auth/me`, {
        withCredentials: true,
      });
      return data;
    } catch (error) {
      handleError(error);
    }
  },
  logout: async () => {
    try {
      await axios.post(`${API_URL}/auth/logout`, {
        withCredentials: true,
      });
    } catch (error) {
      handleError(error);
    }
  },
};
