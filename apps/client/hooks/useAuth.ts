import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { authService } from "@/services/auth/authService";

const AUTH_KEY = ["auth"] as const;

export function useAuth() {
  return useQuery({
    queryKey: AUTH_KEY,
    queryFn: authService.me,
    retry: false, // don't retry on 401
  });
}

export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authService.login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: AUTH_KEY });
    },
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: authService.register,
  });
}

export function useVerifyEmail() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: authService.verify,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: AUTH_KEY });
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await authService.logout();
    },
    onSettled: () => {
      queryClient.setQueryData(AUTH_KEY, null);
    },
  });
}
