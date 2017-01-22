import { Component } from '@angular/core';

@Component({
    selector: 'gp-opstelten',
    template: `
        <simple-notifications [options]="notificationOptions"></simple-notifications>
        <template ngbModalContainer></template>

        <nav>
            <a routerLink="/member/new" routerLinkActive="active">New member</a>
            <a routerLink="/member/1" routerLinkActive="active">Member 1</a>
            <a routerLink="/members" routerLinkActive="active">Members (non functional)</a>
        </nav>
        <router-outlet></router-outlet>`,
})
export class AppComponent {
    notificationOptions: Object = {
        timeOut: 7000,
        position: ["top", "right"],
        showProgressBar: true,
        pauseOnHover: true,
    };
}
