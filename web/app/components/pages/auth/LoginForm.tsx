import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Link } from "@tanstack/react-router";
import { FormWrapper } from "./FormWrapper";
import Client from "@/api/accounts";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";

const formSchema = z.object({
	username: z.string(),
	password: z.string(),
});

export function LoginForm({
	className,
	...props
}: React.ComponentPropsWithoutRef<"div">) {
	const a = new Client("http://localhost:4000");
	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			username: "",
			password: "",
		},
	});

	function onSubmit(data: z.infer<typeof formSchema>) {
		a.accounts.Login(data);
	}

	return (
		<FormWrapper className={className} {...props}>
			<Form {...form}>
				<form
					className="flex flex-col gap-4"
					onSubmit={form.handleSubmit(onSubmit)}
				>
					<div className="grid gap-6">
						<FormField
							control={form.control}
							name="username"
							render={({ field }) => (
								<FormItem className="grid gap-2">
									<FormLabel htmlFor="username">Username</FormLabel>
									<FormControl>
										<Input
											id="username"
											type="text"
											placeholder="username"
											required
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<FormField
							control={form.control}
							name="password"
							render={({ field }) => (
								<FormItem className="grid gap-2">
									<div className="flex items-center">
										<Link
											to="/reset"
											className="ml-auto text-sm underline-offset-4 hover:underline"
										>
											Forgot your password?
										</Link>
									</div>
									<FormLabel htmlFor="password">Password</FormLabel>
									<FormControl>
										<Input
											id="password"
											type="password"
											placeholder="password"
											required
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
						<Button type="submit" className="w-full">
							Login
						</Button>
					</div>
					<div className="text-center text-sm">
						Don&apos;t have an account?{" "}
						<Link to="/signup" className="underline underline-offset-4">
							Sign up
						</Link>
					</div>
				</form>
			</Form>
		</FormWrapper>
	);
}
