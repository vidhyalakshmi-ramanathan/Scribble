import type { InputHTMLAttributes } from "react";
import { CircleAlert } from "lucide-react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

type FormFieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  error?: string;
};

export default function FormField({
    label,
    error,
    id,
    className = "",
    ...inputProps
}: FormFieldProps) {
    return (
        <div className="space-y-2">
            <Label htmlFor={id}>
                {label}
            </Label>

            <Input
                {...inputProps}
                id={id}
                aria-invalid={Boolean(error)}
                aria-describedby={error ? `${id}-error` : undefined}
                className={`${error ? "border-red-500 focus-visible:ring-red-500" : ""} ${className}`}
            />

            {error && (
                <p
                    id={`${id}-error`}
                    role="alert"
                    className="flex items-center gap-1 text-sm text-red-600"
                >
                    <CircleAlert size={14} />
                    {error}
                </p>
            )}
        </div>
    );
}