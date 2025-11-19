import { useState } from "react";

export function useAuthForm() {
    const [name, setName] = useState<string>('')
    const [password, setPassword] = useState<string>('')
    const [error, setError] = useState<string>('')
    return { name, setName, password, setPassword, error,  setError}
}