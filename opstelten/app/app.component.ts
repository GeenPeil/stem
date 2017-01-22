import { Component } from '@angular/core';

@Component({
    selector: 'gp-opstelten',
    template: `
      <nav>
          <a routerLink="/member/1" routerLinkActive="active">Member 1</a>
          <a routerLink="/members" routerLinkActive="active">Members (non functional)</a>
      </nav>
      <router-outlet></router-outlet>`,
})
export class AppComponent { }
