import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { RegisterRoutingModule } from './register-routing.module';

import { RegisterComponent } from './register.component';


@NgModule({
	imports: [
		BrowserModule,
		FormsModule,
		RegisterRoutingModule
	],
	providers: [],
	declarations: [
		RegisterComponent
	]
})
export class RegisterModule { }
