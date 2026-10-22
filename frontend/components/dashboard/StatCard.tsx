import type { LucideIcon } from "lucide-react";
import { ArrowUpRight } from "lucide-react";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

type StatCardProps = {
  title: string;
  value: string | number;
  description: string;
  icon: LucideIcon;
  iconColor: string;
  iconBackground: string;
};

export default function StatCard({
  title,
  value,
  description,
  icon: Icon,
  iconColor,
  iconBackground,
}: StatCardProps) {
  return (
    <Card className="gap-0 rounded-2xl border-slate-200 py-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
      <CardHeader className="flex flex-row items-start justify-between space-y-0 px-5">
        <div>
          <p className="text-sm font-medium text-slate-500">{title}</p>

          <CardTitle className="mt-3 text-3xl font-bold tracking-tight text-slate-900">
            {value}
          </CardTitle>
        </div>

        <div
          className={`flex h-11 w-11 items-center justify-center rounded-xl ${iconBackground} ${iconColor}`}
        >
          <Icon size={21} />
        </div>
      </CardHeader>

      <CardContent className="mt-4 px-5">
        <div className="flex items-center gap-1 text-xs text-slate-500">
          <ArrowUpRight size={14} className="shrink-0 text-emerald-600" />
          <span>{description}</span>
        </div>
      </CardContent>
    </Card>
  );
}