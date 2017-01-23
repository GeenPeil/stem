import { Component } from '@angular/core';

@Component({
    selector: 'gp-opstelten',
    template: `
        <simple-notifications [options]="notificationOptions"></simple-notifications>
        <template ngbModalContainer></template>

        <nav class="navbar navbar-toggleable-md navbar-light bg-faded">
            <button class="navbar-toggler navbar-toggler-right" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <a class="navbar-brand" href="#">GeenPeil Backoffice</a>
            <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
                <div class="navbar-nav">
                    <a class="nav-item nav-link" routerLinkActive="active" routerLink="/members" >Ledenlijst</a>
                    <a class="nav-item nav-link" routerLinkActive="active" routerLink="/member/new" >Nieuw lid</a>
                </div>
            </div>
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
