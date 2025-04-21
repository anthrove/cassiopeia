export class NetworkError extends Error {
	status: number;
	body: any;

	constructor(message: string, status: number, body?: any) {
		super(message);
		this.name = 'NetworkError';
		this.status = status;
		this.body = body;
	}
}
