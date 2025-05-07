import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Link } from "@tanstack/react-router";
import { FormWrapper } from "./FormWrapper";

export function SignupForm({
	className,
	...props
}: React.ComponentPropsWithoutRef<"div">) {
	return (
		<FormWrapper className={className} {...props}>
			<div className="flex flex-col gap-4">
				<div className="grid gap-6">
					<div className="grid gap-2">
						<Label htmlFor="email">Email</Label>
						<Input
							id="email"
							type="email"
							placeholder="m@example.com"
							required
						/>
					</div>
					<div className="grid gap-2">
						<div className="flex items-center">
							<Label htmlFor="password">Password</Label>
						</div>
						<Input id="password" type="password" required />
					</div>
					<Button type="submit" className="w-full">
						Signup
					</Button>
				</div>
				<div className="text-center text-sm">
					Already have an account?{" "}
					<Link to="/login" className="underline underline-offset-4">
						Login
					</Link>
				</div>
			</div>
		</FormWrapper>
	);
}
