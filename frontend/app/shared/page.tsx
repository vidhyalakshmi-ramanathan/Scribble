"use client";

import { useAppShell } from "@/components/layouts/AppShell";
import Link from "next/link";
import { Search, FileText, Clock3, Users, FolderOpen, MoreHorizontal, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const sharedDocuments = [
  {
    id: "1",
    title: "Project Architecture",
    category: "Engineering",
    owner: "Alex Rivera",
    initials: "AR",
    role: "Editor",
    updated: "10 minutes ago",
    color: "bg-blue-100 text-blue-700",
  },
  {
    id: "2",
    title: "Sprint Planning",
    category: "Team Workspace",
    owner: "Sarah Kim",
    initials: "SK",
    role: "Viewer",
    updated: "2 hours ago",
    color: "bg-purple-100 text-purple-700",
  },
  {
    id: "3",
    title: "Product Requirements",
    category: "Product",
    owner: "Michael Jones",
    initials: "MJ",
    role: "Editor",
    updated: "Yesterday",
    color: "bg-green-100 text-green-700",
  },
];

export default function SharedPage() {
    const { searchQuery, setSearchQuery } = useAppShell();

    const filteredDocuments = sharedDocuments.filter((doc) =>
        `${doc.title} ${doc.owner} ${doc.category}`
        .toLowerCase().includes(searchQuery.trim().toLowerCase())
    );
    
    return (
    <main className="min-h-screen p-5 sm:p-8">
        <div className="mx-auto max-w-6xl">
            <div className="mb-8">
                <h1 className="text-2xl font-bold text-slate-800">
                    Shared With Me
                </h1>
                <p className="mt-2 text-sm text-slate-500">
                    Documents your teammates have shared with you.
                </p>
            </div>

            <div className="relative mb-6">
                <Search
                    size={18}
                    className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400"
                />
                <Input
                    type="search"
                    value={searchQuery}
                    onChange={(event) => setSearchQuery(event.target.value)}
                    placeholder="Search documents or owners..."
                    aria-label="Search shared documents"
                    className="h-12 border-slate-200 bg-white pl-11"
                />
            </div>
                
            <div className="mb-4 flex items-center justify-between">
                <h2 className="font-semibold text-slate-800">
                    Shared documents
                </h2>
                <Badge variant="secondary">
                    {filteredDocuments.length}{" "}
                    {filteredDocuments.length === 1 ? "document" : "documents"}
                </Badge>
            </div>
            
            <Card className="gap-0 overflow-hidden rounded-xl border-slate-200 py-0 shadow-sm">
                <CardContent className="p-0">
                    {filteredDocuments.length > 0 ? (
                    filteredDocuments.map((doc, index) => (
                        <div
                            key={doc.id}
                            className={`flex items-center gap-3 p-4 transition hover:bg-slate-50 sm:gap-4 sm:p-5 ${
                                index < filteredDocuments.length - 1
                                ? "border-b border-slate-100"
                                : ""
                            }`}
                        >
                            <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                                <FileText size={21} />
                            </div>

                            <div className="min-w-0 flex-1">
                                <Link
                                    href={`/documents/${doc.id}`}
                                    className="block truncate text-sm font-medium text-slate-800 hover:text-blue-600"
                                >
                                    {doc.title}
                                </Link>
                                <p className="mt-1 truncate text-xs text-slate-500">
                                    {doc.category} · Shared by {doc.owner}
                                </p>

                                <div className="mt-2 flex items-center gap-1.5 text-xs text-slate-400">
                                    <Clock3 size={13} />
                                    <span>Updated {doc.updated}</span>
                                </div>

                                {/* Show role on small screens too */}
                                <div className="mt-2 sm:hidden">
                                    <Badge
                                        variant="secondary"
                                        className={
                                        doc.role === "Editor"
                                            ? "bg-blue-50 text-blue-700"
                                            : "bg-slate-100 text-slate-600"
                                        }
                                    >
                                        {doc.role}
                                    </Badge>
                                </div>
                            </div>

                            <div className="hidden shrink-0 text-right sm:block">
                                <Badge
                                variant="secondary"
                                className={
                                    doc.role === "Editor"
                                    ? "bg-blue-50 text-blue-700"
                                    : "bg-slate-100 text-slate-600"
                                }
                                >
                                    {doc.role}
                                </Badge>

                                <div className="mt-2 flex items-center justify-end gap-1 text-xs text-slate-400">
                                    <Users size={13} />
                                    <span>{doc.initials}</span>
                                </div>
                            </div>

                            <DropdownMenu>
                                <DropdownMenuTrigger className="h-9 w-9 shrink-0 text-slate-400 hover:text-slate-700">
                                    <MoreHorizontal size={20} />
                                </DropdownMenuTrigger>

                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem>
                                        <Link href={`/documents/${doc.id}`}>
                                            <ExternalLink size={16} />
                                            Open document
                                        </Link>
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    ))) : (
                    <div className="flex flex-col items-center px-5 py-12 text-center">
                        <FolderOpen
                            size={30}
                            className="text-slate-400"
                        />
                        <h3 className="mt-3 text-sm font-semibold text-slate-700">
                            No shared documents found
                        </h3>
                        <p className="mt-1 text-sm text-slate-500">
                            Try a different search term.
                        </p>

                        {searchQuery.trim() && (
                        <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            className="mt-4"
                            onClick={() => setSearchQuery("")}
                        >
                            Clear search
                        </Button>
                        )}
                    </div>
                    )}
                </CardContent>
            </Card>
        </div>
    </main>
    );
}


        

       

        