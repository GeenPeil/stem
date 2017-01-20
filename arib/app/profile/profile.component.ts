import { Component } from '@angular/core';

import { AuthHttp } from '../auth/auth-http.service';

@Component({
    template: `
        error: {{error}}<br/>
        data:{{data}}
    `,
})
export class ProfileComponent {
    data: string;
    error: string;

    constructor(private http: AuthHttp) {
        http.get("private/profile").subscribe(
            data => this.data = JSON.stringify(data),
            error => this.error = error._body
        );
    }
}
