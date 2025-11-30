import { useUserStore } from "@/stores/userStore";

export default function ProfilePage() {
  const { user } = useUserStore();

  return (
    <div>
      <div>Hi welcome</div>
      {user && (
        <ul>
          {Object.entries(user).map(([key, val]) => (
            <li>
              {key}: {String(val)}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
