// import { useAuthForm } from "@/hooks/userAuthForm"
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "./ui/form";
import { useForm } from "react-hook-form";
import { Input } from "./ui/input";


export const LoginForm = () => {
    const form = useForm()
    // const { name, setName, password, setPassword, error,  setError } = useAuthForm(); 
    return (
    <Form>
        <FormField
            control={form.control}
            name="username"
            render={({field}) => (
                <FormItem>
                    <FormLabel>Username</FormLabel>
                    <FormControl>
                        <Input placeholder="กรุณากรอกชื่อของท่าน" {...field}></Input>
                    </FormControl>
                    <FormDescription>This is you publicy display name.</FormDescription>
                    <FormMessage/>
                </FormItem>
            )}
        />
    </Form>
  )
}
