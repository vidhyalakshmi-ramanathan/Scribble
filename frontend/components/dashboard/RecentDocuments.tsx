"use client";

import { FileText, Clock3, MoreHorizontal, Eye } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

type DocumentItem = {
    title: string;
    category: string;
    updated: string;
    initials: string;
    color: string;
};

type RecentDocumentsProps = {
    documents: DocumentItem[];
    searchQuery: string;
};

const colorStyles: Record<string, string> = {
    blue: "bg-blue-100 text-blue-700",
    purple: "bg-purple-100 text-purple-700",
    green: "bg-emerald-100 text-emerald-700",
    orange: "bg-orange-100 text-orange-700",
};

export default function RecentDocuments({
    documents,
    searchQuery,
}: RecentDocumentsProps) {
    const filteredDocuments = documents.filter((doc) =>
        doc.title.toLowerCase().includes(searchQuery.trim().toLowerCase()),
    );

    return (
        <Card
            id="recent"
            className="mt-3 gap-0 rounded-2xl border-slate-200 py-0 shadow-sm"
        >
        <CardHeader className="flex flex-row items-center justify-between gap-3 border-b border-slate-100 px-5 py-5">
            <div>
                <CardTitle className="text-base font-bold text-slate-900">
                    Recent documents
                </CardTitle>
                <p className="mt-1 text-sm text-slate-500">
                    Pick up where you left off
                </p>
            </div>
        </CardHeader>

        <CardContent className="p-0">
            {filteredDocuments.length > 0 ? (
                <div className="divide-y divide-slate-100">
                    {filteredDocuments.map((doc) => (
                        <div
                            key={doc.title}
                            className="flex items-center gap-3 px-5 py-4 transition hover:bg-slate-50 sm:gap-4"
                        >
                            <div
                                className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl ${
                                    colorStyles[doc.color] ??
                                    "bg-slate-100 text-slate-600"
                                }`}
                            >
                                <FileText size={21} />
                            </div>

                            <div className="min-w-0 flex-1">
                                <h3 className="truncate text-sm font-semibold text-slate-800">
                                    {doc.title}
                                </h3>
                                <p className="mt-1 truncate text-xs text-slate-500">
                                    {doc.category}
                                </p>
                            </div>

                            <div className="hidden items-center gap-1.5 text-xs text-slate-500 md:flex">
                                <Clock3 size={14} />
                                <span>{doc.updated}</span>
                            </div>

                            <div className="hidden items-center -space-x-2 sm:flex">
                                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full border-2 border-white bg-amber-100 text-xs font-semibold text-slate-700">
                                    {doc.initials}
                                </div>

                                <div
                                    aria-label="Additional collaborators"
                                    className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full border-2 border-white bg-slate-100 text-xs font-semibold text-slate-600"
                                >
                                    +
                                </div>
                            </div>
                            <DropdownMenu>
                                <DropdownMenuTrigger className="h-9 w-9 shrink-0 text-slate-400 hover:text-slate-700">
                                    <MoreHorizontal size={19} />
                                </DropdownMenuTrigger>

                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem
                                        onSelect={() =>
                                            document
                                            .getElementById("my-documents")
                                            ?.scrollIntoView({ behavior: "smooth" })
                                        }
                                    >
                                        <Eye size={15} />
                                        View all documents
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    ))}
                </div>
            ) : (
            <div className="flex flex-col items-center px-5 py-12 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-slate-500">
                    <FileText size={23} />
                </div>

                <h3 className="mt-4 text-sm font-semibold text-slate-800">
                    No documents match your search.
                </h3>

                <p className="mt-1 text-sm text-slate-500">
                    Try another search term.
                </p>
            </div>
            )}

            <div className="flex items-center justify-between gap-3 border-t border-slate-100 px-5 py-4">
                <p className="text-xs text-slate-500">
                    Showing {filteredDocuments.length} of {documents.length} documents
                </p>

                <Button
                    type="button"
                    variant="link"
                    className="h-auto shrink-0 p-0 text-sm font-semibold text-blue-600"
                    onClick={() => {
                    document.getElementById("my-documents")?.scrollIntoView({
                        behavior: "smooth",
                    });
                    }}
                >
                    View all
                </Button>
            </div>
        </CardContent>
        </Card>
    );
}