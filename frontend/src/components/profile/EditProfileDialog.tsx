import { Pencil } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

export function EditProfileDialog() {
  return (
    <Link to="/user/edit-profile" className="w-full">
      <Button variant="outline" className="w-full">
        <Pencil className="mr-2 h-4 w-4" />
        Edit Profile
      </Button>
    </Link>
  );
}
