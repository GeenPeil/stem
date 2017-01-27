import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AuthModule } from './auth/auth.module';
import { ProfileModule } from './profile/profile.module';

import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './not-found.component';

// Add the RxJS Observable operators we need in this app.
import './rxjs-operators';

@NgModule({
	imports: [
		BrowserModule,
		FormsModule,
		AuthModule,
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
