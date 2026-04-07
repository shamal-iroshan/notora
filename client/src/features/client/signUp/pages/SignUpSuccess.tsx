import { Link } from "react-router-dom";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AlertCircle } from "lucide-react";

export function SignUpSuccessPage() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <Card className="border-primary/20">
          <CardHeader className="text-center">
            <div className="flex justify-center mb-4">
              <div className="rounded-full bg-primary/10 p-3">
                <AlertCircle className="w-6 h-6 text-primary" />
              </div>
            </div>

            <CardTitle className="text-2xl">Account Created!</CardTitle>
            <CardDescription>
              Your account is pending admin approval
            </CardDescription>
          </CardHeader>

          <CardContent>
            <div className="flex flex-col gap-4">
              <p className="text-sm text-muted-foreground">
                Thank you for signing up! Your account has been created and is
                now pending approval from an administrator. You will receive an
                email notification once your account has been approved.
              </p>

              <div className="bg-muted/50 border border-border rounded-lg p-3 text-xs">
                <p className="font-medium mb-1">What happens next?</p>
                <ul className="space-y-1 text-muted-foreground">
                  <li>• Admin will review your account</li>
                  <li>• You&apos;ll receive an approval email</li>
                  <li>• Then you can log in and start taking notes</li>
                </ul>
              </div>

              <div className="text-center text-sm space-y-2">
                <p className="text-muted-foreground">
                  Try the demo account to explore the app:
                </p>
                <p className="font-mono text-xs bg-muted p-2 rounded">
                  user@example.com
                </p>
                <p className="font-mono text-xs bg-muted p-2 rounded">
                  password123
                </p>
              </div>

              <div className="text-center text-sm pt-2">
                <Link
                  to="/auth/login"
                  className="text-primary underline underline-offset-4 hover:text-primary/80"
                >
                  Back to login
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
