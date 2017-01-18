import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import { ConfigService } from '../common/config.service';

@Injectable()
export class Auth {

    constructor(private http: Http, private config: ConfigService) { }

    // session is stored internally
    private session: any; // TODO: type
    getSessionToken(): string {
        if (!this.session) {
            return ``;
        }
        return this.session.token;
    }

    // tryLogin tries to obtain a valid session from the server using given credentials
    // in a later stage we'll simply extend this with a user/pass check, probably with SRP.
    tryLogin(id: number) {
        return this.http.post(this.config.apiURL + "/api/login", JSON.stringify({ id: id }))
            .map((data: Response) => {
                this.session = data.json().session;
                return true;
            }).catch((error: any) => {
                return Observable.throw(error._body);
            });

    }
}