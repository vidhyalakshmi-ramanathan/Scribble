"use client";

import { useState } from "react";
import type { ComponentProps } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { FileText } from "lucide-react";
import { signup } from "../lib/auth";
import FormField from "@/components/ui/FormField";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function SignupPage() {
  const router = useRouter();

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({
    username: "",
    email: "",
    password: "",
  });
  
  function validateForm() {
    const nextErrors = {
      username: "",
      email: "",
      password: "",
    };

    if (!username.trim()) {
      nextErrors.username = "Username is required";
    } 

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

    return !nextErrors.username && !nextErrors.email && !nextErrors.password;
    }

  const handleSubmit: ComponentProps<"form">["onSubmit"] = async (event) => {
    event.preventDefault();

    if (!validateForm()) return;

    setLoading(true);

    try {
      await signup({ username, email, password });
      toast.success("Account created successfully!")
      router.push("/login");
    }
    catch (err) {
      const message = err instanceof Error ? err.message : "Unable to create your account.";
      toast.error(message);
    }
    finally {
      setLoading(false);
    }
  };

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
          Create your account
        </h1>
        <p className="mt-2 text-sm text-slate-500">
          Sign up to start creating and collaborating on documents.
        </p>

        <form onSubmit={handleSubmit} noValidate className="mt-8 space-y-5">
          <div>
            <FormField
              id="username"
              label="Username"
              type="text"
              autoComplete="username"
              placeholder="Enter your name"
              value={username}
              onChange={(event) => {
                  setUsername(event.target.value);
                  setErrors((previous) => ({ ...previous, username: "" }));
              }}
              error={errors.username}
             />
          </div>

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
              autoComplete="new-password"
              placeholder="Enter your password"
              value={password}
              onChange={(event) => {
                  setPassword(event.target.value);
                  setErrors((previous) => ({ ...previous, password: "" }));
              }}
              error={errors.password}
            />
          </div>

          <Button
            type="submit"
            disabled={loading}
            className="w-full"
          >
            {loading ? "Creating account..." : "Create account"}
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-500">
          Already have an account?{" "}
          <Link href="/login" className="font-medium text-indigo-600 hover:text-indigo-700">
            Sign in
          </Link>
        </p>
      </div>
    </main>
  );
}