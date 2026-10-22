"use client";

import { Menu, Search, Bell, ChevronDown, UserRound, LogOut } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { useRouter } from "next/navigation";
import { useAuth } from "@/app/lib/authContext";

type HeaderProps = {
  onMenuClick: () => void;
  searchQuery: string;
  onSearchChange: (value: string) => void;
};

export default function Header({ onMenuClick, searchQuery, onSearchChange }: HeaderProps) {
  const { user, logout } = useAuth();
  const router = useRouter();

  const initials = user?.username?.trim().split(/\s+/).map((part) => part[0]).join("").slice(0, 2).toUpperCase() || "U";

  function handleLogout() {
    logout();
    router.replace("/login");
  }

  return (
    <header className="sticky top-0 z-30 flex h-20 items-center justify-between border-b border-slate-200 bg-white/90 px-4 backdrop-blur sm:px-6 lg:px-8">
      <div className="flex items-center gap-3">
        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Open navigation"
          onClick={onMenuClick}
          className="lg:hidden"
        >
          <Menu size={21} />
        </Button>

        <div>
          <p className="text-xs text-slate-400">Workspace</p>
          <h2 className="text-sm font-semibold text-slate-800 sm:text-base">
            My workspace
          </h2>
        </div>
      </div>

      <div className="flex items-center gap-2 sm:gap-4">
        <div className="relative hidden w-44 sm:block md:w-64">
          <Search
            size={18}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
          />

          <Input
            type="search"
            value={searchQuery}
            onChange={(event) => onSearchChange(event.target.value)}
            placeholder="Search documents..."
            aria-label="Search documents"
            className="h-10 border-slate-200 bg-slate-50 pl-10 focus-visible:bg-white"
          />
        </div>

        <Button
          type="button"
          variant="ghost"
          size="icon"
          aria-label="Notifications"
          className="relative"
        >
          <Bell size={19} />
          <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-blue-600 ring-2 ring-white" />
        </Button>

        <DropdownMenu>
          <DropdownMenuTrigger
            className="flex h-auto items-center gap-2 rounded-full px-1.5 py-1 outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            aria-label="Open account menu"
          >
            <Avatar className="h-9 w-9 border-2 border-white shadow-sm">
              <AvatarFallback className="bg-blue-100 font-bold text-blue-700">
                {initials}
              </AvatarFallback>
            </Avatar>

            <ChevronDown size={15} className="hidden text-slate-500 sm:block" />
          </DropdownMenuTrigger>

          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuGroup>
                <DropdownMenuLabel className="font-normal">
                  <p className="truncate text-sm font-semibold">
                    {user?.username || "Your account"}
                  </p>
                  <p className="truncate text-xs text-slate-500">
                    {user?.email || ""}
                  </p>
              </DropdownMenuLabel>
            </DropdownMenuGroup>

            <DropdownMenuSeparator />

            <DropdownMenuItem disabled>
              <UserRound size={16} />
              My account
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem
              onClick={handleLogout}
              className="text-red-600 focus:text-red-600"
            >
              <LogOut size={16} />
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}