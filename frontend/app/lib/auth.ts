import { apiRequest } from "./api";

export type User = {
  userID: string;
  username: string;
  email: string;
};

export type LoginResult = {
  user: User;
  token: string;
};

type ApiResponse<T> = {
  message: string;
  data?: T;
};

type SignupInput = {
  username: string;
  email: string;
  password: string;
};

type LoginInput = {
  email: string;
  password: string;
};

export async function signup( input: SignupInput ) {
    return apiRequest<ApiResponse<null>>("/auth/signup", {
        method: "POST",
        body: JSON.stringify(input),
    });
}

export async function login( input: LoginInput ): Promise<LoginResult> {
    const response = await apiRequest<ApiResponse<LoginResult>>(
        "/auth/login",
        {
            method: "POST",
            body: JSON.stringify(input),
        },
    );

    if (!response.data?.token || !response.data.user) {
        throw new Error("Login response is missing the user or token.");
    }

    return response.data;
}