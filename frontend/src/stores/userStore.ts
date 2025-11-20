import type { User } from '@/utils/types/user.type';
import {create} from 'zustand';

interface UserState {
    user: User | null;
    setUser:(user:User | Partial<User> | null, key?: keyof User, value?:unknown) => void;
    isAuth: boolean;
    setAuth: (b:boolean) => void;
    // clearAuth: () => Promise<void>;
}

export const useUserStore = create<UserState>((set) => ({
    user: null,
    
    setUser:(userOrPartial, key, value):void => {
        if (userOrPartial && typeof userOrPartial === 'object' && !key) {
            set({ user: userOrPartial as User})
        }

        else if(key){
            set((state) => ({
                user: state.user
                ? { ...state.user, [key]:value}
                : state.user
            }))
        }
        else if(userOrPartial === null){
            set({ user:null })
        }
    },

    isAuth:false,
    setAuth: (b: boolean) => {
        set({
            isAuth: b
        })
    },
    // clearAuth: async () => {
        
    // },

}))