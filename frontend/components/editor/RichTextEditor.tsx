"use client";

import { EditorContent, useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { Bold, Italic, UnderlineIcon, Heading1, Heading2, List, ListOrdered, Link2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Separator } from "@/components/ui/separator";
import { useEffect } from "react";

type RichTextEditorProps = {
    initialContent: string;
    onChange: (content: string) => void;
};

export default function RichTextEditor({
    initialContent,
    onChange,
}: RichTextEditorProps) {
    const editor = useEditor({
        extensions: [
            StarterKit,
        ],
    content: initialContent,
    immediatelyRender: false,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
    editorProps: {
      attributes: {
        class:
          "min-h-[60vh] w-full outline-none text-base leading-8 text-slate-700 prose prose-slate max-w-none",
        },
    },
    });

    useEffect(() => {
        if (!editor) return;

        if (editor.getHTML() !== initialContent) {
            editor.commands.setContent(initialContent || "");
        }
    }, [editor, initialContent]);

    if (!editor) {
        return (
        <div className="min-h-[60vh] text-sm text-slate-400">
            Loading editor...
        </div>
        );
    }

    const toolbarButtonClass = (active: boolean) =>
        `h-9 w-9 shrink-0 ${
        active
            ? "bg-blue-100 text-blue-700 hover:bg-blue-100 hover:text-blue-700"
            : "text-slate-600"
        }`;

    const setLink = () => {
        const previousUrl = editor.getAttributes("link").href;
        const url = window.prompt("Enter a URL", previousUrl || "https://");

        if (url === null) return;

        if (url.trim() === "") {
            editor.chain().focus().extendMarkRange("link").unsetLink().run();
            return;
        }

    editor
        .chain()
        .focus()
        .extendMarkRange("link")
        .setLink({ href: url.trim() })
        .run();
    };

    const toolbarItems = [
        {
            label: "Bold",
            icon: Bold,
            active: editor.isActive("bold"),
            action: () => editor.chain().focus().toggleBold().run(),
            canRun: editor.can().chain().focus().toggleBold().run(),
        },
        {
            label: "Italic",
            icon: Italic,
            active: editor.isActive("italic"),
            action: () => editor.chain().focus().toggleItalic().run(),
            canRun: editor.can().chain().focus().toggleItalic().run(),
        },
        {
            label: "Underline",
            icon: UnderlineIcon,
            active: editor.isActive("underline"),
            action: () => editor.chain().focus().toggleUnderline().run(),
            canRun: editor.can().chain().focus().toggleUnderline().run(),
        },
        {
            label: "Heading 1",
            icon: Heading1,
            active: editor.isActive("heading", { level: 1 }),
            action: () =>
                editor.chain().focus().toggleHeading({ level: 1 }).run(),
            canRun: editor.can().chain().focus().toggleHeading({ level: 1 }).run(),
        },
        {
            label: "Heading 2",
            icon: Heading2,
            active: editor.isActive("heading", { level: 2 }),
            action: () =>
                editor.chain().focus().toggleHeading({ level: 2 }).run(),
            canRun: editor.can().chain().focus().toggleHeading({ level: 2 }).run(),
        },
        {
            label: "Bulleted list",
            icon: List,
            active: editor.isActive("bulletList"),
            action: () => editor.chain().focus().toggleBulletList().run(),
            canRun: editor.can().chain().focus().toggleBulletList().run(),
        },
        {
            label: "Numbered list",
            icon: ListOrdered,
            active: editor.isActive("orderedList"),
            action: () => editor.chain().focus().toggleOrderedList().run(),
            canRun: editor.can().chain().focus().toggleOrderedList().run(),
        },
    ];

    return (
    <div>
        <TooltipProvider>
            <div
                role="toolbar"
                aria-label="Text formatting"
                className="mb-5 flex flex-wrap items-center gap-1 border-b border-slate-100 pb-3"
            >
                {toolbarItems.slice(0, 3).map((item) => {
                    const Icon = item.icon;

                    return (
                    <Tooltip key={item.label}>
                        <TooltipTrigger>
                            <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                aria-label={item.label}
                                aria-pressed={item.active}
                                disabled={!item.canRun}
                                onClick={item.action}
                                className={toolbarButtonClass(item.active)}
                            >
                                <Icon size={17} />
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>{item.label}</TooltipContent>
                    </Tooltip>
                    );
                })}
                <Separator orientation="vertical" className="mx-2 h-5" />

                {toolbarItems.slice(3, 5).map((item) => {
                    const Icon = item.icon;

                    return (
                    <Tooltip key={item.label}>
                        <TooltipTrigger>
                            <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                aria-label={item.label}
                                aria-pressed={item.active}
                                disabled={!item.canRun}
                                onClick={item.action}
                                className={toolbarButtonClass(item.active)}
                            >
                                <Icon size={18} />
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>{item.label}</TooltipContent>
                    </Tooltip>
                    );
                })}
                <Separator orientation="vertical" className="mx-2 h-5" />

                {toolbarItems.slice(5).map((item) => {
                    const Icon = item.icon;

                    return (
                    <Tooltip key={item.label}>
                        <TooltipTrigger>
                            <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                aria-label={item.label}
                                aria-pressed={item.active}
                                disabled={!item.canRun}
                                onClick={item.action}
                                className={toolbarButtonClass(item.active)}
                            >
                                <Icon size={17} />
                            </Button>
                        </TooltipTrigger>
                        <TooltipContent>{item.label}</TooltipContent>
                    </Tooltip>
                    );
                })}
            </div>
        </TooltipProvider>

        <EditorContent editor={editor} />
    </div>
    );
}