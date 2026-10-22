"use client";

import { useEffect, useRef, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import dynamic from "next/dynamic";
import { toast } from "sonner";
import { ArrowLeft, MoreHorizontal, Save, Share2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { useAuth } from "@/app/lib/authContext";
import { DocumentVersion, getDocument, getVersions, restoreVersion, saveVersion, updateDocument } from "@/app/lib/doc";
import { History, RotateCcw } from "lucide-react";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import VersionHistoryDialog from "@/components/editor/VersionHistoryDialog";

const RichTextEditor = dynamic(
    () => import("@/components/editor/RichTextEditor"),
    {
        ssr: false,
        loading: () => (
            <div className="min-h-[60vh] p-4 text-sm text-slate-400">
                Loading editor...
            </div>
        ),
    }
);

export default function DocumentEditorPage() {
    const params = useParams<{ id: string }>();
    const docID = params.id;
    const { token } = useAuth();

    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [isDocumentLoaded, setIsDocumentLoaded] = useState(false);
    const [loadError, setLoadError] = useState<string | null>(null);
    const [saveStatus, setSaveStatus] = useState<"saved" | "saving" | "error">("saved");
    const [isSavingVersion, setIsSavingVersion] = useState(false);
    const [isHistoryOpen, setIsHistoryOpen] = useState(false);

    // Remember the last successfully saved document state.
    const lastSavedRef = useRef({
        title: "",
        content: "",
    });

    // Load the document identified by the URL.
    useEffect(() => {
        if (!token || !docID) return;

        let cancelled = false;

        async function loadDocument() {
            setIsLoading(true);
            setIsDocumentLoaded(false);
            setLoadError(null);

            try {
                const doc = await getDocument(docID, token!);

                if (cancelled) return;

                const loadedTitle = doc.title;
                const loadedContent = doc.content ?? "";

                setTitle(loadedTitle);
                setContent(loadedContent);

                lastSavedRef.current = {
                    title: loadedTitle,
                    content: loadedContent,
                };

                setSaveStatus("saved");
                setIsDocumentLoaded(true);
            } catch (err) {
                if (!cancelled) {
                    setLoadError(err instanceof Error ? err.message : "Failed to load document.");
                }
            } finally {
                if (!cancelled) {
                    setIsLoading(false);
                }
            }
        }

        loadDocument();

        return () => {
            cancelled = true;
        };
    }, [docID, token]);

    // Autosave only after 5 minutes without further changes.
    useEffect(() => {
    if (!token || !docID || !isDocumentLoaded) return;

    if ( title === lastSavedRef.current.title && content === lastSavedRef.current.content ) {
        return;
    }
    setSaveStatus("saving");
    const timeout = setTimeout(async () => {
        try {
            await updateDocument(docID, title, content, token);

            lastSavedRef.current = { title, content };
            setSaveStatus("saved");
        } catch (err) {
            setSaveStatus("error");
            toast.error(err instanceof Error ? `Save failed: ${err.message}` : "Failed to save document.");
        }
    }, 1 * 60 * 1000);

    return () => {
        clearTimeout(timeout);
    };
}, [title, content, docID, token, isDocumentLoaded]);

    const wordCount = content.replace(/<[^>]*>/g, " ").trim().split(/\s+/).filter(Boolean).length;

    if (isLoading) {
        return (
            <main className="p-8 text-sm text-slate-500">
                Loading document...
            </main>
        );
    }

    if (loadError) {
        return (
            <main className="space-y-4 p-8">
                <p className="text-sm text-destructive">{loadError}</p>
                <Button
                    variant="outline"
                    onClick={() => window.location.reload()}
                >
                    Try again
                </Button>
            </main>
        );
    }

    const handleSaveVersion = async () => {
        if (!token || !docID || isSavingVersion) return;

        setIsSavingVersion(true);

        try {
            // First persist the current editor content.
            await updateDocument(docID, title, content, token);
            lastSavedRef.current = { title, content };
            setSaveStatus("saved");

            await saveVersion(docID, token);
            toast.success("Version saved successfully!");
        } catch (err) {
            toast.error(err instanceof Error ? `Failed to save version: ${err.message}` : "Failed to save version.");
        } finally {
            setIsSavingVersion(false);
        }
    };
    return (
        <main className="min-h-screen bg-slate-50">
            <header className="sticky top-0 z-20 border-b border-slate-200 bg-white">
                <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 sm:px-6">
                    <div className="flex min-w-0 flex-1 items-center gap-3">
                        <Button
                            variant="ghost"
                            size="icon"
                            aria-label="Back to documents"
                            className="shrink-0 text-slate-500"
                        >
                            <Link href="/documents">
                                <ArrowLeft size={19} />
                            </Link>
                        </Button>

                        <div className="min-w-0 flex-1 sm:flex-none">
                            <Input
                                value={title}
                                onChange={(event) => setTitle(event.target.value)}
                                aria-label="Document title"
                                className="h-9 max-w-64 border-transparent px-2 text-sm font-semibold shadow-none hover:border-slate-200 focus-visible:border-input"
                            />
                        </div>
                    </div>

                    <div className="flex shrink-0 items-center gap-2">
                        <span
                            className={`hidden text-xs sm:inline ${
                                saveStatus === "error"
                                    ? "text-red-600"
                                    : "text-slate-400"
                            }`}
                            role="status"
                            aria-live="polite"
                        >
                            {saveStatus === "saving"
                                ? "Changes pending..."
                                : saveStatus === "saved"
                                  ? "Saved"
                                  : "Save failed"}
                        </span>

                        <div
                            className="hidden items-center -space-x-2 sm:flex"
                            aria-label="Document collaborators"
                        >
                            <Avatar className="h-8 w-8 border-2 border-white">
                                <AvatarFallback className="bg-blue-100 text-xs font-semibold text-blue-700">
                                    AR
                                </AvatarFallback>
                            </Avatar>

                            <Avatar className="h-8 w-8 border-2 border-white">
                                <AvatarFallback className="bg-purple-100 text-xs font-semibold text-purple-700">
                                    SK
                                </AvatarFallback>
                            </Avatar>

                            <Avatar className="h-8 w-8 border-2 border-white">
                                <AvatarFallback className="bg-green-100 text-xs font-semibold text-green-700">
                                    MJ
                                </AvatarFallback>
                            </Avatar>
                        </div>

                        <Button
                            type="button"
                            variant="outline"
                            onClick={handleSaveVersion}
                            disabled={isSavingVersion || !isDocumentLoaded}
                            className="gap-2"
                        >
                            <Save size={16} />
                            <span>
                                {isSavingVersion ? "Saving..." : "Save Version"}
                            </span>
                        </Button>
                        <Button
                            type="button"
                            onClick={() => toast.info("Document sharing coming soon!")}
                            className="gap-2 bg-blue-600 hover:bg-blue-700"
                        >
                            <Share2 size={16} />
                            <span className="hidden sm:inline">Share</span>
                        </Button>

                        <DropdownMenu>
                            <DropdownMenuTrigger className="text-slate-500">
                                    <MoreHorizontal size={19} />
                            </DropdownMenuTrigger>

                            <DropdownMenuContent align="end">
                                <DropdownMenuItem onSelect={() => toast.info("Rename using the title field above.")}>
                                    Rename document
                                </DropdownMenuItem>

                                <DropdownMenuItem onClick={() => setIsHistoryOpen(true)}>
                                    <History className="mr-2 h-4 w-4" />
                                    Version history
                                </DropdownMenuItem>

                                <DropdownMenuItem onSelect={() => window.print()}>
                                    Print document
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                </div>
            </header>

            {/* Writing area */}
            <div className="mx-auto max-w-5xl px-3 py-6 sm:px-8 sm:py-10">
                <div className="mx-auto min-h-[70vh] max-w-3xl rounded-xl border border-slate-200 bg-white px-6 py-8 shadow-sm sm:px-12 sm:py-12">
                    <RichTextEditor
                        key={docID}
                        initialContent={content}
                        onChange={setContent}
                    />
                </div>

                <div className="mx-auto mt-4 flex max-w-3xl flex-wrap items-center justify-between gap-2 px-1 text-xs text-slate-400">
                    <span>{wordCount} words</span>
                </div>
            </div>

            {token && (
                <VersionHistoryDialog
                    docID={docID}
                    token={token}
                    open={isHistoryOpen}
                    onOpenChange={setIsHistoryOpen}
                    onRestored={(restoredDocument) => {
                        setTitle(restoredDocument.title);
                        setContent(restoredDocument.content ?? "");

                        lastSavedRef.current = {
                            title: restoredDocument.title,
                            content: restoredDocument.content ?? "",
                        };
                        setSaveStatus("saved");
                    }}
                />
            )}
        </main>
    );
}