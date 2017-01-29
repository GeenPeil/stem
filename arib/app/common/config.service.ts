import { Injectable } from '@angular/core';

@Injectable()
export class ConfigService {
	public apiURL: string = window.location.protocol + `//` + window.location.host;
}
