// Mirrors server/internal/audit/dto.go + pagination.Page.

export interface PageMeta {
	page: number;
	limit: number;
	total: number;
	total_pages: number;
}

export interface Page<T> {
	items: T[];
	meta: PageMeta;
}

export interface AuditEvent {
	id: string;
	action: string;
	actor_id?: string;
	actor_name?: string;
	actor_avatar_url?: string;
	target_type?: string;
	target_id?: string;
	target_name?: string;
	details?: Record<string, unknown>;
	device_id?: string;
	ip?: string;
	created_at: string;
}

export interface ListAuditParams {
	page?: number;
	limit?: number;
	action?: string;
	action_prefix?: string;
	actor_id?: string;
	device_id?: string;
	from?: string;
	to?: string;
}
