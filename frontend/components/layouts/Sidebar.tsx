"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { FileText, X, LayoutDashboard, FolderOpen, Users, Clock3, Settings, LogOut, ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/app/lib/authContext";

type SidebarProps = {
    open: boolean;
    onClose: () => void;
};

const navigationItems = [
    { label: "Dashboard", href: "/", icon: LayoutDashboard },
    { label: "My documents", href: "/documents", icon: FolderOpen },
    { label: "Shared with me", href: "/shared", icon: Users },
    { label: "Recent activity", href: "/activity", icon: Clock3 },
];

export default function Sidebar({ open, onClose }: SidebarProps) {
    const pathname = usePathname();
    const router = useRouter();
    const { user, logout } = useAuth();

    function isActive(href: string) {
        if (href === "/") return pathname === "/";
            return pathname === href || pathname.startsWith(`${href}/`);
    }

    function handleLogout() {
        logout();
        onClose();
        router.replace("/login");
    }

    const initials = user?.username?.trim().split(/\s+/).map((part) => part[0]).join("").slice(0, 2).toUpperCase() || "U";

    return (
        <>
        {open && (
            <button
                type="button"
                aria-label="Close navigation"
                onClick={onClose}
                className="fixed inset-0 z-40 bg-slate-950/40 lg:hidden"
            />
        )}
        <aside
            className={`fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-slate-200 bg-white transition-transform duration-200 lg:translate-x-0 ${
            open ? "translate-x-0" : "-translate-x-full"
            }`}
        >
            <div className="flex h-20 items-center justify-between px-6">
                <Link
                    href="/"
                    onClick={onClose}
                    className="flex items-center gap-3"
                >
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-600 text-white shadow-sm shadow-blue-200">
                        <FileText size={22} />
                    </div>

                    <span className="text-xl font-bold tracking-tight">
                        Sync<span className="text-blue-600">Write</span>
                    </span>
                </Link>

                <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    aria-label="Close navigation"
                    onClick={onClose}
                    className="lg:hidden"
                >
                    <X size={20} />
                </Button>
            </div>

            <div className="px-4 pt-5">
                <p className="mb-3 px-3 text-xs font-semibold uppercase tracking-widest text-slate-400">
                    Workspace
                </p>

                <nav aria-label="Main navigation" className="space-y-1">
                    {navigationItems.map(({ label, href, icon: Icon }) => {
                        const active = isActive(href);

                        return (
                            <Link
                                key={href}
                                href={href}
                                onClick={onClose}
                                aria-current={active ? "page" : undefined}
                                className={`flex items-center gap-3 rounded-xl px-3 py-3 text-sm transition ${
                                    active
                                    ? "bg-blue-50 font-semibold text-blue-700"
                                    : "font-medium text-slate-600 hover:bg-slate-50"
                                }`}
                            >
                                <Icon size={19} />
                                <span>{label}</span>
                            </Link>
                        );
                    })}
                </nav>
            </div>
            <div className="mt-auto border-t border-slate-100 p-4">
                <Button
                    type="button"
                    variant="ghost"
                    disabled
                    className="w-full justify-start gap-3 rounded-xl px-3 py-3 text-sm text-slate-600"
                >
                    <Settings size={19} />
                    Settings
                </Button>

                <Button
                    type="button"
                    variant="ghost"
                    onClick={handleLogout}
                    className="w-full justify-start gap-3 rounded-xl px-3 py-3 text-sm text-slate-600 hover:text-red-600"
                >
                    <LogOut size={19} />
                    Log out
                </Button>

                <div className="mt-3 flex items-center gap-3 rounded-xl bg-slate-50 p-3">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-blue-100 text-sm font-bold text-blue-700">
                        {initials}
                    </div>

                    <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-semibold">
                            {user?.username || "Your account"}
                        </p>
                        <p className="truncate text-xs text-slate-500">
                            {user?.email || "Free workspace"}
                        </p>
                    </div>
                </div>
            </div>
        </aside>
        </>
    );
}