import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { SimpleNotificationsModule } from 'angular2-notifications';

import { MembersModule } from './members/members.module';
import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './not-found.component';

// Add the RxJS Observable operators we need in this app.
import './rxjs-operators';

@NgModule({
    imports: [
        BrowserModule,
        
        NgbModule.forRoot(),
        SimpleNotificationsModule,
        
        MembersModule,
        AppRoutingModule
    ],
    declarations: [
        AppComponent,
        PageNotFoundComponent
    ],
    bootstrap: [AppComponent]
})
export class AppModule { }
