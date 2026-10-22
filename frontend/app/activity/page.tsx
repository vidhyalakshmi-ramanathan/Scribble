"use client";

import { useAppShell } from "@/components/layouts/AppShell";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { Pencil, Share2, Plus, Activity, Search, Clock3, History } from "lucide-react";

const activities = [
  {
    id: 1,
    action: "edited",
    document: "Project Architecture",
    detail: "Updated the API design section",
    person: "You",
    initials: "YO",
    time: "10 minutes ago",
    color: "bg-blue-100 text-blue-700",
    icon: "edit",
  },
  {
    id: 2,
    action: "shared",
    document: "Sprint Planning",
    detail: "Shared the document with you",
    person: "Sarah Kim",
    initials: "SK",
    time: "2 hours ago",
    color: "bg-purple-100 text-purple-700",
    icon: "share",
  },
  {
    id: 3,
    action: "created",
    document: "Product Requirements",
    detail: "Created a new document",
    person: "You",
    initials: "YO",
    time: "Yesterday",
    color: "bg-green-100 text-green-700",
    icon: "create",
  },
  {
    id: 4,
    action: "restored a version of",
    document: "Project Architecture",
    detail: "Restored an earlier document version",
    person: "Alex Rivera",
    initials: "AR",
    time: "2 days ago",
    color: "bg-amber-100 text-amber-700",
    icon: "history",
  },
];

const activityIcons = {
    edit: Pencil,
    share: Share2,
    create: Plus,
    history: History,
} as const;

export default function ActivityPage() {
    const { searchQuery, setSearchQuery } = useAppShell();

    const filteredActivities = activities.filter((item) =>
        `${item.action} ${item.document} ${item.detail} ${item.person}`
        .toLowerCase()
        .includes(searchQuery.trim().toLowerCase()),
    );

    return (
    <main className="min-h-screen p-5 sm:p-8">
        <div className="mx-auto max-w-5xl">
            <div className="mb-8">
                <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-xl bg-blue-100 text-blue-700">
                    <Activity size={22} />
                </div>

                <h1 className="text-2xl font-bold text-slate-800">
                    Recent Activity
                </h1>

                <p className="mt-2 text-sm text-slate-500">
                    Keep track of changes and activity across your documents.
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
                    placeholder="Search activity..."
                    aria-label="Search recent activity"
                    className="h-12 border-slate-200 bg-white pl-11"
                />
            </div>

            <div className="mb-4 flex items-center justify-between">
                <h2 className="font-semibold text-slate-800">
                    Activity timeline
                </h2>

                <Badge variant="secondary">
                    {filteredActivities.length}{" "}
                    {filteredActivities.length === 1 ? "event" : "events"}
                </Badge>
            </div>

            <Card className="rounded-xl border-slate-200 py-0 shadow-sm">
                <CardHeader className="sr-only">
                    <CardTitle>Activity timeline</CardTitle>
                </CardHeader>

                <CardContent className="p-5 sm:p-7">
                    {filteredActivities.length > 0 ? (
                        <div className="space-y-7">
                            {filteredActivities.map((item, index) => {
                                const Icon = activityIcons[ item.icon as keyof typeof activityIcons ];

                                return (
                                    <div key={item.id} className="relative flex gap-4">
                                        {index < filteredActivities.length - 1 && (
                                            <div className="absolute bottom-[-28px] left-5 top-11 w-px bg-slate-200" />
                                        )}
                                        <div
                                            className={`relative z-10 flex h-10 w-10 shrink-0 items-center justify-center rounded-full ${item.color}`}
                                        >
                                            <Icon size={18} />
                                        </div>

                                        <div className="min-w-0 flex-1 pt-0.5">
                                            <p className="text-sm leading-6 text-slate-700">
                                                <span className="font-medium text-slate-900">
                                                    {item.person}
                                                </span>{" "}
                                                {item.action}{" "}
                                                <span className="font-medium text-slate-900">
                                                    {item.document}
                                                </span>
                                            </p>
                                            <p className="mt-1 text-sm text-slate-500">
                                                {item.detail}
                                            </p>

                                            <div className="mt-2 flex items-center gap-1.5 text-xs text-slate-400">
                                                <Clock3 size={13} />
                                                <span>{item.time}</span>
                                            </div>
                                        </div>

                                        <div
                                            className={`hidden h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-semibold sm:flex ${item.color}`}
                                            aria-label={`Activity by ${item.person}`}
                                        >
                                            {item.initials}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    ) : (
                    <div className="flex flex-col items-center py-10 text-center">
                        <Activity
                            size={30}
                            className="text-slate-400"
                        />
                        <h3 className="mt-3 text-sm font-semibold text-slate-700">
                            No activity found
                        </h3>
                        <p className="mt-1 text-sm text-slate-500">
                            Try searching for another document or activity.
                        </p>
                    </div>
                    )}
                </CardContent>
            </Card>
        </div>
    </main>
    );
}