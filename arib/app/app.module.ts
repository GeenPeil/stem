import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

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
		FormsModule,

		NgbModule.forRoot(),

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
