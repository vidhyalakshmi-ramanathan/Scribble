"use client";

import { Plus, Search, FileText, Clock3, Users, MoreHorizontal, ExternalLink } from "lucide-react";
import { useAppShell } from "@/components/layouts/AppShell";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
    Card,
    CardContent,
} from "@/components/ui/card";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "../lib/authContext";
import { useEffect, useState } from "react";
import { createDocument, Document, getMyDocuments } from "../lib/doc";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const initials =  "AR"
const color =  "bg-blue-100 text-blue-700"

export default function DocumentsPage() {
    const router = useRouter();
    const [isCreating, setIsCreating] = useState(false);
    const { searchQuery, setSearchQuery } = useAppShell();
    const { token } = useAuth();
    
    const [documents, setDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    
    useEffect(() => {
        if (!token) return;

        let cancelled = false;

        async function fetchDocuments() {
            try {
                setIsLoading(true);
                setError(null);

                const result = await getMyDocuments(token!);
                if (!cancelled) {
                    setDocuments(result);
                }
            } catch (err) {
                if (!cancelled) {
                   setError(
                        err instanceof Error ? err.message : "Failed to load documents.",
                    );
                }
            } finally {
                if (!cancelled) {
                    setIsLoading(false);
                }
            }
        }
        fetchDocuments();
        return () => {
            cancelled = true;
        };
    }, [token]);

    async function handleCreateDocument() {
        if (!token || isCreating) return;

        try {
            setIsCreating(true);

            const document = await createDocument("Untitled document", "", token);

            if (!document.docID) {
                throw new Error("The server did not return a document ID.");
            }
            toast.success("Document created successfully!");
            router.push(`/documents/${document.docID}`);
        } catch (err) {
            toast.error(err instanceof Error ? err.message : "Failed to create document.");
        } finally {
            setIsCreating(false);
        }
    }

    const filteredDocuments = documents.filter((doc) =>
        doc.title.toLowerCase().includes(searchQuery.trim().toLowerCase()),
    );

    return (
    <main className="min-h-screen p-5 sm:p-8">
        <div className="mx-auto max-w-6xl">
            <div className="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
                <div>
                    <h1 className="text-2xl font-bold text-slate-800">
                        My Documents
                    </h1>

                    <p className="mt-2 text-sm text-slate-500">
                        Create, organize, and manage your documents.
                    </p>
                </div>

                <Button
                    type="button"
                    onClick={handleCreateDocument}
                    disabled={isCreating}
                    className="gap-2 rounded-lg bg-blue-600 hover:bg-blue-700"
                >
                    <Plus size={18} />
                    {isCreating ? "Creating..." : "New Document"}
                </Button>
            </div>

            <div className="relative mb-6">
                <Search
                    size={18}
                    className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400"
                  />

                <Input
                    type="search"
                    value={searchQuery}
                    placeholder="Search documents..."
                    aria-label="Search documents"
                    onChange={(event) => setSearchQuery(event.target.value)}
                    className="h-12 border-slate-200 bg-white pl-11"
                />
            </div>

            <div className="mb-4 flex items-center justify-between">
                <h2 className="font-semibold text-slate-800">
                    All documents
                </h2>

                <Badge variant="secondary">
                    {filteredDocuments.length}{" "}
                    {filteredDocuments.length === 1 ? "document" : "documents"}
                </Badge>
            </div>
            
            {isLoading ? (
                <p className="py-8 text-center text-sm text-muted-foreground">
                    Loading documents...
                </p>
            ) : error ? (
                <p className="py-8 text-center text-sm text-destructive">
                    {error}
                </p>
            ) : (
                <Card className="gap-0 overflow-hidden rounded-xl border-slate-200 py-0 shadow-sm">
                    <CardContent className="p-0">
                        {filteredDocuments.length > 0 ? (
                            filteredDocuments.map((doc, index) => (
                                <div
                                    key={doc.docID}
                                    className={`flex items-center gap-3 p-4 transition hover:bg-slate-50 sm:gap-4 sm:p-5 ${
                                        index < documents.length - 1
                                        ? "border-b border-slate-100"
                                        : ""
                                    }`}
                                >
                                    <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                                        <FileText size={21} />
                                    </div>

                                    <div className="min-w-0 flex-1">
                                        <Link
                                            href={`/documents/${doc.docID}`}
                                            className="block truncate text-sm font-medium text-slate-800 transition hover:text-blue-600"
                                        >
                                            {doc.title}
                                        </Link>

                                        <p className="mt-1 truncate text-xs text-slate-500">
                                            Engineering
                                        </p>

                                        <div className="mt-2 flex items-center gap-1.5 text-xs text-slate-400">
                                            <Clock3 size={13} />
                                            <span>Updated {formatRelativeTime(doc.updatedAt)}</span>
                                        </div>
                                    </div>

                                    <div className="hidden items-center gap-2 sm:flex">
                                        <div
                                            className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-semibold ${color}`}
                                            aria-label={`Collaborator ${initials}`}
                                        >
                                            {initials}
                                        </div>

                                        <Users size={16} className="text-slate-400" />
                                    </div>

                                    <DropdownMenu>
                                        <DropdownMenuTrigger className="h-9 w-9 shrink-0 text-slate-400 hover:text-slate-700">
                                            <MoreHorizontal size={20} />
                                        </DropdownMenuTrigger>

                                        <DropdownMenuContent align="end">
                                            <DropdownMenuItem>
                                                <Link href={`/documents/${doc.docID}`}>
                                                    <ExternalLink size={16} />
                                                    Open document
                                                </Link>
                                            </DropdownMenuItem>
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                </div>
                            ))
                        ) : (
                        <div className="flex flex-col items-center px-5 py-12 text-center">
                            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-500">
                                <FileText size={24} />
                            </div>

                            <p className="mt-4 text-sm font-semibold text-slate-800">
                                {searchQuery.trim()
                                    ? "No documents match your search"
                                    : "No documents yet"}
                            </p>

                            <p className="mt-1 text-sm text-slate-500">
                                {searchQuery.trim()
                                    ? "Try another search term."
                                    : "Create your first document to get started."}
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
            )}
        </div>
    </main>
    );
}

function formatRelativeTime(dateString: string): string {
    const date = new Date(dateString);

    if (Number.isNaN(date.getTime())) {
        return "recently";
    }

    const diffSeconds = Math.floor(
        (Date.now() - date.getTime()) / 1000,
    );

    if (diffSeconds < 0) {
        return "just now";
    }

    if (diffSeconds < 60) {
        return "just now";
    }

    const minutes = Math.floor(diffSeconds / 60);
    if (minutes < 60) {
        return `${minutes} minute${minutes === 1 ? "" : "s"} ago`;
    }

    const hours = Math.floor(minutes / 60);
    if (hours < 24) {
        return `${hours} hour${hours === 1 ? "" : "s"} ago`;
    }

    const days = Math.floor(hours / 24);
    if (days === 1) {
        return "yesterday";
    }

    if (days < 7) {
        return `${days} days ago`;
    }

    return date.toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
    });
}