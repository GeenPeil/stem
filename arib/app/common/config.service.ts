import { Injectable } from '@angular/core';

@Injectable()
export class ConfigService {
	public apiURL: string = `http://localhost:8002`;
}
