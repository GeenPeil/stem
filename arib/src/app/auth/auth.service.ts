import { Injectable } from "@angular/core";
import { Http, Headers, Response } from "@angular/http";
import { Observable } from "rxjs/Observable";

import { ConfigService } from "../common/config.service";

interface Session {
	id: number;
	token: string;
}

@Injectable()
export class Auth {

	private session: Session = null;
	private userId = 0;

	constructor(private http: Http, private config: ConfigService) { }

	public getUserId(): number {
		return this.userId;
	}

	public getSessionToken(): string {
		return (this.session ? this.session.token : "");
	}

	public isLoggedIn(): boolean {
		return !!this.session;
	}

	// tryLogin tries to obtain a valid session from the server using given credentials
	// in a later stage we"ll simply extend this with a user/pass check, probably with SRP.
	public tryLogin(id: number) {
		return this.http.post(this.config.apiURL + "/api/login", JSON.stringify({ id }))
			.map((data: Response) => {
				this.userId = id;
				this.session = data.json().session;
				return true;
			}).catch((error: any) => {
				return Observable.throw(error._body);
			});
	}

	public logout() {
		return this.http.post(this.config.apiURL + "/api/logout", JSON.stringify(this.session))
			.map((data: Response) => {
				this.session = null;
				return data.json();
			}).catch((error: any) => {
				return Observable.throw(error._body);
			});
	}
}
