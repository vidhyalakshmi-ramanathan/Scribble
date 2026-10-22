"use client";

import Link from "next/link";
import {
    Plus,
    FileText,
    Users,
    Clock3,
    ArrowRight,
    Mail,
    Sparkles,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";

import StatCard from "./StatCard";
import RecentDocuments from "./RecentDocuments";

type DocumentItem = {
    title: string;
    category: string;
    updated: string;
    initials: string;
    color: string;
};

type DashboardContentProps = {
    documents: DocumentItem[];
    searchQuery: string;
    onCreateDocument: () => void;
};

export default function DashboardContent({
    documents,
    searchQuery,
    onCreateDocument,
}: DashboardContentProps) {
    return (
        <main
            id="dashboard"
            className="mx-auto max-w-7xl space-y-8 p-4 sm:p-6 lg:p-8"
        >
            {/* Welcome section */}
            <section className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
                <div>
                    <p className="mb-2 text-sm font-medium text-blue-600">
                        YOUR WORKSPACE
                    </p>

                    <h1 className="text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl">
                        Welcome back!
                    </h1>

                    <p className="mt-2 text-sm text-slate-500 sm:text-base">
                        Your ideas, documents, and team — all in one place. Here's what's happening with your documents. Here's what's happening with your documents.
                    </p>
                </div>

                <Button
                    type="button"
                    onClick={onCreateDocument}
                    className="gap-2 rounded-xl bg-blue-600 px-5 py-3 shadow-lg shadow-blue-600/20 hover:bg-blue-700"
                >
                    <Plus size={18} />
                    Create document
                </Button>
            </section>

            {/* Statistics */}
            <section className="grid grid-cols-1 gap-5 sm:grid-cols-2 xl:grid-cols-3">
                <StatCard
                    title="My documents"
                    value={documents.length}
                    description="Your personal documents"
                    icon={FileText}
                    iconColor="text-blue-600"
                    iconBackground="bg-blue-50"
                />

                <StatCard
                    title="Shared with me"
                    value={8}
                    description="Documents shared with you"
                    icon={Users}
                    iconColor="text-purple-600"
                    iconBackground="bg-purple-50"
                />

                <StatCard
                    title="Recently updated"
                    value={5}
                    description="Updated in the last 7 days"
                    icon={Clock3}
                    iconColor="text-emerald-600"
                    iconBackground="bg-emerald-50"
                />
            </section>

            {/* Main dashboard columns */}
            <section className="grid grid-cols-1 items-start gap-6 xl:grid-cols-3">
                <div className="space-y-6 xl:col-span-2">
                    {/* Recent documents */}
                    <div id="my-documents" className="scroll-mt-24">
                        <RecentDocuments
                            documents={documents}
                            searchQuery={searchQuery}
                        />
                    </div>

                    {/* Shared documents */}
                    <Card
                        id="shared"
                        className="gap-0 rounded-2xl border-slate-200 py-0 shadow-sm"
                    >
                        <CardHeader className="flex flex-row items-center justify-between gap-3 px-5 pt-5">
                           <div>
                                <CardTitle className="text-base font-bold text-slate-900">
                                    Shared with me
                                </CardTitle>

                                <p className="mt-1 text-sm text-slate-500">
                                    Documents your teammates shared with you
                                </p>
                            </div>

                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-purple-50 text-purple-600">
                                <Users size={20} />
                            </div>
                        </CardHeader>

                        <CardContent className="p-5">
                            <div className="flex items-center gap-3 rounded-xl bg-slate-50 p-4">
                                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-white text-slate-500 shadow-sm">
                                    <FileText size={19} />
                                </div>

                                <div className="min-w-0 flex-1">
                                    <p className="truncate text-sm font-semibold text-slate-800">
                                        Team meeting notes
                                    </p>

                                    <p className="mt-1 text-xs text-slate-500">
                                        Shared by Alex Rivera
                                    </p>
                                </div>

                                <Badge
                                    variant="secondary"
                                    className="shrink-0 bg-purple-100 text-purple-700 hover:bg-purple-100"
                                >
                                    Editor
                                </Badge>
                            </div>

                            <div className="flex flex-row items-center justify-center bg-slate-50">
                                <Link href="/shared" className="flex flex-row items-center gap-3 mt-2 h-auto text-sm font-semibold text-blue-600">
                                    Browse shared documents
                                    <ArrowRight size={16} />
                                </Link>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Right column */}
                <div className="space-y-6">
                    {/* Invitations */}
                    <Card className="gap-0 rounded-2xl border-slate-200 py-0 shadow-sm">
                        <CardHeader className="flex flex-row items-center gap-3 px-5 pt-5">
                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-orange-50 text-orange-600">
                                <Mail size={20} />
                            </div>

                            <div>
                                <CardTitle className="text-base font-bold text-slate-900">
                                    Invitations
                                </CardTitle>

                                <p className="mt-1 text-sm text-slate-500">
                                    Pending document invites
                                </p>
                            </div>
                        </CardHeader>

                        <CardContent className="p-5">
                            <div className="rounded-xl border border-dashed border-slate-200 px-4 py-6 text-center">
                                <p className="text-sm font-medium text-slate-700">
                                    You're all caught up!
                                </p>

                                <p className="mt-1 text-xs text-slate-500">
                                    New invitations will appear here.
                                </p>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Collaboration promo */}
                    <Card className="gap-0 overflow-hidden rounded-2xl border-0 bg-gradient-to-br from-blue-600 to-indigo-700 py-0 text-white shadow-sm">
                        <CardContent className="p-5">
                            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white/15">
                                <Sparkles size={21} />
                            </div>

                            <h2 className="mt-4 text-lg font-bold">Better together</h2>

                            <p className="mt-2 text-sm leading-6 text-blue-100">
                                Invite your teammates and work on documents together in real time.
                            </p>

                            <Button
                                type="button"
                                onClick={onCreateDocument}
                                variant="secondary"
                                className="mt-5 gap-2 rounded-xl bg-white text-blue-700 hover:bg-blue-50"
                            >
                                Get started
                                <ArrowRight size={16} />
                            </Button>
                        </CardContent>
                    </Card>
                </div>
            </section>
        </main>
    );
}