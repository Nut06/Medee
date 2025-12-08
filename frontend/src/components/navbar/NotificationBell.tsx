import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";

export function NotificationBell() {
  // TODO: Replace with actual notification count from store/API
  const notificationCount = 2;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-5 w-5" />
          {notificationCount > 0 && (
            <Badge
              variant="destructive"
              className="absolute -right-1 -top-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs"
            >
              {notificationCount}
            </Badge>
          )}
          <span className="sr-only">การแจ้งเตือน</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80">
        <DropdownMenuLabel>การแจ้งเตือน</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem>
          <div className="flex flex-col gap-1">
            <p className="text-sm font-medium">โปรเจกต์ใหม่ที่เหมาะกับคุณ</p>
            <p className="text-xs text-muted-foreground">5 นาทีที่แล้ว</p>
          </div>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <div className="flex flex-col gap-1">
            <p className="text-sm font-medium">มีข้อความใหม่จากบริษัท ABC</p>
            <p className="text-xs text-muted-foreground">1 ชั่วโมงที่แล้ว</p>
          </div>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem className="text-center text-sm text-primary cursor-pointer">
          ดูการแจ้งเตือนทั้งหมด
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
