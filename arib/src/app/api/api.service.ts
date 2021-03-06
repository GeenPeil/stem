import { Injectable } from "@angular/core";
import { Http, Headers, Response } from "@angular/http";

import { Observable } from "rxjs/Rx";
import { TypedJSON } from "typedjson-npm";

import { CreateResponse } from "./create-response";
import { UpdateResponse } from "./update-response";

import { Auth } from "../auth/auth.service";
import { ConfigService } from "../common/config.service";

@Injectable()
export class Api {
	private apiPath = `/api/`;

	constructor(
		private http: Http,
		private config: ConfigService,
		private auth: Auth
	) { }

	public post(url: string, data: any): Observable<CreateResponse> {
		let headers = new Headers();
		this.createAuthorizationHeader(headers);
		return this.http.post(this.config.apiURL + this.apiPath + url, data, { headers })
			.map((res: Response) => TypedJSON.parse(res.text(), CreateResponse));
	}

	public get<T>(clazz: { new (): T }, url: string): Observable<T> {
		let headers = new Headers();
		this.createAuthorizationHeader(headers);
		return this.http.get(this.config.apiURL + this.apiPath + url, { headers })
			.map((res: Response) => TypedJSON.parse(res.text(), clazz));
	}

	public put(url: string, data: any): Observable<UpdateResponse> {
		let headers = new Headers();
		this.createAuthorizationHeader(headers);
		return this.http.put(this.config.apiURL + this.apiPath + url, data, { headers })
			.map((res: Response) => TypedJSON.parse(res.text(), UpdateResponse));
	}

	// TODO: no session token? Error before even making the request, if Api is used, the request will not succeed without session token.
	private createAuthorizationHeader(headers: Headers) {
		headers.append("x-gp-sessiontoken", this.auth.getSessionToken());
	}
}
