import { useUserStore } from "@/stores/userStore";
import ExperienceSection from "@/components/profile/ExperienceSection";

export default function ProfilePage() {
  const { user } = useUserStore();

  return (
    <div className="container mx-auto py-10 space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-muted-foreground">
          Manage your public profile and work experience.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left Column: User Info (Placeholder for now) */}
        <div className="md:col-span-1 space-y-6">
          <div className="bg-card rounded-lg border p-6 shadow-sm">
            <div className="flex flex-col items-center space-y-4">
              <div className="h-24 w-24 rounded-full bg-muted flex items-center justify-center text-2xl font-bold">
                {user?.firstName?.[0]}
                {user?.lastName?.[0]}
              </div>
              <div className="text-center">
                <h2 className="text-xl font-semibold">
                  {user?.firstName} {user?.lastName}
                </h2>
                <p className="text-sm text-muted-foreground">{user?.email}</p>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Details */}
        <div className="md:col-span-2 space-y-6">
          <ExperienceSection />
          {/* Skills and Projects will go here */}
        </div>
      </div>
    </div>
  );
}
