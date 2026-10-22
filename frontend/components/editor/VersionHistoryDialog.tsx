"use client";

import { useEffect, useState } from "react";
import { History, RotateCcw, Clock, FileText, X } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { DocumentVersion, getVersions, restoreVersion, getDocument } from "@/app/lib/doc";

type VersionHistoryDialogProps = {
    docID: string;
    token: string;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onRestored: (document: Document) => void;
};

export default function VersionHistoryDialog({
    docID,
    token,
    open,
    onOpenChange,
    onRestored,
}: VersionHistoryDialogProps) {
    const [versions, setVersions] = useState<DocumentVersion[]>([]);
    const [selectedVersionID, setSelectedVersionID] = useState<number | null>(
        null
    );
    const [isLoading, setIsLoading] = useState(false);
    const [restoringVersionID, setRestoringVersionID] = useState<number | null>(
        null
    );

    useEffect(() => {
        if (!open) return;

        let cancelled = false;

        async function loadVersions() {
            setIsLoading(true);

            try {
                const data = await getVersions(docID, token);

                if (!cancelled) {
                    // Display the newest version first.
                    const sortedVersions = [...(data ?? [])].sort(
                        (a, b) => b.versionID - a.versionID
                    );

                    setVersions(sortedVersions);
                    setSelectedVersionID(
                        sortedVersions[0]?.versionID ?? null
                    );
                }
            } catch (err) {
                if (!cancelled) {
                    toast.error(
                        err instanceof Error
                            ? `Failed to load versions: ${err.message}`
                            : "Failed to load version history."
                    );
                }
            } finally {
                if (!cancelled) {
                    setIsLoading(false);
                }
            }
        }

        loadVersions();

        return () => {
            cancelled = true;
        };
    }, [open, docID, token]);

    const selectedVersion = versions.find(
        (version) => version.versionID === selectedVersionID
    );

    const handleRestore = async () => {
        if (selectedVersionID === null || restoringVersionID !== null) return;

        const versionID = selectedVersionID;
        setRestoringVersionID(versionID);

        try {
            await restoreVersion(docID, versionID, token);

            const restoredDocument = await getDocument(docID, token);

            onRestored(restoredDocument);
            toast.success(`Restored version ${versionID}.`);
            onOpenChange(false);
        } catch (err) {
            toast.error(
                err instanceof Error
                    ? `Failed to restore version: ${err.message}`
                    : "Failed to restore version."
            );
        } finally {
            setRestoringVersionID(null);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent
                showCloseButton={false}
                className="flex h-[85vh] max-h-[850px] w-[95vw] max-w-6xl flex-col gap-0 overflow-hidden rounded-2xl border border-slate-200 bg-white p-0 text-slate-900 shadow-2xl dark:border-slate-800 dark:bg-slate-950 dark:text-slate-100 sm:max-w-6xl"
            >
                {/* Header */}
                <DialogHeader className="flex shrink-0 flex-row items-center justify-between border-b border-slate-200 px-5 py-4 text-left dark:border-slate-800 sm:px-7">
                    <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-100 dark:bg-slate-800">
                            <History className="h-5 w-5" />
                        </div>

                        <div className="space-y-1">
                            <DialogTitle className="text-lg font-semibold">
                                Version history
                            </DialogTitle>
                            <DialogDescription className="text-xs sm:text-sm">
                                Browse previous snapshots and restore your work.
                            </DialogDescription>
                        </div>
                    </div>

                    <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        aria-label="Close version history"
                        onClick={() => onOpenChange(false)}
                    >
                        <X className="h-5 w-5" />
                    </Button>
                </DialogHeader>

                {/* Main two-panel layout */}
                <div className="grid min-h-0 flex-1 grid-cols-1 md:grid-cols-[280px_minmax(0,1fr)]">
                    {/* Left: Version timeline */}
                    <aside className="flex min-h-0 flex-col border-b border-slate-200 bg-slate-50 dark:border-slate-800 dark:bg-slate-900/50 md:border-b-0 md:border-r">
                        <div className="shrink-0 px-5 pb-3 pt-5">
                            <p className="text-xs font-semibold uppercase tracking-wider text-slate-500">
                                Saved versions
                            </p>
                            <p className="mt-1 text-sm text-slate-500">
                                {versions.length}{" "}
                                {versions.length === 1 ? "snapshot" : "snapshots"}
                            </p>
                        </div>

                        <div className="min-h-0 flex-1 overflow-y-auto px-3 pb-4">
                            {isLoading ? (
                                <p className="px-3 py-6 text-sm text-slate-500">
                                    Loading versions...
                                </p>
                            ) : versions.length === 0 ? (
                                <div className="px-3 py-6">
                                    <History className="mb-3 h-8 w-8 text-slate-400" />
                                    <p className="text-sm font-medium">
                                        No saved versions yet
                                    </p>
                                    <p className="mt-1 text-xs leading-5 text-slate-500">
                                        Click Save Version in your editor to
                                        create your first snapshot.
                                    </p>
                                </div>
                            ) : (
                                <div className="relative space-y-1">
                                    <div className="absolute bottom-5 left-[23px] top-5 w-px bg-slate-200 dark:bg-slate-700" />

                                    {versions.map((version, index) => {
                                        const isSelected =
                                            version.versionID ===
                                            selectedVersionID;

                                        return (
                                            <button
                                                key={version.versionID}
                                                type="button"
                                                onClick={() =>
                                                    setSelectedVersionID(
                                                        version.versionID
                                                    )
                                                }
                                                className={`relative flex w-full gap-3 rounded-xl p-3 text-left transition-colors ${
                                                    isSelected
                                                        ? "bg-white shadow-sm ring-1 ring-slate-200 dark:bg-slate-800 dark:ring-slate-700"
                                                        : "hover:bg-white/70 dark:hover:bg-slate-800/50"
                                                }`}
                                            >
                                                <div className="relative z-10 mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 border-slate-300 bg-slate-50 dark:border-slate-600 dark:bg-slate-900">
                                                    <div
                                                        className={`h-2 w-2 rounded-full ${
                                                            isSelected
                                                                ? "bg-slate-900 dark:bg-white"
                                                                : "bg-slate-300 dark:bg-slate-600"
                                                        }`}
                                                    />
                                                </div>

                                                <div className="min-w-0 flex-1">
                                                    <div className="flex flex-wrap items-center gap-2">
                                                        <span className="text-sm font-semibold">
                                                            Version{" "}
                                                            {version.versionID}
                                                        </span>

                                                        {index === 0 && (
                                                            <span className="rounded-full bg-emerald-100 px-2 py-0.5 text-[10px] font-medium text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">
                                                                Latest
                                                            </span>
                                                        )}
                                                    </div>

                                                    <p className="mt-1 truncate text-xs text-slate-500">
                                                        {version.modifiedByUsername}
                                                    </p>

                                                    <p className="mt-1.5 text-xs leading-5 text-slate-400">
                                                        {new Date(
                                                            version.createdAt
                                                        ).toLocaleString()}
                                                    </p>
                                                </div>
                                            </button>
                                        );
                                    })}
                                </div>
                            )}
                        </div>
                    </aside>

                    {/* Right: Selected version preview */}
                    <section className="flex min-h-0 min-w-0 flex-col">
                        {selectedVersion ? (
                            <>
                                <div className="flex shrink-0 flex-wrap items-center justify-between gap-3 border-b border-slate-200 px-5 py-4 dark:border-slate-800 sm:px-7">
                                    <div className="min-w-0">
                                        <p className="text-xs font-medium text-slate-500">
                                            PREVIEWING VERSION{" "}
                                            {selectedVersion.versionID}
                                        </p>
                                        <h3 className="mt-1 truncate text-lg font-semibold">
                                            {selectedVersion.title ||
                                                "Untitled document"}
                                        </h3>
                                        <p className="mt-1 text-xs text-slate-500">
                                            <Clock className="mr-1 inline h-3.5 w-3.5" />
                                            {new Date(
                                                selectedVersion.createdAt
                                            ).toLocaleString()}
                                        </p>
                                    </div>

                                    <Button
                                        type="button"
                                        onClick={handleRestore}
                                        disabled={
                                            restoringVersionID !== null
                                        }
                                        className="shrink-0"
                                    >
                                        <RotateCcw className="mr-2 h-4 w-4" />
                                        {restoringVersionID ===
                                        selectedVersion.versionID
                                            ? "Restoring..."
                                            : "Restore this version"}
                                    </Button>
                                </div>

                                <div className="min-h-0 flex-1 overflow-y-auto bg-white px-5 py-6 dark:bg-slate-950 sm:px-8 sm:py-8">
                                    <article className="mx-auto max-w-3xl">
                                        <div className="mb-6 flex items-center gap-2 text-xs text-slate-400">
                                            <FileText className="h-4 w-4" />
                                            Read-only preview
                                        </div>

                                        {selectedVersion.content?.trim() ? (
                                            <div
                                                className="prose prose-slate max-w-none break-words dark:prose-invert"
                                                dangerouslySetInnerHTML={{
                                                    __html: selectedVersion.content,
                                                }}
                                            />
                                        ) : (
                                            <p className="text-sm italic text-slate-400">
                                                This version has no content.
                                            </p>
                                        )}
                                    </article>
                                </div>
                            </>
                        ) : (
                            <div className="flex flex-1 flex-col items-center justify-center px-6 text-center">
                                <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-slate-100 dark:bg-slate-800">
                                    <FileText className="h-7 w-7 text-slate-400" />
                                </div>

                                <h3 className="font-medium">
                                    {isLoading
                                        ? "Loading version history..."
                                        : "Select a version"}
                                </h3>
                                <p className="mt-2 max-w-xs text-sm text-slate-500">
                                    Choose a snapshot from the timeline to
                                    preview its document content.
                                </p>
                            </div>
                        )}
                    </section>
                </div>
            </DialogContent>
        </Dialog>
    );
}