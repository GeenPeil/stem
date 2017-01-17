import { Component } from '@angular/core';

import { Auth } from './auth.service';

@Component({
    template: `
        {{errorMessage}} | {{successMessage}}
        <input [(ngModel)]="id" type="number" >
        <button (click)="login()" >Log in</button>
    `,
})
export class LoginComponent {
    id: number;

    // feedback messages
    successMessage: string;
    errorMessage: string;

    constructor(private auth: Auth) { }

    // login tries to authenticate using the current form values
    login() {
        this.successMessage = ``;
        this.errorMessage = ``;
        this.auth.tryLogin(this.id).subscribe(
            () => this.successMessage = 'all good!',
            (error: any) => this.errorMessage = error
        );
    }
}
