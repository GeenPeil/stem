import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import { ConfigService } from '../common/config.service';

interface Session {
	id: number;
	token: string;
}

@Injectable()
export class Auth {

	constructor(private http: Http, private config: ConfigService) { }

	private session: Session = null;
	private userId = 0;

	getUserId(): number {
		return this.userId;
	}

	getSessionToken(): string {
		return (this.session ? this.session.token : '');
	}

	isLoggedIn(): boolean {
		return !!this.session;
	}

	// tryLogin tries to obtain a valid session from the server using given credentials
	// in a later stage we'll simply extend this with a user/pass check, probably with SRP.
	tryLogin(id: number) {
		return this.http.post(this.config.apiURL + "/api/login", JSON.stringify({ id: id }))
			.map((data: Response) => {
				this.userId = id;
				this.session = data.json().session;
				return true;
			}).catch((error: any) => {
				return Observable.throw(error._body);
			});
	}
	logout() {
		return this.http.post(this.config.apiURL + "/api/logout", JSON.stringify(this.session))
			.map((data: Response) => {
				this.session = null;
				return data.json();
			}).catch((error: any) => {
				return Observable.throw(error._body);
			});
	}
}
