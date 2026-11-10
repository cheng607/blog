import { http } from "../utils/request";

export const getUser = () => {
    return http.get("/user")
}
export interface Profile {
    id: number
    username: string
    signature: string
    avatar: string
}
export const saveProfile = (profile: Profile) => profile.id
    ? http.put('/user/update', profile)
    : http.post('/user/create', profile)
