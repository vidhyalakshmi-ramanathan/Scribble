"use client";

import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import Sidebar from "./Sidebar";
import Header from "./Header";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/app/lib/authContext";

type AppShellContextType = {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
};

const AppShellContext = createContext<AppShellContextType | null>(null);

export function useAppShell() {
    const context = useContext(AppShellContext);

    if (!context) {
        throw new Error("useAppShell must be used inside AppShell");
    }

    return context;
}

export default function AppShell({ children }: { children: ReactNode }) {
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");

    const pathname = usePathname();
    const router = useRouter();
    const { user, isLoading } = useAuth();

    const isPublicRoute = pathname === "/login" || pathname === "/signup";

    useEffect(() => {
        if (!isLoading && !user && !isPublicRoute) {
            router.replace("/login");
        }
    }, [isLoading, user, isPublicRoute, router]);

    // Login and signup remain accessible without authentication.
    if (isPublicRoute) {
        return <>{children}</>;
    }

    // Wait for session restoration before checking authentication.
    if (isLoading || !user) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-slate-50">
                <p className="text-sm text-slate-500">
                    {isLoading ? "Restoring your session..." : "Redirecting to login..."}
                </p>
            </div>
        );
    }

    return (
    <AppShellContext.Provider
        value={{
            searchQuery,
            setSearchQuery,
        }}
    >
        <Sidebar
            open={sidebarOpen}
            onClose={() => setSidebarOpen(false)}
        />

        <div className="min-h-screen lg:pl-64">
            <Header
                onMenuClick={() => setSidebarOpen(true)}
                searchQuery={searchQuery}
                onSearchChange={setSearchQuery}
            />

            {children}

            <footer className="sticky mt-auto border-t border-slate-200 px-4 py-6 text-center text-xs text-slate-400">
                © 2026 Scribble. Built for better collaboration · Your ideas, in scribble.
            </footer>
        </div>
    </AppShellContext.Provider>
    );
}