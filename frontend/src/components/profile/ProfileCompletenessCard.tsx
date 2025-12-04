import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { useUserStore } from "@/stores/userStore";
import {
  calculateProfileCompleteness,
  getMissingFields,
} from "@/utils/profileUtils";
import { useMemo } from "react";

export default function ProfileCompletenessCard() {
  const { user } = useUserStore();

  const completeness = useMemo(
    () => calculateProfileCompleteness(user),
    [user]
  );
  const missingFields = useMemo(() => getMissingFields(user), [user]);

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="text-base">Profile Completeness</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="font-medium">{completeness}%</span>
          </div>
          <Progress value={completeness} className="h-2" />
        </div>
        <p className="text-xs text-muted-foreground">
          Complete your profile to increase your visibility to companies.
        </p>
        {missingFields.length > 0 && completeness < 100 && (
          <div className="pt-2 border-t">
            <p className="text-xs font-medium mb-2">Missing:</p>
            <div className="flex flex-wrap gap-1">
              {missingFields.slice(0, 5).map((field, index) => (
                <span
                  key={index}
                  className="text-xs px-2 py-1 bg-muted rounded-md text-muted-foreground"
                >
                  {field}
                </span>
              ))}
              {missingFields.length > 5 && (
                <span className="text-xs px-2 py-1 bg-muted rounded-md text-muted-foreground">
                  +{missingFields.length - 5} more
                </span>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
