"use client";

import { useState } from "react";
import Sidebar from "@/components/layouts/Sidebar";
import Header from "@/components/layouts/Header";
import DashboardContent from "@/components/dashboard/DashboardContents";
import { useAppShell } from "@/components/layouts/AppShell";

const documents = [
  {
    title: "Project Architecture",
    category: "Engineering",
    updated: "10 minutes ago",
    initials: "AR",
    color: "blue",
  },
  {
    title: "Sprint Planning",
    category: "Team Workspace",
    updated: "2 hours ago",
    initials: "SK",
    color: "purple",
  },
  {
    title: "Product Requirements",
    category: "Product",
    updated: "Yesterday",
    initials: "MJ",
    color: "green",
  },
]

export default function Home() {

  const {searchQuery} = useAppShell();

  const handleCreateDocument = () => {
        alert("Create document functionality coming soon!");
    };

  return (
    <>
      <DashboardContent
        documents={documents}
        searchQuery={searchQuery}
        onCreateDocument={handleCreateDocument}
      />
    </>
  );
}
