import { Component } from '@angular/core';

@Component({
    selector: 'my-app',
    template: `
        <nav>
            <a routerLink="/does-not-exist" routerLinkActive="active">Does not exist</a>
            <a routerLink="/profile" routerLinkActive="active">Profile</a>
            <a routerLink="/login" routerLinkActive="active">Login</a>
        </nav>
        <router-outlet></router-outlet>
    `,
})
export class AppComponent { }
