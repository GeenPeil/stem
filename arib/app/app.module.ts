import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { Ng2Webstorage } from 'ng2-webstorage';

import { AppRoutingModule } from './app-routing.module';
import { AuthModule } from './auth/auth.module';
import { RegisterModule } from './register/register.module';
import { ProfileModule } from './profile/profile.module';

import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './not-found.component';

// Add the RxJS Observable operators we need in this app.
import './rxjs-operators';

@NgModule({
	imports: [
		BrowserModule,

		NgbModule.forRoot(),
		Ng2Webstorage.forRoot({ prefix: 'gpiac-arib', separator: '.' }),

		AuthModule,
		RegisterModule,
		ProfileModule,
		AppRoutingModule
	],
	providers: [],
	declarations: [
		AppComponent,
		PageNotFoundComponent
	],
	bootstrap: [AppComponent]
})
export class AppModule { }
