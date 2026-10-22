"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { FileText } from "lucide-react";
import { login } from "../lib/auth";
import { useAuth } from "../lib/authContext";
import FormField from "@/components/ui/FormField";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function LoginPage() {
    const router = useRouter();
    const { loginSession } = useAuth();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const [errors, setErrors] = useState({
        email: "",
        password: "",
    });

    function validateForm() {
        const nextErrors = {
            email: "",
            password: "",
        };

        if (!email.trim()) {
            nextErrors.email = "Email is required";
        } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
            nextErrors.email = "Enter a valid email address";
        }

        if (!password) {
            nextErrors.password = "Password is required";
        } else if (password.length < 8) {
            nextErrors.password = "Password must be at least 8 characters";
        }

        setErrors(nextErrors);

        return !nextErrors.email && !nextErrors.password;
    }

    const handleSubmit: React.ComponentProps<"form">["onSubmit"] = async (event) => {
        if (!event) return;

        event.preventDefault();
        
        if (!validateForm()) return;
        
        setError("");
        setLoading(true);

        try {
            const result = await login({ email, password });
            loginSession(result);
            toast.success("Welcome back to Scribble!");
            router.replace("/");
        }
        catch (err) {
            const message = err instanceof Error ? err.message : "Unable to log in.";
            toast.error(message);
        }
        finally {
            setLoading(false);
        }
    }

    return (
        <main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
            <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
                <div className="flex items-center justify-center text-sm text-slate-500">
                    <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-600 text-white">
                        <FileText size={23} />
                    </div>
                    <span className="text-xl px-3 py-6 font-semibold text-slate-900">
                        Scribble
                    </span>
                </div>

                <h1 className="text-2xl font-semibold text-slate-900">
                    Welcome back!
                </h1>
                <p className="mt-2 text-sm text-slate-500">
                    Sign in to continue to your documents.
                </p>

                <form onSubmit={handleSubmit} noValidate className="mt-8 space-y-5">
                    <div>
                        <FormField
                            id="email"
                            label="Email"
                            type="email"
                            autoComplete="email"
                            placeholder="you@example.com"
                            value={email}
                            onChange={(event) => {
                                setEmail(event.target.value);
                                setErrors((previous) => ({ ...previous, email: "" }));
                            }}
                            error={errors.email}
                        />
                    </div>

                    <div>
                        <FormField
                            id="password"
                            label="Password"
                            type="password"
                            autoComplete="current-password"
                            placeholder="Enter your password"
                            value={password}
                            onChange={(event) => {
                                setPassword(event.target.value);
                                setErrors((previous) => ({ ...previous, password: "" }));
                            }}
                            error={errors.password}
                        />
                    </div>

                    {/* {error && (
                        <AppAlert
                            title="Login failed"
                            message={error}
                            severity="error"
                            onClose={() => setError("")}
                        />
                    )} */}

                    <Button
                        type="submit"
                        disabled={loading}
                        className="w-full"
                        >
                        {loading ? "Signing in..." : "Sign in"}
                    </Button>
                </form>

                <p className="mt-6 text-center text-sm text-slate-500">
                    Don&apos;t have an account?{" "}
                    <Link
                        href="/signup"
                        className="font-medium text-indigo-600 hover:text-indigo-700"
                    >
                        Create an account
                    </Link>
                </p>
            </div>
        </main>
    );
}