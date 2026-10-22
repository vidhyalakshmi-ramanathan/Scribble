"use client";

import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import type { LoginResult, User } from "./auth";

type AuthContextType = {
    user: User | null;
    token: string | null;
    isLoading: boolean;
    loginSession: (result: LoginResult) => void;
    logout: () => void;
};

const AuthContext = createContext<AuthContextType | null>(null);

const SESSION_KEY = "scribble-session";

type StoredSession = {
    user: User;
    token: string;
};

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        try {
            const storedSession = sessionStorage.getItem(SESSION_KEY);

            if (storedSession) {
                const session = JSON.parse(storedSession) as StoredSession;

                if (session.user && session.token) {
                    setUser(session.user);
                    setToken(session.token);
                } else {
                    sessionStorage.removeItem(SESSION_KEY);
                }
            }
        } catch {
            sessionStorage.removeItem(SESSION_KEY);
        } finally {
            setIsLoading(false);
        }
    }, []);

    function loginSession(result: LoginResult) {
        sessionStorage.setItem(
            SESSION_KEY,
            JSON.stringify({
                user: result.user,
                token: result.token,
            }),
        );
        setUser(result.user);
        setToken(result.token);
    }
    function logout() {
        sessionStorage.removeItem(SESSION_KEY);
        setUser(null);
        setToken(null);
    }

    return (
        <AuthContext.Provider
            value={{ user, token, isLoading, loginSession, logout }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error("useAuth must be used inside AuthProvider");
    }

    return context;
}