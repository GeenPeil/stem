import { Component } from '@angular/core';

import { ConfigService } from './common/config.service';
import { AuthHttp } from './auth/auth-http.service';

@Component({
    template: `
        error: {{error}}<br/>
        data:{{data}}
    `,
})
export class ProfileComponent {
    data: string;
    error: string;

    constructor(private http: AuthHttp, private config: ConfigService) {
        http.get(this.config.apiURL + "/api/private/profile").subscribe(
            data => this.data = JSON.stringify(data),
            error => this.error = error._body
        );
    }
}
