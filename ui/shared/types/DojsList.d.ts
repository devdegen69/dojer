export type Doj = {
	id: string
	name: string
	title: string
	parodies: string
	characters: string
	tags: string
	artists: string
	groups: string
	languages: string
	categories: string
	pages: number
	uploaded: string
	createdAt: string
}

export type Pagination = {
	currentPage: number
	totalPages: number
	totalResults: number
	pageSize: number
	pages: number[]
}

export type DojsList = {
	dojs: Doj[]
	pagination: Pagination
}

export type Reader = {
	doj: Doj
	counters: {
		[field: string]: {
			[name: string]: number
		}
	}
	images: string[]
	sizes: string[]
	title: string
}
