import { Component } from '@angular/core';

import { Auth } from '../auth/auth.service';
import { AuthHttp } from '../auth/auth-http.service';

@Component({
    template: `
        <div *ngIf="error">
            <h1>Could not load your profile</h1>
            <p>error: {{error}}</p>
            <button (click)="fetchProfile()">Retry</button>
        </div>
        <div>
            <h1 *ngIf="nickname">{{nickname}} (member ID: {{userId}})</h1>
            <p *ngIf="email">Email: {{email}}</p>
        </div>
    `,
})
export class ProfileComponent {
    email = '';
    error: Error = null;
    nickname = '';
    userId = 0;

    constructor(private auth: Auth, private http: AuthHttp) {
        this.userId = auth.getUserId();

        this.fetchProfile();
    }

    private fetchProfile() {
        this.http.get("private/profile").subscribe(
            data => {
                this.nickname = data.profile.nickname;
                this.email = data.profile.email;
            },
            error => this.error = error
        );
    }
}
