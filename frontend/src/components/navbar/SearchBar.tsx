import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";

export function SearchBar() {
  return (
    <div className="relative hidden md:block md:w-80">
      <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        type="search"
        placeholder="ค้นหาโปรเจกต์..."
        className="pl-10 bg-muted/50"
      />
    </div>
  );
}
