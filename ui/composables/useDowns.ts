export function useDowns(endpoint: string) {
	const config = useRuntimeConfig()
	return `${config.public.apiUrl}/downs/${endpoint}`
};
