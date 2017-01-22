import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';

import { Observable } from 'rxjs/Rx';
import { TypedJSON } from 'typedjson-npm';

import { APIResponse } from '../common/api-response';

import { ConfigService } from '../common/config.service';


@Injectable()
export class AuthHttp {

    constructor(
        private http: Http,
        private config: ConfigService
    ) { }

    // TODO: no session token? Error before even making the request, if AuthHttp is used, the request will not succeed without session token.
    // createAuthorizationHeader(headers: Headers) {
    //     headers.append('x-gp-sessiontoken', this.auth.getSessionToken());
    // }

    get<T>(clazz: { new (): T }, url: string): Observable<T> {
        let headers = new Headers();
        // this.createAuthorizationHeader(headers);
        return this.http.get(this.config.apiURL + "/backoffice-api/" + url, { headers })
            .map((res: Response) => TypedJSON.parse(res.text(), clazz));
    }

    put(url: string, data: any): Observable<APIResponse> {
        let headers = new Headers();
        // this.createAuthorizationHeader(headers);
        return this.http.put(this.config.apiURL + "/backoffice-api/" + url, data, { headers })
            .map((res: Response) => TypedJSON.parse(res.text(), APIResponse));
    }

    post(url: string, data: any): Observable<APIResponse> {
        let headers = new Headers();
        // this.createAuthorizationHeader(headers);
        return this.http.post(this.config.apiURL + "/backoffice-api/" + url, data, { headers })
            .map((res: Response) => TypedJSON.parse(res.text(), APIResponse));
    }
}
