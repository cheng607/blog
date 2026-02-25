import { http } from "../utils/request";

export const getUser = () => {
    return http.get("/user")
}