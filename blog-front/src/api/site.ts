import { http } from '../utils/request'

export const getSiteSummary = () => http.get('/site/summary')
