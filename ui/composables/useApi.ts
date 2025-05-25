export function useApi<Result = undefined, Err = undefined>(endpoint: string, options = {}) {
	const config = useRuntimeConfig()
	const apiBaseUrl = config.public.apiUrl + "/api"
	return useFetch<Result, Err>(`${apiBaseUrl}${endpoint}`, options);
};
