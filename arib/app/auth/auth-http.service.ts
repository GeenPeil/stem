import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';

import { Auth } from './auth.service';
import { ConfigService } from '../common/config.service';


@Injectable()
export class AuthHttp {

    constructor(
        private http: Http,
        private auth: Auth,
        private config: ConfigService
    ) { }

    // TODO: error before even making the request, if AuthHttp is used, the request will not succeed without session token.

    createAuthorizationHeader(headers: Headers) {
        headers.append('x-gp-sessiontoken', this.auth.getSessionToken());
    }

    get(url: string) {
        let headers = new Headers();
        this.createAuthorizationHeader(headers);
        return this.http.get(this.config.apiURL + "/api/" + url, {
            headers: headers
        });
    }

    post(url: string, data: any) {
        let headers = new Headers();
        this.createAuthorizationHeader(headers);
        return this.http.post(this.config.apiURL + "/api/" + url, data, {
            headers: headers
        });
    }
}
