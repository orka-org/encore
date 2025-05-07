import { SignupForm } from "@/components/pages/auth/SignupForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(auth)/signup")({
	component: RouteComponent,
});

function RouteComponent() {
	return <SignupForm />;
}
